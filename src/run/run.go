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
