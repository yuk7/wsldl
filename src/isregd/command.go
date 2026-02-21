package isregd

import (
	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
)

func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"isregd"},
		Run:   execute,
	}
}

// execute is default isregd entrypoint. Exits with registerd status
func execute(name string, args []string) error {
	if wsllib.WslIsDistributionRegistered(name) {
		return nil
	}
	return utils.NewExitCodeError(1, false)
}
