package isregd

import (
	"os"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsllib-go"
)

func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"isregd"},
		Run:   execute,
	}
}

// execute is default isregd entrypoint. Exits with registerd status
func execute(name string, args []string) {
	if wsllib.WslIsDistributionRegistered(name) {
		os.Exit(0)
	}
	os.Exit(1)
}
