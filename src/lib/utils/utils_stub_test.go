//go:build !windows

package utils

import (
	"errors"
	"testing"
)

func TestCreateProcessAndWaitStub_ReturnsUnsupported(t *testing.T) {
	t.Parallel()

	code, err := CreateProcessAndWait("echo hello")
	if code != 0 {
		t.Fatalf("exit code = %d, want %d", code, 0)
	}
	if !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("error = %v, want %v", err, errUnsupportedPlatform)
	}
}
