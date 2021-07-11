package clean

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/yuk7/wsldl/lib/wslapi"
)

//Execute is default run entrypoint.
func Execute(name string, args []string) {
	showProgress := true
	switch len(args) {
	case 0:
		var in string
		fmt.Printf("This will remove this distro (%s) from the filesystem.\n", name)
		fmt.Printf("Are you sure you would like to proceed? (This cannot be undone)\n")
		fmt.Printf("Type \"y\" to continue:")
		fmt.Scan(&in)

		if in != "y" {
			fmt.Printf("Accepting is required to proceed.")
			os.Exit(1)
		}

	case 1:
		showProgress = false
		if args[0] == "-y" {
			showProgress = false
		} else {
			fmt.Println("Invalid Arg.")
			os.Exit(1)
		}

	default:
		fmt.Println("Invalid Arg.")
		os.Exit(1)
	}

	Clean(name, showProgress)
}

//Clean cleans distribution
func Clean(name string, showProgress bool) {
	if showProgress {
		fmt.Println("Unregistering...")
	}
	err := wslapi.WslUnregisterDistribution(name)
	if showProgress {
		if err != nil {
			fmt.Println("ERR: Failed to Unregister")
			var errno syscall.Errno
			if errors.As(err, &errno) {
				fmt.Printf("Code: 0x%x\n", int(errno))
				log.Fatal(err)
			}
		} else {
			fmt.Println("Unregistration complete")
		}
	}
	if err != nil {
		var errno syscall.Errno
		if errors.As(err, &errno) {
			os.Exit(int(errno))
		}
		os.Exit(1)
	}
	os.Exit(0)
}
