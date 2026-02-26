package run

import (
	"errors"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestGetCommandWithNoArgsWithDeps_HelpVisibility(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}
	cmd := GetCommandWithNoArgsWithDeps(wsl, wsllib.MockWslReg{})
	if got := cmd.Help("Arch", true); got != "" {
		t.Fatalf("Help(list query) = %q, want empty", got)
	}
	if got := cmd.Help("Arch", false); got == "" {
		t.Fatal("Help(non-list query) should not be empty")
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
	if got := cmd.Help("Arch", true); got != "" {
		t.Fatalf("Help(list query) = %q, want empty", got)
	}
	if got := cmd.Help("Arch", false); got == "" {
		t.Fatal("Help(non-list query) should not be empty")
	}
}

func TestGetCommandPWithDeps_HelpVisibility(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}
	cmd := GetCommandPWithDeps(wsl)
	if got := cmd.Help("Arch", true); got != "" {
		t.Fatalf("Help(list query) = %q, want empty", got)
	}
	if got := cmd.Help("Arch", false); got == "" {
		t.Fatal("Help(non-list query) should not be empty")
	}
}

func TestExecute_LaunchInteractiveSuccess(t *testing.T) {
	t.Parallel()

	called := 0
	wsl := wsllib.MockWslLib{
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if command != ` echo "hello world"` {
				t.Fatalf("command = %q, want %q", command, ` echo "hello world"`)
			}
			if !inheritPath {
				t.Fatal("inheritPath = false, want true")
			}
			return 0, nil
		},
	}

	if err := execute(wsl, "Arch", []string{"echo", "hello world"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("LaunchInteractive call count = %d, want 1", called)
	}
}

func TestExecute_LaunchErrorReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("launch failed")
	wsl := wsllib.MockWslLib{
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			return 0, wantErr
		},
	}

	err := execute(wsl, "Arch", []string{"echo"})
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_NonZeroExitReturnsExitCodeError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			return 37, nil
		},
	}

	err := execute(wsl, "Arch", []string{"echo"})
	ee := assertExitCodeError(t, err)
	if ee.Code != 37 {
		t.Fatalf("exit code = %d, want %d", ee.Code, 37)
	}
	if ee.Pause {
		t.Fatal("pause = true, want false")
	}
}

func TestExecuteP_NoPathTranslation_DelegatesToExecute(t *testing.T) {
	t.Parallel()

	called := 0
	wsl := wsllib.MockWslLib{
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if command != ` ls -la` {
				t.Fatalf("command = %q, want %q", command, ` ls -la`)
			}
			if !inheritPath {
				t.Fatal("inheritPath = false, want true")
			}
			return 0, nil
		},
	}

	if err := executeP(wsl, "Arch", []string{"ls", "-la"}); err != nil {
		t.Fatalf("executeP returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("LaunchInteractive call count = %d, want 1", called)
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

func assertExitCodeError(t *testing.T, err error) *errutil.ExitCodeError {
	t.Helper()
	var ee *errutil.ExitCodeError
	if !errors.As(err, &ee) {
		t.Fatalf("error type = %T, want *errutil.ExitCodeError", err)
	}
	return ee
}
