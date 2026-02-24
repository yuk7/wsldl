package wslapi

import "github.com/yuk7/wsldl/lib/wsllib"

// GetConfig gets distribution configuration values.
func GetConfig(wsl wsllib.WslLib, distributionName string) (uid uint64, flags uint32, err error) {
	_, uid, flags, err = wsl.GetDistributionConfiguration(distributionName)
	return
}
