package clean

import (
	"fmt"
	"os"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

type cleanOptions struct {
	showProgress        bool
	requireConfirmation bool
}

// GetCommand returns the clean command structure
func GetCommand() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithDeps(deps.Wsl)
}

// GetCommandWithDeps returns the clean command structure with injectable dependencies.
func GetCommandWithDeps(wsl wsllib.WslLib) cmdline.Command {
	return cmdline.Command{
		Names: []string{"clean"},
		Help: func(distroName string, isListQuery bool) string {
			if wsl.IsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			}
			return ""
		},
		Run: func(name string, args []string) error {
			return execute(wsl, name, args)
		},
	}
}

// execute is default run entrypoint.
func execute(wsl wsllib.WslLib, name string, args []string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	return executeWithOptions(wsl, name, opts)
}

func parseArgs(args []string) (cleanOptions, error) {
	opts := cleanOptions{
		showProgress: true,
	}
	switch len(args) {
	case 0:
		opts.requireConfirmation = true

	case 1:
		if args[0] == "-y" {
			opts.showProgress = false
		} else {
			return cleanOptions{}, os.ErrInvalid
		}

	default:
		return cleanOptions{}, os.ErrInvalid
	}

	return opts, nil
}

func executeWithOptions(wsl wsllib.WslLib, name string, opts cleanOptions) error {
	if opts.requireConfirmation {
		var in string
		fmt.Printf("This will remove this distro (%s) from the filesystem.\n", name)
		fmt.Printf("Are you sure you would like to proceed? (This cannot be undone)\n")
		fmt.Printf("Type \"y\" to continue:")
		fmt.Scan(&in)

		if in != "y" {
			fmt.Fprintf(os.Stderr, "Accepting is required to proceed.")
			return errutil.NewDisplayError(os.ErrInvalid, false, true, false)
		}
	}

	err := Clean(wsl, name, opts.showProgress)
	if err != nil {
		return errutil.NewDisplayError(err, opts.showProgress, true, false)
	}
	return nil
}
