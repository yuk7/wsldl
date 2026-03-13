package install

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/preset"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestParseArgs_Nil_IsAutoFromNoArg(t *testing.T) {
	t.Parallel()

	parsed, err := parseArgs(nil)
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if parsed.mode != installModeAuto {
		t.Fatalf("mode = %v, want %v", parsed.mode, installModeAuto)
	}
	if !parsed.fromNoArgCall {
		t.Fatal("fromNoArgCall = false, want true")
	}
}

func TestParseArgs_EmptySlice_IsAutoButNotNoArgCall(t *testing.T) {
	t.Parallel()

	parsed, err := parseArgs([]string{})
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if parsed.mode != installModeAuto {
		t.Fatalf("mode = %v, want %v", parsed.mode, installModeAuto)
	}
	if parsed.fromNoArgCall {
		t.Fatal("fromNoArgCall = true, want false")
	}
}

func TestParseArgs_RootFlag_IsRootMode(t *testing.T) {
	t.Parallel()

	parsed, err := parseArgs([]string{"--root"})
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if parsed.mode != installModeRoot {
		t.Fatalf("mode = %v, want %v", parsed.mode, installModeRoot)
	}
	if parsed.fromNoArgCall {
		t.Fatal("fromNoArgCall = true, want false")
	}
}

func TestParseArgs_CustomPath_IsPathMode(t *testing.T) {
	t.Parallel()

	parsed, err := parseArgs([]string{"rootfs.tar"})
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if parsed.mode != installModePath {
		t.Fatalf("mode = %v, want %v", parsed.mode, installModePath)
	}
	if parsed.inputPath != "rootfs.tar" {
		t.Fatalf("inputPath = %q, want %q", parsed.inputPath, "rootfs.tar")
	}
}

func TestParseArgs_TooManyArgs_ReturnsInvalid(t *testing.T) {
	t.Parallel()

	if _, err := parseArgs([]string{"a", "b"}); !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("err = %v, want %v", err, os.ErrInvalid)
	}
}

func TestResolveOptions_PathMode_UsesInputPath(t *testing.T) {
	t.Parallel()

	opts := resolveOptions(installArgs{
		mode:          installModePath,
		inputPath:     "custom.tar.gz",
		fromNoArgCall: false,
	})
	if opts.rootPath != "custom.tar.gz" {
		t.Fatalf("rootPath = %q, want %q", opts.rootPath, "custom.tar.gz")
	}
	if opts.showProgress {
		t.Fatal("showProgress = true, want false")
	}
	if opts.pauseAfterRun {
		t.Fatal("pauseAfterRun = true, want false")
	}
}

func TestResolveOptions_AutoMode_FallbacksWhenDetectFails(t *testing.T) {
	orig := detectRootfsFilesFunc
	detectRootfsFilesFunc = func() (string, error) {
		return "", errors.New("detect failed")
	}
	t.Cleanup(func() {
		detectRootfsFilesFunc = orig
	})

	opts := resolveOptions(installArgs{
		mode: installModeAuto,
	})
	if opts.rootPath != "rootfs.tar.gz" {
		t.Fatalf("rootPath = %q, want %q", opts.rootPath, "rootfs.tar.gz")
	}
}

func TestResolveOptions_AutoMode_UsesDetectedPathWhenAvailable(t *testing.T) {
	orig := detectRootfsFilesFunc
	detectRootfsFilesFunc = func() (string, error) {
		return "install.tar", nil
	}
	t.Cleanup(func() {
		detectRootfsFilesFunc = orig
	})

	opts := resolveOptions(installArgs{
		mode: installModeAuto,
	})
	if opts.rootPath != "install.tar" {
		t.Fatalf("rootPath = %q, want %q", opts.rootPath, "install.tar")
	}
}

func TestResolveOptions_AutoMode_PresetInstallFileOverridesDetectedPath(t *testing.T) {
	origDetect := detectRootfsFilesFunc
	detectRootfsFilesFunc = func() (string, error) {
		return "install.tar", nil
	}
	origReadPreset := readParsePresetFunc
	readParsePresetFunc = func() (preset.Preset, error) {
		return preset.Preset{
			InstallFile:       "override.tar.gz",
			InstallFileSha256: "abc123",
		}, nil
	}
	t.Cleanup(func() {
		detectRootfsFilesFunc = origDetect
		readParsePresetFunc = origReadPreset
	})

	opts := resolveOptions(installArgs{
		mode: installModeAuto,
	})
	if opts.rootPath != "override.tar.gz" {
		t.Fatalf("rootPath = %q, want %q", opts.rootPath, "override.tar.gz")
	}
	if opts.rootFileSHA256 != "abc123" {
		t.Fatalf("rootFileSHA256 = %q, want %q", opts.rootFileSHA256, "abc123")
	}
}

func TestExecute_AlreadyRegistered_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return true
		},
	}
	err := execute(wsl, wsllib.MockWslReg{}, "Arch", nil)
	de := assertDisplayError(t, err)
	if !errors.Is(de, os.ErrInvalid) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), os.ErrInvalid)
	}
}

func TestExecute_InvalidArgs_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	err := execute(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", []string{"a", "b"})
	de := assertDisplayError(t, err)
	if !errors.Is(de, os.ErrInvalid) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), os.ErrInvalid)
	}
}

func TestExecute_PathMode_Success(t *testing.T) {
	t.Parallel()

	called := 0
	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
		RegisterDistributionFunc: func(name, rootPath string) error {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if rootPath != "custom.tar.gz" {
				t.Fatalf("rootPath = %q, want %q", rootPath, "custom.tar.gz")
			}
			return nil
		},
	}

	err := execute(wsl, wsllib.MockWslReg{}, "Arch", []string{"custom.tar.gz"})
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("RegisterDistribution call count = %d, want 1", called)
	}
}

func TestExecuteWithOptions_InstallError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("register failed")
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			return wantErr
		},
	}

	err := executeWithOptions(wsl, wsllib.MockWslReg{}, "Arch", installOptions{
		rootPath:      "rootfs.tar",
		showProgress:  false,
		pauseAfterRun: false,
		presetVersion: 0,
	})
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecuteWithOptionsAndDeps_RepairError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("repair failed")
	deps := installCommandDeps{
		isInstalledFilesExist: func() bool { return true },
		readInput:             func() string { return "y" },
		repairRegistry: func(wsllib.WslReg, string) error {
			return wantErr
		},
		install: func(context.Context, wsllib.WslLib, wsllib.WslReg, string, string, string, bool) error {
			t.Fatal("install should not be called when repair is chosen")
			return nil
		},
		setVersion: func(string, int) error {
			t.Fatal("setVersion should not be called when repair is chosen")
			return nil
		},
		waitForEnter: func() {},
	}

	err := executeWithOptionsAndDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", installOptions{
		pauseAfterRun: true,
		showProgress:  false,
	}, deps)
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecuteWithOptionsAndDeps_RepairAccepted_ReturnsEarly(t *testing.T) {
	t.Parallel()

	repairCalled := false
	deps := installCommandDeps{
		isInstalledFilesExist: func() bool { return true },
		readInput:             func() string { return "y" },
		repairRegistry: func(reg wsllib.WslReg, name string) error {
			repairCalled = true
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return nil
		},
		install: func(context.Context, wsllib.WslLib, wsllib.WslReg, string, string, string, bool) error {
			t.Fatal("install should not be called when repair is accepted")
			return nil
		},
		setVersion: func(string, int) error {
			t.Fatal("setVersion should not be called when repair is accepted")
			return nil
		},
		waitForEnter: func() {
			t.Fatal("waitForEnter should not be called when repair is accepted")
		},
	}

	err := executeWithOptionsAndDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", installOptions{
		pauseAfterRun: true,
		showProgress:  false,
	}, deps)
	if err != nil {
		t.Fatalf("executeWithOptionsAndDeps returned error: %v", err)
	}
	if !repairCalled {
		t.Fatal("repairRegistry was not called")
	}
}

func TestExecuteWithOptionsAndDeps_RepairDeclined_ContinuesInstall(t *testing.T) {
	t.Parallel()

	installCalled := 0
	deps := installCommandDeps{
		isInstalledFilesExist: func() bool { return true },
		readInput:             func() string { return "n" },
		repairRegistry: func(wsllib.WslReg, string) error {
			t.Fatal("repairRegistry should not be called when answer is n")
			return nil
		},
		install: func(ctx context.Context, wsl wsllib.WslLib, reg wsllib.WslReg, name, rootPath, sha256Sum string, showProgress bool) error {
			installCalled++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if rootPath != "rootfs.tar.gz" {
				t.Fatalf("rootPath = %q, want %q", rootPath, "rootfs.tar.gz")
			}
			if sha256Sum != "" {
				t.Fatalf("sha256Sum = %q, want empty", sha256Sum)
			}
			if showProgress {
				t.Fatal("showProgress = true, want false")
			}
			return nil
		},
		setVersion: func(string, int) error { return nil },
		waitForEnter: func() {
			t.Fatal("waitForEnter should not be called when pauseAfterRun=false")
		},
	}

	err := executeWithOptionsAndDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", installOptions{
		rootPath:      "rootfs.tar.gz",
		showProgress:  false,
		pauseAfterRun: false,
		presetVersion: 0,
	}, deps)
	if err != nil {
		t.Fatalf("executeWithOptionsAndDeps returned error: %v", err)
	}
	if installCalled != 1 {
		t.Fatalf("install call count = %d, want 1", installCalled)
	}
}

func TestExecuteWithOptionsAndDeps_SetVersionError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("set-version failed")
	setVersionCalled := 0
	deps := installCommandDeps{
		isInstalledFilesExist: func() bool { return false },
		readInput:             func() string { return "" },
		repairRegistry: func(wsllib.WslReg, string) error {
			t.Fatal("repairRegistry should not be called")
			return nil
		},
		install: func(context.Context, wsllib.WslLib, wsllib.WslReg, string, string, string, bool) error {
			return nil
		},
		setVersion: func(name string, version int) error {
			setVersionCalled++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if version != 2 {
				t.Fatalf("version = %d, want %d", version, 2)
			}
			return wantErr
		},
		waitForEnter: func() {},
	}

	err := executeWithOptionsAndDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", installOptions{
		rootPath:      "rootfs.tar.gz",
		showProgress:  false,
		pauseAfterRun: false,
		presetVersion: 2,
	}, deps)
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
	if setVersionCalled != 1 {
		t.Fatalf("setVersion call count = %d, want 1", setVersionCalled)
	}
}

func TestExecuteWithOptionsAndDeps_SetVersionSuccess(t *testing.T) {
	t.Parallel()

	setVersionCalled := 0
	deps := installCommandDeps{
		isInstalledFilesExist: func() bool { return false },
		readInput:             func() string { return "" },
		repairRegistry:        func(wsllib.WslReg, string) error { return nil },
		install: func(context.Context, wsllib.WslLib, wsllib.WslReg, string, string, string, bool) error {
			return nil
		},
		setVersion: func(string, int) error {
			setVersionCalled++
			return nil
		},
		waitForEnter: func() {},
	}

	err := executeWithOptionsAndDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", installOptions{
		rootPath:      "rootfs.tar.gz",
		showProgress:  false,
		pauseAfterRun: false,
		presetVersion: 1,
	}, deps)
	if err != nil {
		t.Fatalf("executeWithOptionsAndDeps returned error: %v", err)
	}
	if setVersionCalled != 1 {
		t.Fatalf("setVersion call count = %d, want 1", setVersionCalled)
	}
}

func TestExecuteWithOptionsAndDeps_PauseAfterRun_WaitsForEnter(t *testing.T) {
	t.Parallel()

	installCalled := false
	waitCalled := 0
	deps := installCommandDeps{
		isInstalledFilesExist: func() bool { return false },
		readInput:             func() string { return "" },
		repairRegistry: func(wsllib.WslReg, string) error {
			t.Fatal("repairRegistry should not be called")
			return nil
		},
		install: func(ctx context.Context, wsl wsllib.WslLib, reg wsllib.WslReg, name, rootPath, sha256Sum string, showProgress bool) error {
			installCalled = true
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if rootPath != "rootfs.tar.gz" {
				t.Fatalf("rootPath = %q, want %q", rootPath, "rootfs.tar.gz")
			}
			if sha256Sum != "" {
				t.Fatalf("sha256Sum = %q, want empty", sha256Sum)
			}
			if showProgress {
				t.Fatal("showProgress = true, want false")
			}
			return nil
		},
		setVersion: func(string, int) error {
			t.Fatal("setVersion should not be called when presetVersion=0")
			return nil
		},
		waitForEnter: func() {
			waitCalled++
		},
	}

	err := executeWithOptionsAndDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", installOptions{
		rootPath:      "rootfs.tar.gz",
		showProgress:  false,
		pauseAfterRun: true,
		presetVersion: 0,
	}, deps)
	if err != nil {
		t.Fatalf("executeWithOptionsAndDeps returned error: %v", err)
	}
	if !installCalled {
		t.Fatal("install was not called")
	}
	if waitCalled != 1 {
		t.Fatalf("waitForEnter call count = %d, want 1", waitCalled)
	}
}

func TestExecuteWithOptionsAndDeps_ShowProgressTrue_PrintsCompletion(t *testing.T) {
	t.Parallel()

	deps := installCommandDeps{
		isInstalledFilesExist: func() bool { return false },
		readInput:             func() string { return "" },
		repairRegistry:        func(wsllib.WslReg, string) error { return nil },
		install: func(context.Context, wsllib.WslLib, wsllib.WslReg, string, string, string, bool) error {
			return nil
		},
		setVersion: func(string, int) error { return nil },
		waitForEnter: func() {
			t.Fatal("waitForEnter should not be called when pauseAfterRun=false")
		},
	}

	err := executeWithOptionsAndDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", installOptions{
		rootPath:      "rootfs.tar.gz",
		showProgress:  true,
		pauseAfterRun: false,
		presetVersion: 0,
	}, deps)
	if err != nil {
		t.Fatalf("executeWithOptionsAndDeps returned error: %v", err)
	}
}

func TestGetCommandWithNoArgsWithDeps_HelpVisibility(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}
	cmd := GetCommandWithNoArgsWithDeps(wsl, wsllib.MockWslReg{})
	if cmd.Visible == nil {
		t.Fatal("Visible is nil")
	}
	if !cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = false, want true")
	}
	if cmd.HelpText == nil {
		t.Fatal("HelpText is nil")
	}
	if got := cmd.HelpText(); got == "" {
		t.Fatal("HelpText should not be empty")
	}

	err := cmd.Run("Arch", []string{"a", "b"})
	assertDisplayError(t, err)
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
	if !cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = false, want true")
	}
	if cmd.HelpText == nil {
		t.Fatal("HelpText is nil")
	}
	if got := cmd.HelpText(); got == "" {
		t.Fatal("HelpText should not be empty")
	}

	err := cmd.Run("Arch", []string{"a", "b"})
	assertDisplayError(t, err)
}

func TestGetCommand_DefaultDeps_BasicShape(t *testing.T) {
	t.Parallel()

	cmd := GetCommand()
	if len(cmd.Names) == 0 {
		t.Fatal("Names should not be empty")
	}
	if cmd.Run == nil {
		t.Fatal("Run is nil")
	}
	if cmd.Visible == nil {
		t.Fatal("Visible is nil")
	}
	if !cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = false, want true with unit-test mock deps")
	}
}

func TestGetCommandWithNoArgs_DefaultDeps_BasicShape(t *testing.T) {
	t.Parallel()

	cmd := GetCommandWithNoArgs()
	if cmd.Run == nil {
		t.Fatal("Run is nil")
	}
	if cmd.Visible == nil {
		t.Fatal("Visible is nil")
	}
	if !cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = false, want true with unit-test mock deps")
	}
}

func TestDefaultInstallCommandDeps_InstallReturnsContextError(t *testing.T) {
	t.Parallel()

	deps := defaultInstallCommandDeps()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := deps.install(ctx, wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", "rootfs.tar", "", false)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("install error = %v, want %v", err, context.Canceled)
	}
}

func TestDefaultInstallCommandDeps_SetVersionReturnsErrorForMissingDistro(t *testing.T) {
	t.Parallel()

	deps := defaultInstallCommandDeps()
	err := deps.setVersion("wsldl-test-missing-distro-0f4f74a4", 2)
	if err == nil {
		t.Fatal("setVersion succeeded unexpectedly")
	}
}

func TestDefaultInstallCommandDeps_ReadInputAndWaitForEnter(t *testing.T) {
	deps := defaultInstallCommandDeps()
	if deps.readInput == nil {
		t.Fatal("readInput is nil")
	}
	if deps.waitForEnter == nil {
		t.Fatal("waitForEnter is nil")
	}

	origStdin := os.Stdin
	origStdout := os.Stdout
	inR, inW, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdin pipe failed: %v", err)
	}
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe failed: %v", err)
	}
	os.Stdin = inR
	os.Stdout = outW
	t.Cleanup(func() {
		os.Stdin = origStdin
		os.Stdout = origStdout
		_ = inR.Close()
		_ = inW.Close()
		_ = outR.Close()
		_ = outW.Close()
	})

	if _, err := io.WriteString(inW, "yes\n\n"); err != nil {
		t.Fatalf("write stdin data failed: %v", err)
	}
	if err := inW.Close(); err != nil {
		t.Fatalf("close stdin writer failed: %v", err)
	}

	got := deps.readInput()
	if got != "yes" {
		t.Fatalf("readInput = %q, want %q", got, "yes")
	}

	deps.waitForEnter()
	if err := outW.Close(); err != nil {
		t.Fatalf("close stdout writer failed: %v", err)
	}
	output, err := io.ReadAll(outR)
	if err != nil {
		t.Fatalf("read stdout failed: %v", err)
	}
	if !strings.Contains(string(output), "Press enter to continue...") {
		t.Fatalf("stdout = %q, want prompt text", string(output))
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
