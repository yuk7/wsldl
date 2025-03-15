package isregd

import (
	"os"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsllib-go"
)

func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"isregd"},
		Run: func(distroName string, args []string) {
			Execute(distroName)
		},
	}
}

// Execute is default isregd entrypoint. Exits with registerd status
func Execute(name string) {
	if wsllib.WslIsDistributionRegistered(name) {
		os.Exit(0)
	}
	os.Exit(1)
}
