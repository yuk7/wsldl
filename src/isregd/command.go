package isregd

import (
	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func GetCommand() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithDeps(deps.Wsl)
}

func GetCommandWithDeps(wsl wsllib.WslLib) cmdline.Command {
	return cmdline.Command{
		Names: []string{"isregd"},
		Run: func(name string, args []string) error {
			return execute(wsl, name, args)
		},
	}
}

// execute is default isregd entrypoint. Exits with registerd status
func execute(wsl wsllib.WslLib, name string, args []string) error {
	if wsl.IsDistributionRegistered(name) {
		return nil
	}
	return errutil.NewExitCodeError(1, false)
}
