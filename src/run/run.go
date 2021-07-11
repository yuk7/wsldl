package run

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wslapi"
)

//Execute is default run entrypoint.
func Execute(name string, args []string) {
	command := ""
	for _, s := range args {
		command = command + " " + utils.DQEscapeString(s)
	}

	exitCode, err := wslapi.WslLaunchInteractive(name, command, true)
	var errno syscall.Errno
	if errors.As(err, &errno) {
		fmt.Printf("ERR: Launch Process failed\n")
		fmt.Printf("Code: 0x%x\nExit Code:0x%x", int(errno), exitCode)
		log.Fatal(err)
	} else {
		os.Exit(int(exitCode))
	}
}

//ExecuteP runs Execute function with Path Translator
func ExecuteP(name string, args []string) {
	var convArgs []string
	for _, s := range args {
		if strings.Contains(s, "\\") {
			s = strings.Replace(s, "\\", "/", -1)
			s = utils.DQEscapeString(s)
			out, exitCode, err := ExecRead(name, "wslpath -u "+s)
			if err != nil || exitCode != 0 {
				fmt.Println("ERR: Failed to Path Translation")
				var errno syscall.Errno
				if errors.As(err, &errno) {
					fmt.Printf("Code: 0x%x\n", int(errno))
				}
				fmt.Printf("ExitCode: 0x%x\n", exitCode)
				if err != nil {
					log.Fatal(err)
				}
				os.Exit(int(exitCode))
			}
			convArgs = append(convArgs, out)
		} else {
			convArgs = append(convArgs, s)
		}
	}

	Execute(name, convArgs)
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

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    <no args>")
	println("      - Open a new shell with your default settings.")
	println()
	println("    run <command line>")
	println("      - Run the given command line in that distro. Inherit current directory.")
	println()
	println("    runp <command line (includes windows path)>")
	println("      - Run the path translated command line in that distro.")
}
