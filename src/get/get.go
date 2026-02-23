package get

import (
	"github.com/yuk7/wsldl/lib/wsllib"
)

// WslGetConfig is getter of distribution configuration
func WslGetConfig(wsl wsllib.WslLib, distributionName string) (uid uint64, flags uint32, err error) {
	_, uid, flags, err = wsl.GetDistributionConfiguration(distributionName)
	return
}
