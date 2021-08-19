package get

import (
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wslapi"
)

//WslGetConfig is getter of distribution configuration
func WslGetConfig(distributionName string) (uid uint64, flags uint32) {
	var err error
	_, uid, flags, err = wslapi.WslGetDistributionConfiguration(distributionName)
	if err != nil {
		utils.ErrorRedPrintln("ERR: Failed to GetDistributionConfiguration")
		utils.ErrorExit(err, true, true, false)
	}
	return
}
