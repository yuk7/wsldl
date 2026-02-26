package get

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
	"github.com/yuk7/wsldl/lib/wtutils"
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

func TestExecute_DefaultUID_Success(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, wsllib.FlagAppendNTPath, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, nil
		},
	}

	if err := execute(wsl, reg, "Arch", []string{"--default-uid"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecute_InvalidArgLength_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{}

	err := execute(wsl, reg, "Arch", nil)
	assertDisplayError(t, err)
}

func TestExecute_GetConfigError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("config failed")
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 0, 0, 0, wantErr
		},
	}

	err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--default-uid"})
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_LxGuid_ProfileError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("profile failed")
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, wantErr
		},
	}

	err := execute(wsl, reg, "Arch", []string{"--lxguid"})
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_InvalidOption_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, nil
		},
	}

	err := execute(wsl, reg, "Arch", []string{"--unknown"})
	assertDisplayError(t, err)
}

func TestExecute_WslVersionAndDefaultTerm_Success(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, wsllib.FlagEnableWsl2, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{WsldlTerm: wsllib.FlagWsldlTermWT}, nil
		},
	}

	if err := execute(wsl, reg, "Arch", []string{"--wsl-version"}); err != nil {
		t.Fatalf("execute(--wsl-version) returned error: %v", err)
	}
	if err := execute(wsl, reg, "Arch", []string{"--default-term"}); err != nil {
		t.Fatalf("execute(--default-term) returned error: %v", err)
	}
}

func TestExecute_WtProfileName_FindsProfileByGeneratedGUID(t *testing.T) {
	distName := "Arch"
	guid := "{" + wtutils.CreateProfileGUID(distName) + "}"
	json := `{"profiles":{"list":[{"name":"Arch","guid":"` + guid + `"}]}}`

	tmp := t.TempDir()
	settingsPath := tmp + "\\Packages\\" + wtutils.WTPackageName + "\\LocalState\\settings.json"
	if err := os.WriteFile(settingsPath, []byte(json), 0o600); err != nil {
		t.Fatalf("write settings.json failed: %v", err)
	}
	t.Setenv("LOCALAPPDATA", tmp)

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{DistributionName: distName}, nil
		},
	}

	if err := execute(wsl, reg, "ignored", []string{"--wt-profile-name"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecute_LxGuid_EmptyUUIDWithoutProfileError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, nil
		},
	}

	err := execute(wsl, reg, "Arch", []string{"--lxguid"})
	de := assertDisplayError(t, err)
	if !strings.Contains(de.Unwrap().Error(), "lxguid get failed") {
		t.Fatalf("wrapped error = %v, want to contain %q", de.Unwrap(), "lxguid get failed")
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
