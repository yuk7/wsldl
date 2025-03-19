package clean

import (
	"fmt"
	"os"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
)

// GetCommand returns the clean command structure
func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"clean"},
		Help: func(distroName string, isListQuery bool) string {
			if wsllib.WslIsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			}
			return ""
		},
		Run: execute,
	}
}

// execute is default run entrypoint.
func execute(name string, args []string) {
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
