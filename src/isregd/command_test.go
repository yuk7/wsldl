package isregd

import (
	"errors"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestExecute_WhenRegistered_ReturnsNil(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return true
		},
	}

	err := execute(wsl, "Arch", nil)
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecute_WhenNotRegistered_ReturnsExitCodeOne(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}

	err := execute(wsl, "Arch", nil)
	var ex *errutil.ExitCodeError
	if !errors.As(err, &ex) {
		t.Fatalf("execute error type = %T, want *errutil.ExitCodeError", err)
	}
	if ex.Code != 1 || ex.Pause {
		t.Fatalf("exit code error = %+v, want Code=1 Pause=false", ex)
	}
}

func TestGetCommandWithDeps_WiresRun(t *testing.T) {
	t.Parallel()

	cmd := GetCommandWithDeps(wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return true
		},
	})
	if len(cmd.Names) != 1 || cmd.Names[0] != "isregd" {
		t.Fatalf("Names = %v, want [isregd]", cmd.Names)
	}
	if cmd.Run == nil {
		t.Fatal("Run is nil")
	}

	if err := cmd.Run("Arch", nil); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
}

func TestGetCommand_WiresDefaultDeps(t *testing.T) {
	t.Parallel()

	cmd := GetCommand()
	if len(cmd.Names) != 1 || cmd.Names[0] != "isregd" {
		t.Fatalf("Names = %v, want [isregd]", cmd.Names)
	}
	if cmd.Run == nil {
		t.Fatal("Run is nil")
	}

	err := cmd.Run("Arch", nil)
	var ex *errutil.ExitCodeError
	if !errors.As(err, &ex) {
		t.Fatalf("Run error type = %T, want *errutil.ExitCodeError", err)
	}
	if ex.Code != 1 || ex.Pause {
		t.Fatalf("exit code error = %+v, want Code=1 Pause=false", ex)
	}
}
