package backup

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestGetCommandWithDeps_HelpVisibility(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}
	cmd := GetCommandWithDeps(wsl, wsllib.MockWslReg{})
	if got := cmd.Help("Arch", true); got != "" {
		t.Fatalf("Help(list query) = %q, want empty", got)
	}
	if got := cmd.Help("Arch", false); got == "" {
		t.Fatal("Help(non-list query) should not be empty")
	}
}

func TestExecute_InvalidArgLength_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	err := execute(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", []string{"a", "b"})
	assertDisplayError(t, err)
}

func TestExecute_InvalidSingleArg_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	err := execute(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", []string{"--unknown"})
	assertDisplayError(t, err)
}

func TestExecute_RegOption_ProfileError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("profile failed")
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, wantErr
		},
	}

	err := execute(wsllib.MockWslLib{}, reg, "Arch", []string{"--reg"})
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_VhdxCustomPath_Success(t *testing.T) {
	t.Parallel()

	destPath := filepath.Join(t.TempDir(), "backup.ext4.vhdx")
	gotName := ""
	gotDestPath := ""
	called := false

	err := executeWithBackups(
		wsllib.MockWslLib{},
		wsllib.MockWslReg{},
		"Arch",
		[]string{destPath},
		func(wsllib.WslReg, string, string) error {
			t.Fatal("backupReg should not be called")
			return nil
		},
		func(string, string) error {
			t.Fatal("backupTar should not be called")
			return nil
		},
		func(_ wsllib.WslReg, name, dest string) error {
			called = true
			gotName = name
			gotDestPath = dest
			return nil
		},
	)
	if err != nil {
		t.Fatalf("executeWithBackups returned error: %v", err)
	}
	if !called {
		t.Fatal("backupExt4Vhdx should be called")
	}
	if gotName != "Arch" {
		t.Fatalf("name = %q, want %q", gotName, "Arch")
	}
	if gotDestPath != destPath {
		t.Fatalf("dest path = %q, want %q", gotDestPath, destPath)
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
