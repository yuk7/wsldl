package config

import (
	"errors"
	"runtime"
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

func TestExecute_DefaultUID_ConfigureDistribution(t *testing.T) {
	t.Parallel()

	configured := false
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
		ConfigureDistributionFunc: func(name string, uid uint64, flags uint32) error {
			configured = true
			if name != "Arch" || uid != 2000 || flags != 0 {
				t.Fatalf("ConfigureDistribution args: name=%q uid=%d flags=%d", name, uid, flags)
			}
			return nil
		},
	}

	err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--default-uid", "2000"})
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if !configured {
		t.Fatal("ConfigureDistribution was not called")
	}
}

func TestExecute_DefaultTerm_UpdatesProfile(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
		ConfigureDistributionFunc: func(name string, uid uint64, flags uint32) error {
			return nil
		},
	}

	written := wsllib.Profile{}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{DistributionName: name}, nil
		},
		WriteProfileFunc: func(profile wsllib.Profile) error {
			written = profile
			return nil
		},
	}

	err := execute(wsl, reg, "Arch", []string{"--default-term", "wt"})
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if written.WsldlTerm != wsllib.FlagWsldlTermWT {
		t.Fatalf("written term = %d, want %d", written.WsldlTerm, wsllib.FlagWsldlTermWT)
	}
}

func TestExecute_InvalidArgs_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--append-path"})
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
	err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--default-uid", "1000"})
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_ConfigureDistributionError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("configure failed")
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
		ConfigureDistributionFunc: func(name string, uid uint64, flags uint32) error {
			return wantErr
		},
	}
	err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--flags-val", "7"})
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_AppendPathTrue_SetsFlag(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
		ConfigureDistributionFunc: func(name string, uid uint64, flags uint32) error {
			if uid != 1000 {
				t.Fatalf("uid = %d, want %d", uid, 1000)
			}
			if flags&wsllib.FlagAppendNTPath != wsllib.FlagAppendNTPath {
				t.Fatalf("flags = %d, want append-path bit set", flags)
			}
			return nil
		},
	}

	if err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--append-path", "true"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecute_MountDriveFalse_ClearsFlag(t *testing.T) {
	t.Parallel()

	initialFlags := uint32(wsllib.FlagEnableDriveMounting | wsllib.FlagAppendNTPath)
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, initialFlags, nil
		},
		ConfigureDistributionFunc: func(name string, uid uint64, flags uint32) error {
			if flags&wsllib.FlagEnableDriveMounting == wsllib.FlagEnableDriveMounting {
				t.Fatalf("flags = %d, want mount-drive bit cleared", flags)
			}
			if flags&wsllib.FlagAppendNTPath != wsllib.FlagAppendNTPath {
				t.Fatalf("flags = %d, append-path bit should remain set", flags)
			}
			return nil
		},
	}

	if err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--mount-drive", "false"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecute_WslVersion_ValidCallsSetWslVersion(t *testing.T) {
	t.Parallel()

	setCalled := false
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
		ConfigureDistributionFunc: func(name string, uid uint64, flags uint32) error {
			return nil
		},
	}
	reg := wsllib.MockWslReg{
		SetWslVersionFunc: func(name string, version int) error {
			setCalled = true
			if name != "Arch" || version != 2 {
				t.Fatalf("SetWslVersion args: name=%q version=%d", name, version)
			}
			return nil
		},
	}

	if err := execute(wsl, reg, "Arch", []string{"--wsl-version", "2"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if !setCalled {
		t.Fatal("SetWslVersion was not called")
	}
}

func TestExecute_WslVersion_InvalidValue_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--wsl-version", "3"})
	assertDisplayError(t, err)
}

func TestExecute_DefaultUser_OnNonWindows_ReturnsDisplayError(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("non-windows stub behavior only")
	}

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--default-user", "root"})
	assertDisplayError(t, err)
}

func assertDisplayError(t *testing.T, err error) *errutil.DisplayError {
	t.Helper()
	var de *errutil.DisplayError
	if !errors.As(err, &de) {
		t.Fatalf("error type = %T, want *errutil.DisplayError", err)
	}
	return de
}
