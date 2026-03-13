package config

import (
	"errors"
	"os"
	"runtime"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestParseArgs_DefaultTermNumeric_SetsOption(t *testing.T) {
	t.Parallel()

	opts, err := parseArgs([]string{"--default-term", "1"})
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if opts.option != configOptionDefaultTerm {
		t.Fatalf("option = %v, want %v", opts.option, configOptionDefaultTerm)
	}
	if opts.defaultTerm != wsllib.FlagWsldlTermWT {
		t.Fatalf("defaultTerm = %d, want %d", opts.defaultTerm, wsllib.FlagWsldlTermWT)
	}
}

func TestParseArgs_DefaultTermDefaultKeyword_SetsDefaultValue(t *testing.T) {
	t.Parallel()

	opts, err := parseArgs([]string{"--default-term", "default"})
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if opts.option != configOptionDefaultTerm {
		t.Fatalf("option = %v, want %v", opts.option, configOptionDefaultTerm)
	}
	if opts.defaultTerm != wsllib.FlagWsldlTermDefault {
		t.Fatalf("defaultTerm = %d, want %d", opts.defaultTerm, wsllib.FlagWsldlTermDefault)
	}
}

func TestParseArgs_InvalidAppendPathBool_ReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := parseArgs([]string{"--append-path", "not-bool"}); err == nil {
		t.Fatal("parseArgs succeeded unexpectedly")
	}
}

func TestParseArgs_AllSupportedFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want configOptions
	}{
		{
			name: "default uid",
			args: []string{"--default-uid", "1001"},
			want: configOptions{option: configOptionDefaultUID, uid: 1001},
		},
		{
			name: "default user",
			args: []string{"--default-user", "root"},
			want: configOptions{option: configOptionDefaultUser, user: "root"},
		},
		{
			name: "append path true",
			args: []string{"--append-path", "true"},
			want: configOptions{option: configOptionAppendPath, enabled: true},
		},
		{
			name: "mount drive false",
			args: []string{"--mount-drive", "false"},
			want: configOptions{option: configOptionMountDrive, enabled: false},
		},
		{
			name: "wsl version",
			args: []string{"--wsl-version", "2"},
			want: configOptions{option: configOptionWslVersion, wslVersion: 2},
		},
		{
			name: "default term keyword",
			args: []string{"--default-term", "wt"},
			want: configOptions{option: configOptionDefaultTerm, defaultTerm: wsllib.FlagWsldlTermWT},
		},
		{
			name: "default term numeric",
			args: []string{"--default-term", "2"},
			want: configOptions{option: configOptionDefaultTerm, defaultTerm: wsllib.FlagWsldlTermFlute},
		},
		{
			name: "flags value",
			args: []string{"--flags-val", "7"},
			want: configOptions{option: configOptionFlagsVal, flags: 7},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseArgs(tt.args)
			if err != nil {
				t.Fatalf("parseArgs returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("parseArgs = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseArgs_InvalidPatterns_ReturnError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		args          []string
		wantInvalidOS bool
	}{
		{name: "missing value", args: []string{"--append-path"}, wantInvalidOS: true},
		{name: "unknown option", args: []string{"--unknown", "1"}, wantInvalidOS: true},
		{name: "default uid non number", args: []string{"--default-uid", "abc"}},
		{name: "append path non bool", args: []string{"--append-path", "nope"}},
		{name: "mount drive non bool", args: []string{"--mount-drive", "nope"}},
		{name: "wsl version unsupported", args: []string{"--wsl-version", "3"}, wantInvalidOS: true},
		{name: "default term unknown", args: []string{"--default-term", "xterm"}, wantInvalidOS: true},
		{name: "flags val non number", args: []string{"--flags-val", "abc"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseArgs(tt.args)
			if err == nil {
				t.Fatal("parseArgs succeeded unexpectedly")
			}
			if tt.wantInvalidOS && !errors.Is(err, os.ErrInvalid) {
				t.Fatalf("error = %v, want os.ErrInvalid", err)
			}
		})
	}
}

func TestParseArgs_WslVersionNonNumber_ReturnsError(t *testing.T) {
	t.Parallel()

	_, err := parseArgs([]string{"--wsl-version", "x"})
	if err == nil {
		t.Fatal("parseArgs succeeded unexpectedly")
	}
}

func TestUpdateFlag_DisableWhenAlreadyOff_StaysOff(t *testing.T) {
	t.Parallel()

	flags := uint32(wsllib.FlagEnableDriveMounting)
	got := updateFlag(flags, wsllib.FlagAppendNTPath, false)
	if got != flags {
		t.Fatalf("flags = %d, want %d", got, flags)
	}
}

func TestUpdateFlag_DisableWhenOn_ClearsBit(t *testing.T) {
	t.Parallel()

	flags := uint32(wsllib.FlagEnableDriveMounting | wsllib.FlagAppendNTPath)
	got := updateFlag(flags, wsllib.FlagAppendNTPath, false)
	if got&wsllib.FlagAppendNTPath == wsllib.FlagAppendNTPath {
		t.Fatalf("flags = %d, want append-path bit cleared", got)
	}
	if got&wsllib.FlagEnableDriveMounting != wsllib.FlagEnableDriveMounting {
		t.Fatalf("flags = %d, want mount-drive bit preserved", got)
	}
}

func TestGetCommandWithDeps_HelpVisibility(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}
	cmd := GetCommandWithDeps(wsl, wsllib.MockWslReg{})
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
	if len(cmd.Names) != 2 || cmd.Names[0] != "config" || cmd.Names[1] != "set" {
		t.Fatalf("Names = %v, want [config set]", cmd.Names)
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

	err := cmd.Run("Arch", []string{"--bad", "1"})
	assertDisplayError(t, err)
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

func TestExecuteWithOptions_DefaultUser_SuccessUpdatesUID(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 4, nil
		},
		ConfigureDistributionFunc: func(name string, uid uint64, flags uint32) error {
			if name != "Arch" || uid != 2001 || flags != 4 {
				t.Fatalf("ConfigureDistribution args: name=%q uid=%d flags=%d", name, uid, flags)
			}
			return nil
		},
	}

	execRead := func(_ wsllib.WslLib, name string, cmd string) (string, uint32, error) {
		if name != "Arch" {
			t.Fatalf("name = %q, want %q", name, "Arch")
		}
		if cmd != "id -u root" {
			t.Fatalf("cmd = %q, want %q", cmd, "id -u root")
		}
		return "2001", 0, nil
	}

	err := executeWithOptionsAndExecRead(
		wsl,
		wsllib.MockWslReg{},
		"Arch",
		configOptions{option: configOptionDefaultUser, user: "root"},
		execRead,
	)
	if err != nil {
		t.Fatalf("executeWithOptionsAndExecRead returned error: %v", err)
	}
}

func TestExecuteWithOptions_DefaultUser_InvalidUIDOutput_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}

	execRead := func(_ wsllib.WslLib, _ string, _ string) (string, uint32, error) {
		return "not-a-number", 0, nil
	}

	err := executeWithOptionsAndExecRead(
		wsl,
		wsllib.MockWslReg{},
		"Arch",
		configOptions{option: configOptionDefaultUser, user: "root"},
		execRead,
	)
	de := assertDisplayError(t, err)
	if de.Unwrap().Error() != "not-a-number" {
		t.Fatalf("wrapped error = %v, want %q", de.Unwrap(), "not-a-number")
	}
}

func TestExecute_WslVersion_SetError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("set failed")
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{
		SetWslVersionFunc: func(name string, version int) error {
			return wantErr
		},
	}

	err := execute(wsl, reg, "Arch", []string{"--wsl-version", "2"})
	de := assertDisplayError(t, err)
	if !errors.Is(de.Unwrap(), wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_DefaultTerm_GetProfileError_ReturnsDisplayError(t *testing.T) {
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

	err := execute(wsl, reg, "Arch", []string{"--default-term", "wt"})
	de := assertDisplayError(t, err)
	if !errors.Is(de.Unwrap(), wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_DefaultTerm_WriteProfileError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("write failed")
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{DistributionName: name}, nil
		},
		WriteProfileFunc: func(profile wsllib.Profile) error {
			return wantErr
		},
	}

	err := execute(wsl, reg, "Arch", []string{"--default-term", "wt"})
	de := assertDisplayError(t, err)
	if !errors.Is(de.Unwrap(), wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_MountDriveTrue_SetsFlag(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
		ConfigureDistributionFunc: func(name string, uid uint64, flags uint32) error {
			if flags&wsllib.FlagEnableDriveMounting != wsllib.FlagEnableDriveMounting {
				t.Fatalf("flags = %d, want mount-drive bit set", flags)
			}
			return nil
		},
	}
	if err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--mount-drive", "true"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteWithOptions_InvalidOption_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	err := executeWithOptions(wsl, wsllib.MockWslReg{}, "Arch", configOptions{option: configOption(999)})
	de := assertDisplayError(t, err)
	if !errors.Is(de.Unwrap(), os.ErrInvalid) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), os.ErrInvalid)
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
