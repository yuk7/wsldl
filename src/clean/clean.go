package clean

import (
	"fmt"

	"github.com/yuk7/wsldl/lib/wsllib"
)

// Clean cleans distribution
func Clean(wsl wsllib.WslLib, name string, showProgress bool) error {
	if showProgress {
		fmt.Println("Unregistering...")
	}
	err := wsl.UnregisterDistribution(name)
	if err != nil {
		return err
	}
	return nil
}
