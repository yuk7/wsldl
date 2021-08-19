package clean

import (
	"fmt"
	"os"

	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wslapi"
)

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
