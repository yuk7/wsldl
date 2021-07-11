package run

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/yuk7/wsldl/lib/wslapi"
)

//Execute is default run entrypoint.
func Execute(name string, args []string) {
	command := ""
	for _, s := range args {
		command = command + " " + s
	}

	exitCode, err := wslapi.WslLaunchInteractive(name, command, true)
	var errno syscall.Errno
	if errors.As(err, &errno) {
		fmt.Printf("ERR: Launch Process failed\n")
		fmt.Printf("Code: 0x%x\nExit Code:0x%x", int(errno), exitCode)
	} else {
		os.Exit(int(exitCode))
	}
}

//ExecRead execs command and read output
func ExecRead(name, command string) (out string, exitCode uint32, err error) {
	stdin := syscall.Handle(0)
	stdout := syscall.Handle(0)
	stdintmp := syscall.Handle(0)
	stdouttmp := syscall.Handle(0)
	sa := syscall.SecurityAttributes{InheritHandle: 1, SecurityDescriptor: 0}

	syscall.CreatePipe(&stdin, &stdintmp, &sa, 0)
	syscall.CreatePipe(&stdout, &stdouttmp, &sa, 0)

	handle, err := wslapi.WslLaunch(name, command, true, stdintmp, stdouttmp, stdouttmp)
	syscall.WaitForSingleObject(handle, syscall.INFINITE)
	syscall.GetExitCodeProcess(handle, &exitCode)
	buf := make([]byte, syscall.MAX_LONG_PATH)
	var length uint32

	syscall.ReadFile(stdout, buf, &length, nil)

	//[]byte -> string and cut to fit the length
	out = string(buf)[:length]
	if out[len(out)-1:] == "\n" {
		out = out[:len(out)-1]
	}
	return
}
