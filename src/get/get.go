package get

import (
	"github.com/yuk7/wsllib-go"
)

// WslGetConfig is getter of distribution configuration
func WslGetConfig(distributionName string) (uid uint64, flags uint32, err error) {
	_, uid, flags, err = wsllib.WslGetDistributionConfiguration(distributionName)
	return
}
