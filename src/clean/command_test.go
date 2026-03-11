package clean

import (
	"errors"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestParseArgs_NoArgs_RequiresConfirmation(t *testing.T) {
	t.Parallel()

	opts, err := parseArgs(nil)
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if !opts.showProgress {
		t.Fatal("showProgress = false, want true")
	}
	if !opts.requireConfirmation {
		t.Fatal("requireConfirmation = false, want true")
	}
}

func TestParseArgs_WithY_DisablesProgress(t *testing.T) {
	t.Parallel()

	opts, err := parseArgs([]string{"-y"})
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if opts.showProgress {
		t.Fatal("showProgress = true, want false")
	}
	if opts.requireConfirmation {
		t.Fatal("requireConfirmation = true, want false")
	}
}

func TestGetCommandWithDeps_HelpVisibility(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}
	cmd := GetCommandWithDeps(wsl)
	if cmd.Visible == nil {
		t.Fatal("Visible is nil")
	}
	if cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = true, want false")
	}
	if cmd.HelpText == nil {
		t.Fatal("HelpText is nil")
	}
	if got := cmd.HelpText(); got == "" {
		t.Fatal("HelpText should not be empty")
	}
}

func TestExecute_WithY_CallsUnregister(t *testing.T) {
	t.Parallel()

	called := 0
	wsl := wsllib.MockWslLib{
		UnregisterDistributionFunc: func(name string) error {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return nil
		},
	}

	err := execute(wsl, "Arch", []string{"-y"})
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("UnregisterDistribution call count = %d, want 1", called)
	}
}

func TestExecute_InvalidArg_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	err := execute(wsllib.MockWslLib{}, "Arch", []string{"--bad"})
	assertDisplayError(t, err)
}

func TestClean_PropagatesError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("unregister failed")
	wsl := wsllib.MockWslLib{
		UnregisterDistributionFunc: func(name string) error {
			return wantErr
		},
	}
	err := Clean(wsl, "Arch", false)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Clean error = %v, want %v", err, wantErr)
	}
}

func assertDisplayError(t *testing.T, err error) *errutil.DisplayError {
	t.Helper()
	var de *errutil.DisplayError
	if !errors.As(err, &de) {
		t.Fatalf("error type = %T, want *errutil.DisplayError", err)
	}
	return de
}
