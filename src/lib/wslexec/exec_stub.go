//go:build !windows

package wslexec

import (
	"errors"

	"github.com/yuk7/wsldl/lib/wsllib"
)

var errUnsupportedPlatform = errors.New("wslexec is only available on Windows")

// ExecRead executes a command and reads output.
func ExecRead(wsl wsllib.WslLib, name, command string) (out string, exitCode uint32, err error) {
	return "", 0, errUnsupportedPlatform
}
