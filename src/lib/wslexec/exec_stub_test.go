//go:build !windows

package wslexec

import (
	"errors"
	"testing"

	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestExecRead_ReturnsUnsupportedOnNonWindows(t *testing.T) {
	t.Parallel()

	out, code, err := ExecRead(wsllib.MockWslLib{}, "Arch", "id -u")
	if out != "" {
		t.Fatalf("out = %q, want empty", out)
	}
	if code != 0 {
		t.Fatalf("exitCode = %d, want 0", code)
	}
	if !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("err = %v, want %v", err, errUnsupportedPlatform)
	}
}
