package backup

import (
	"errors"
	"os"
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

	tmp := t.TempDir()
	basePath := filepath.Join(tmp, "base")
	srcPath := basePath + "\\ext4.vhdx"
	payload := []byte("fake-vhdx")
	if err := os.WriteFile(srcPath, payload, 0o600); err != nil {
		t.Fatalf("write source vhdx: %v", err)
	}

	destPath := filepath.Join(tmp, "backup.ext4.vhdx")
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: basePath}, nil
		},
	}

	if err := execute(wsllib.MockWslLib{}, reg, "Arch", []string{destPath}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read destination: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("destination = %q, want %q", got, payload)
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
