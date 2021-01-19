package isregd

import (
	"os"

	"github.com/yuk7/wsldl/lib/wslapi"
)

//Execute is default isregd entrypoint. Exits with registerd status
func Execute(name string) {

	if wslapi.WslIsDistributionRegistered(name) {
		os.Exit(0)
	}
	os.Exit(1)
}
