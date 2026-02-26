//go:build !windows

package utils

import "errors"

var errUnsupportedPlatform = errors.New("unsupported platform")

// CreateProcessAndWait is unsupported outside Windows.
func CreateProcessAndWait(commandLine string) (int, error) {
	return 0, errUnsupportedPlatform
}
