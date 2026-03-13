package clean

import (
	"errors"
	"io"
	"os"
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

func TestParseArgs_TooManyArgs_ReturnsError(t *testing.T) {
	t.Parallel()

	_, err := parseArgs([]string{"-y", "extra"})
	if !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("error = %v, want %v", err, os.ErrInvalid)
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

func TestGetCommand_WiresDefaultDeps(t *testing.T) {
	t.Parallel()

	cmd := GetCommand()
	if len(cmd.Names) != 1 || cmd.Names[0] != "clean" {
		t.Fatalf("Names = %v, want [clean]", cmd.Names)
	}
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
	if cmd.Run == nil {
		t.Fatal("Run is nil")
	}

	err := cmd.Run("Arch", []string{"--bad"})
	assertDisplayError(t, err)
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

func TestClean_ShowProgressTrue_Success(t *testing.T) {
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

	if err := Clean(wsl, "Arch", true); err != nil {
		t.Fatalf("Clean returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("UnregisterDistribution call count = %d, want 1", called)
	}
}

func TestExecuteWithOptions_ConfirmationRejected_ReturnsDisplayError(t *testing.T) {
	called := 0
	wsl := wsllib.MockWslLib{
		UnregisterDistributionFunc: func(name string) error {
			called++
			return nil
		},
	}

	origStdin := os.Stdin
	inR, inW, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdin pipe failed: %v", err)
	}
	os.Stdin = inR
	t.Cleanup(func() {
		os.Stdin = origStdin
		_ = inR.Close()
		_ = inW.Close()
	})

	if _, err := io.WriteString(inW, "n\n"); err != nil {
		t.Fatalf("write stdin failed: %v", err)
	}
	if err := inW.Close(); err != nil {
		t.Fatalf("close stdin writer failed: %v", err)
	}

	err = executeWithOptions(wsl, "Arch", cleanOptions{
		showProgress:        false,
		requireConfirmation: true,
	})
	de := assertDisplayError(t, err)
	if !errors.Is(de, os.ErrInvalid) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), os.ErrInvalid)
	}
	if called != 0 {
		t.Fatalf("UnregisterDistribution call count = %d, want 0", called)
	}
}

func TestExecuteWithOptions_ConfirmationAccepted_CallsClean(t *testing.T) {
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

	origStdin := os.Stdin
	inR, inW, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdin pipe failed: %v", err)
	}
	os.Stdin = inR
	t.Cleanup(func() {
		os.Stdin = origStdin
		_ = inR.Close()
		_ = inW.Close()
	})

	if _, err := io.WriteString(inW, "y\n"); err != nil {
		t.Fatalf("write stdin failed: %v", err)
	}
	if err := inW.Close(); err != nil {
		t.Fatalf("close stdin writer failed: %v", err)
	}

	err = executeWithOptions(wsl, "Arch", cleanOptions{
		showProgress:        false,
		requireConfirmation: true,
	})
	if err != nil {
		t.Fatalf("executeWithOptions returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("UnregisterDistribution call count = %d, want 1", called)
	}
}

func TestExecuteWithOptions_CleanError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("unregister failed")
	wsl := wsllib.MockWslLib{
		UnregisterDistributionFunc: func(name string) error {
			return wantErr
		},
	}

	err := executeWithOptions(wsl, "Arch", cleanOptions{
		showProgress:        false,
		requireConfirmation: false,
	})
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
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
