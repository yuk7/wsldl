package clean

import (
	"fmt"
	"os"

	"github.com/yuk7/wsldl/lib/utils"
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
			fmt.Fprintf(os.Stderr, "Accepting is required to proceed.")
			utils.ErrorExit(os.ErrInvalid, false, true, false)
		}

	case 1:
		showProgress = false
		if args[0] == "-y" {
			showProgress = false
		} else {
			utils.ErrorExit(os.ErrInvalid, true, true, false)
		}

	default:
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}

	Clean(name, showProgress)
}

//Clean cleans distribution
func Clean(name string, showProgress bool) {
	if showProgress {
		fmt.Println("Unregistering...")
	}
	err := wslapi.WslUnregisterDistribution(name)
	if err != nil {
		utils.ErrorExit(err, showProgress, true, false)
	}
	os.Exit(0)
}

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    clean")
	println("      - Uninstall that instance.")
}
