package install

import (
	"strings"
	"testing"
)

func TestGetHelpMessageNoArgs_NotEmpty(t *testing.T) {
	t.Parallel()

	got := getHelpMessageNoArgs()
	if got == "" {
		t.Fatal("getHelpMessageNoArgs returned empty string")
	}
	if !strings.Contains(got, "Install a new instance") {
		t.Fatalf("help text = %q, want to contain %q", got, "Install a new instance")
	}
}

func TestGetHelpMessage_NotEmpty(t *testing.T) {
	t.Parallel()

	got := getHelpMessage()
	if got == "" {
		t.Fatal("getHelpMessage returned empty string")
	}
	if !strings.Contains(got, "install [rootfs file]") {
		t.Fatalf("help text = %q, want to contain %q", got, "install [rootfs file]")
	}
}
