package clean

import (
	"fmt"

	"github.com/yuk7/wsllib-go"
)

// Clean cleans distribution
func Clean(name string, showProgress bool) error {
	if showProgress {
		fmt.Println("Unregistering...")
	}
	err := wsllib.WslUnregisterDistribution(name)
	if err != nil {
		return err
	}
	return nil
}
