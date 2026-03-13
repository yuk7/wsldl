package run

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/yuk7/wsldl/lib/wsllib"
	"github.com/yuk7/wsldl/lib/wtutils"
)

func TestExecWindowsTerminal_NonWindows_ReturnsDisplayError(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("skip on windows to avoid launching real terminal process")
	}

	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{
				DistributionName: "ArchLinux",
			}, nil
		},
	}

	err := ExecWindowsTerminal(reg, "arch")
	de := assertDisplayError(t, err)
	if !strings.Contains(de.Error(), "unsupported platform") {
		t.Fatalf("error = %q, want to contain %q", de.Error(), "unsupported platform")
	}
}

func TestExecWindowsTerminalWithDeps_UsesProfileNameFromGUID(t *testing.T) {
	t.Parallel()

	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{DistributionName: "Arch"}, nil
		},
	}
	conf := wtutils.Config{}
	conf.Profiles.ProfileList = []wtutils.Profile{
		{GUID: "{guid-1234}", Name: "Arch Profile"},
	}

	gotCommand := ""
	deps := execWindowsTerminalDeps{
		readParseWTConfig: func() (wtutils.Config, error) {
			return conf, nil
		},
		createProfileGUID: func(name string) string {
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return "guid-1234"
		},
		getenv:         func(string) string { return `C:\Users\user\AppData\Local` },
		mustExecutable: func() string { t.Fatal("mustExecutable should not be called"); return "" },
		createProcessAndWait: func(commandLine string) (int, error) {
			gotCommand = commandLine
			return 0, nil
		},
		allocConsole: func() {},
	}

	err := execWindowsTerminalWithDeps(reg, "arch", deps)
	if err != nil {
		t.Fatalf("execWindowsTerminalWithDeps returned error: %v", err)
	}
	if !strings.Contains(gotCommand, `wt.exe -p "Arch Profile"`) {
		t.Fatalf("commandLine = %q, want to contain %q", gotCommand, `wt.exe -p "Arch Profile"`)
	}
}

func TestExecWindowsTerminalWithDeps_NoProfileUsesExecutableRun(t *testing.T) {
	t.Parallel()

	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, nil
		},
	}

	gotCommand := ""
	deps := execWindowsTerminalDeps{
		readParseWTConfig: func() (wtutils.Config, error) {
			return wtutils.Config{}, errors.New("no config")
		},
		createProfileGUID: func(string) string { return "unused" },
		getenv:            func(string) string { return `C:\Users\user\AppData\Local` },
		mustExecutable: func() string {
			return `C:\Program Files\wsldl\arch.exe`
		},
		createProcessAndWait: func(commandLine string) (int, error) {
			gotCommand = commandLine
			return 0, nil
		},
		allocConsole: func() {},
	}

	err := execWindowsTerminalWithDeps(reg, "arch", deps)
	if err != nil {
		t.Fatalf("execWindowsTerminalWithDeps returned error: %v", err)
	}
	if !strings.Contains(gotCommand, `"C:\Program Files\wsldl\arch.exe" run`) {
		t.Fatalf("commandLine = %q, want to contain executable run path", gotCommand)
	}
}

func TestExecWindowsTerminalWithDeps_FallbackToCaseInsensitiveProfileName(t *testing.T) {
	t.Parallel()

	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{DistributionName: "Arch"}, nil
		},
	}
	conf := wtutils.Config{}
	conf.Profiles.ProfileList = []wtutils.Profile{
		{GUID: "{other-guid}", Name: "ARCH"},
	}

	gotCommand := ""
	deps := execWindowsTerminalDeps{
		readParseWTConfig: func() (wtutils.Config, error) {
			return conf, nil
		},
		createProfileGUID: func(name string) string {
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return "guid-1234"
		},
		getenv:         func(string) string { return `C:\Users\user\AppData\Local` },
		mustExecutable: func() string { t.Fatal("mustExecutable should not be called"); return "" },
		createProcessAndWait: func(commandLine string) (int, error) {
			gotCommand = commandLine
			return 0, nil
		},
		allocConsole: func() {},
	}

	err := execWindowsTerminalWithDeps(reg, "arch", deps)
	if err != nil {
		t.Fatalf("execWindowsTerminalWithDeps returned error: %v", err)
	}
	if !strings.Contains(gotCommand, `wt.exe -p ARCH`) {
		t.Fatalf("commandLine = %q, want to contain %q", gotCommand, `wt.exe -p ARCH`)
	}
}

func TestExecWindowsTerminalWithDeps_ProcessErrorReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("spawn failed")
	allocCalled := 0
	deps := execWindowsTerminalDeps{
		readParseWTConfig:    func() (wtutils.Config, error) { return wtutils.Config{}, errors.New("no config") },
		createProfileGUID:    func(string) string { return "" },
		getenv:               func(string) string { return `C:\Users\user\AppData\Local` },
		mustExecutable:       func() string { return `C:\wsldl\arch.exe` },
		createProcessAndWait: func(string) (int, error) { return 0, wantErr },
		allocConsole: func() {
			allocCalled++
		},
	}

	err := execWindowsTerminalWithDeps(wsllib.MockWslReg{}, "arch", deps)
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
	if allocCalled != 1 {
		t.Fatalf("allocConsole call count = %d, want 1", allocCalled)
	}
}

func TestExecWindowsTerminalWithDeps_NonZeroExitReturnsExitCodeError(t *testing.T) {
	t.Parallel()

	deps := execWindowsTerminalDeps{
		readParseWTConfig:    func() (wtutils.Config, error) { return wtutils.Config{}, errors.New("no config") },
		createProfileGUID:    func(string) string { return "" },
		getenv:               func(string) string { return `C:\Users\user\AppData\Local` },
		mustExecutable:       func() string { return `C:\wsldl\arch.exe` },
		createProcessAndWait: func(string) (int, error) { return 23, nil },
		allocConsole:         func() {},
	}

	err := execWindowsTerminalWithDeps(wsllib.MockWslReg{}, "arch", deps)
	ee := assertExitCodeError(t, err)
	if ee.Code != 23 {
		t.Fatalf("exit code = %d, want %d", ee.Code, 23)
	}
	if ee.Pause {
		t.Fatal("pause = true, want false")
	}
}
