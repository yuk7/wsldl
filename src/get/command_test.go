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

func TestParseArgs_WtProfileAlias(t *testing.T) {
	t.Parallel()

	opts, err := parseArgs([]string{"--wt-pn"})
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if opts.option != getOptionWTProfileName {
		t.Fatalf("option = %v, want %v", opts.option, getOptionWTProfileName)
	}
}

func TestParseArgs_InvalidOption_ReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := parseArgs([]string{"--bad"}); err == nil {
		t.Fatal("parseArgs succeeded unexpectedly")
	}
}

func TestParseArgs_AllSupportedOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  string
		want getOption
	}{
		{name: "default uid", arg: "--default-uid", want: getOptionDefaultUID},
		{name: "append path", arg: "--append-path", want: getOptionAppendPath},
		{name: "mount drive", arg: "--mount-drive", want: getOptionMountDrive},
		{name: "wsl version", arg: "--wsl-version", want: getOptionWslVersion},
		{name: "lxguid", arg: "--lxguid", want: getOptionLXGuid},
		{name: "lxuid alias", arg: "--lxuid", want: getOptionLXGuid},
		{name: "default terminal", arg: "--default-terminal", want: getOptionDefaultTerm},
		{name: "wt profile name", arg: "--wt-profile-name", want: getOptionWTProfileName},
		{name: "wt profile alias", arg: "--wt-profilename", want: getOptionWTProfileName},
		{name: "flags val", arg: "--flags-val", want: getOptionFlagsVal},
		{name: "flags bits", arg: "--flags-bits", want: getOptionFlagsBits},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseArgs([]string{tt.arg})
			if err != nil {
				t.Fatalf("parseArgs returned error: %v", err)
			}
			if got.option != tt.want {
				t.Fatalf("option = %v, want %v", got.option, tt.want)
			}
		})
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
	if len(cmd.Names) != 1 || cmd.Names[0] != "get" {
		t.Fatalf("Names = %v, want [get]", cmd.Names)
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
	conf := wtutils.Config{}
	conf.Profiles.ProfileList = []wtutils.Profile{
		{Name: "Arch", GUID: guid},
	}

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

	if err := executeWithWTConfigReader(wsl, reg, "ignored", []string{"--wt-profile-name"}, func() (wtutils.Config, error) {
		return conf, nil
	}); err != nil {
		t.Fatalf("executeWithWTConfigReader returned error: %v", err)
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

func TestExecute_ParseArgsLengthError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	err := execute(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", []string{"--default-uid", "1000"})
	de := assertDisplayError(t, err)
	if !errors.Is(de.Unwrap(), os.ErrInvalid) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), os.ErrInvalid)
	}
}

func TestExecute_GetOptionVariants_Success(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, wsllib.FlagAppendNTPath | wsllib.FlagEnableDriveMounting, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{
				UUID:             "{abc}",
				WsldlTerm:        wsllib.FlagWsldlTermFlute,
				DistributionName: "",
			}, nil
		},
	}

	for _, arg := range []string{"--append-path", "--mount-drive", "--flags-val", "--flags-bits", "--default-term", "--lxguid"} {
		if err := execute(wsl, reg, "Arch", []string{arg}); err != nil {
			t.Fatalf("execute(%s) returned error: %v", arg, err)
		}
	}
}

func TestExecute_WslVersion_FlagUnset_PrintsVersion1WithoutError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	if err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"--wsl-version"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecute_DefaultTerm_DefaultBranchWithoutError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{WsldlTerm: 9999}, nil
		},
	}

	if err := execute(wsl, reg, "Arch", []string{"--default-term"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecute_WtProfileName_ReadConfigError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("read wt config failed")
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{DistributionName: "Arch"}, nil
		},
	}

	err := executeWithWTConfigReader(wsl, reg, "Arch", []string{"--wt-profile-name"}, func() (wtutils.Config, error) {
		return wtutils.Config{}, wantErr
	})
	de := assertDisplayError(t, err)
	if !errors.Is(de.Unwrap(), wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_WtProfileName_ProfileNotFound_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 2, 1000, 0, nil
		},
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{DistributionName: "Arch"}, nil
		},
	}

	err := executeWithWTConfigReader(wsl, reg, "Arch", []string{"--wt-profile-name"}, func() (wtutils.Config, error) {
		return wtutils.Config{}, nil
	})
	de := assertDisplayError(t, err)
	if !strings.Contains(de.Unwrap().Error(), "profile not found") {
		t.Fatalf("wrapped error = %v, want to contain %q", de.Unwrap(), "profile not found")
	}
}

func TestExecuteWithOptions_InvalidOption_ReturnsDisplayError(t *testing.T) {
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

	err := executeWithOptions(wsl, reg, "Arch", getOptions{option: getOption(999)}, func() (wtutils.Config, error) {
		t.Fatal("readWTConfig should not be called for invalid option")
		return wtutils.Config{}, nil
	})
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
