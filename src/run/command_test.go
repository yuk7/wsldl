package run

import (
	"errors"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/fileutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestParseRunArgs_Nil_UsesCurrentDirRule(t *testing.T) {
	t.Parallel()

	opts := parseRunArgs(nil)

	wantInherit := !fileutil.IsCurrentDirSpecial()
	if opts.inheritPath != wantInherit {
		t.Fatalf("inheritPath = %v, want %v", opts.inheritPath, wantInherit)
	}
	if len(opts.commandArgs) != 0 {
		t.Fatalf("command args len = %d, want 0", len(opts.commandArgs))
	}
}

func TestParseRunNoArgs_NonEmpty_ReturnsInvalid(t *testing.T) {
	t.Parallel()

	if _, err := parseRunNoArgs([]string{"extra"}); !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("err = %v, want %v", err, os.ErrInvalid)
	}
}

func TestParseRunNoArgs_Empty_Succeeds(t *testing.T) {
	t.Parallel()

	if _, err := parseRunNoArgs(nil); err != nil {
		t.Fatalf("parseRunNoArgs returned error: %v", err)
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

func TestGetCommandWithDeps_HelpVisibility(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}
	cmd := GetCommandWithDeps(wsl)
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

func TestGetCommandPWithDeps_HelpVisibility(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}
	cmd := GetCommandPWithDeps(wsl)
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

func TestGetCommandWithDeps_RunDelegatesToExecute(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return true
		},
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if command != " echo" {
				t.Fatalf("command = %q, want %q", command, " echo")
			}
			if !inheritPath {
				t.Fatal("inheritPath = false, want true")
			}
			return 0, nil
		},
	}

	cmd := GetCommandWithDeps(wsl)
	if !cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = false, want true")
	}
	if err := cmd.Run("Arch", []string{"echo"}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
}

func TestGetCommandPWithDeps_RunDelegatesToExecuteP(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return true
		},
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if command != " ls -la" {
				t.Fatalf("command = %q, want %q", command, " ls -la")
			}
			if !inheritPath {
				t.Fatal("inheritPath = false, want true")
			}
			return 0, nil
		},
	}

	cmd := GetCommandPWithDeps(wsl)
	if !cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = false, want true")
	}
	if err := cmd.Run("Arch", []string{"ls", "-la"}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
}

func TestGetCommandWithNoArgsWithDeps_RunDelegatesToExecuteNoArgs(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return true
		},
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if command != "" {
				t.Fatalf("command = %q, want empty", command)
			}
			wantInherit := !fileutil.IsCurrentDirSpecial()
			if inheritPath != wantInherit {
				t.Fatalf("inheritPath = %v, want %v", inheritPath, wantInherit)
			}
			return 0, nil
		},
	}

	cmd := GetCommandWithNoArgsWithDeps(wsl, wsllib.MockWslReg{})
	if !cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = false, want true")
	}
	if err := cmd.Run("Arch", nil); err != nil {
		t.Fatalf("Run returned error: %v", err)
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

func TestExecuteP_WithPathTranslation_NonWindows_ReturnsDisplayError(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("non-windows stub behavior only")
	}

	err := executeP(wsllib.MockWslLib{}, "Arch", []string{`C:\Users\user\file.txt`})
	de := assertDisplayError(t, err)
	if !strings.Contains(de.Error(), "only available on Windows") {
		t.Fatalf("error = %q, want to contain %q", de.Error(), "only available on Windows")
	}
}

func TestExecutePWithOptionsWithExecRead_PathTranslationError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("wslpath failed")
	execReadCalled := 0
	wsl := wsllib.MockWslLib{
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			t.Fatal("LaunchInteractive should not be called when path translation fails")
			return 0, nil
		},
	}

	err := executePWithOptionsWithExecRead(
		wsl,
		"Arch",
		runPOptions{commandArgs: []string{`C:\Users\user\file.txt`}},
		func(_ wsllib.WslLib, name, command string) (string, uint32, error) {
			execReadCalled++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if command != "wslpath -u C:/Users/user/file.txt" {
				t.Fatalf("command = %q, want %q", command, "wslpath -u C:/Users/user/file.txt")
			}
			return "", 0, wantErr
		},
	)

	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
	if execReadCalled != 1 {
		t.Fatalf("execRead call count = %d, want 1", execReadCalled)
	}
}

func TestExecutePWithOptionsWithExecRead_PathTranslationNonZeroExit_ReturnsExitCodeError(t *testing.T) {
	t.Parallel()

	execReadCalled := 0
	wsl := wsllib.MockWslLib{
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			t.Fatal("LaunchInteractive should not be called when path translation exits non-zero")
			return 0, nil
		},
	}

	err := executePWithOptionsWithExecRead(
		wsl,
		"Arch",
		runPOptions{commandArgs: []string{`C:\Users\user\file.txt`}},
		func(wsllib.WslLib, string, string) (string, uint32, error) {
			execReadCalled++
			return "", 5, nil
		},
	)

	ee := assertExitCodeError(t, err)
	if ee.Code != 5 {
		t.Fatalf("exit code = %d, want %d", ee.Code, 5)
	}
	if ee.Pause {
		t.Fatal("pause = true, want false")
	}
	if execReadCalled != 1 {
		t.Fatalf("execRead call count = %d, want 1", execReadCalled)
	}
}

func TestExecutePWithOptionsWithExecRead_PathTranslationSuccess_UsesTranslatedArgs(t *testing.T) {
	t.Parallel()

	execReadCalled := 0
	wsl := wsllib.MockWslLib{
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if !strings.Contains(command, " /mnt/c/Users/user/file.txt") {
				t.Fatalf("command = %q, want translated path", command)
			}
			return 0, nil
		},
	}

	err := executePWithOptionsWithExecRead(
		wsl,
		"Arch",
		runPOptions{commandArgs: []string{`C:\Users\user\file.txt`}},
		func(_ wsllib.WslLib, name, command string) (string, uint32, error) {
			execReadCalled++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if command != "wslpath -u C:/Users/user/file.txt" {
				t.Fatalf("command = %q, want %q", command, "wslpath -u C:/Users/user/file.txt")
			}
			return "/mnt/c/Users/user/file.txt", 0, nil
		},
	)
	if err != nil {
		t.Fatalf("executePWithOptionsWithExecRead returned error: %v", err)
	}
	if execReadCalled != 1 {
		t.Fatalf("execRead call count = %d, want 1", execReadCalled)
	}
}

func TestExecuteNoArgs_WithExtraArg_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	err := executeNoArgs(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", []string{"extra"})
	assertDisplayError(t, err)
}

func TestExecuteNoArgsWithOptions_DefaultDeps_DelegatesToExecute(t *testing.T) {
	t.Parallel()

	called := 0
	wsl := wsllib.MockWslLib{
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if command != "" {
				t.Fatalf("command = %q, want empty", command)
			}
			wantInherit := !fileutil.IsCurrentDirSpecial()
			if inheritPath != wantInherit {
				t.Fatalf("inheritPath = %v, want %v", inheritPath, wantInherit)
			}
			return 0, nil
		},
	}

	err := executeNoArgsWithOptions(wsl, wsllib.MockWslReg{}, "Arch", runNoArgsOptions{})
	if err != nil {
		t.Fatalf("executeNoArgsWithOptions returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("LaunchInteractive call count = %d, want 1", called)
	}
}

func TestExecuteNoArgsWithOptionsAndDeps_RepairAccepted_ReturnsPauseExit(t *testing.T) {
	t.Parallel()

	repairCalled := false
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{
				BasePath:         "/missing",
				DistributionName: name,
			}, nil
		},
	}
	deps := runNoArgsDeps{
		mustExecutable: func() string { return "/tmp/Arch.exe" },
		stat: func(string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		},
		isInstalledFilesExist: func() bool { return true },
		readInput:             func() string { return "y" },
		repairRegistry: func(_ wsllib.WslReg, profile wsllib.Profile) error {
			repairCalled = true
			if profile.BasePath != "/missing" {
				t.Fatalf("BasePath = %q, want %q", profile.BasePath, "/missing")
			}
			return nil
		},
		isParentConsole: func() (bool, error) {
			t.Fatal("isParentConsole should not be called when repair is accepted")
			return false, nil
		},
		execute: func(wsllib.WslLib, string, []string) error {
			t.Fatal("execute should not be called when repair is accepted")
			return nil
		},
	}

	err := executeNoArgsWithOptionsAndDeps(wsllib.MockWslLib{}, reg, "Arch", runNoArgsOptions{}, deps)
	ee := assertExitCodeError(t, err)
	if ee.Code != 0 {
		t.Fatalf("exit code = %d, want 0", ee.Code)
	}
	if !ee.Pause {
		t.Fatal("pause = false, want true")
	}
	if !repairCalled {
		t.Fatal("repairRegistry was not called")
	}
}

func TestExecuteNoArgsWithOptionsAndDeps_RepairError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("repair failed")
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{
				BasePath:         "/missing",
				DistributionName: name,
			}, nil
		},
	}
	deps := newRunNoArgsDepsForTest()
	deps.stat = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	deps.isInstalledFilesExist = func() bool { return true }
	deps.readInput = func() string { return "y" }
	deps.repairRegistry = func(wsllib.WslReg, wsllib.Profile) error { return wantErr }

	err := executeNoArgsWithOptionsAndDeps(wsllib.MockWslLib{}, reg, "Arch", runNoArgsOptions{}, deps)
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
	if !de.Pause {
		t.Fatal("pause = false, want true")
	}
}

func TestExecuteNoArgsWithOptionsAndDeps_RepairDeclined_FallsBackToExecute(t *testing.T) {
	t.Parallel()

	repairCalled := false
	executeCalled := 0
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{
				BasePath:         "/missing",
				DistributionName: name,
			}, nil
		},
	}
	deps := runNoArgsDeps{
		mustExecutable: func() string { return "/tmp/Arch.exe" },
		stat: func(string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		},
		isInstalledFilesExist: func() bool { return true },
		readInput:             func() string { return "n" },
		repairRegistry: func(wsllib.WslReg, wsllib.Profile) error {
			repairCalled = true
			return nil
		},
		isParentConsole: func() (bool, error) { return true, nil },
		execute: func(_ wsllib.WslLib, name string, args []string) error {
			executeCalled++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if args != nil {
				t.Fatalf("args = %v, want nil", args)
			}
			return nil
		},
	}

	err := executeNoArgsWithOptionsAndDeps(wsllib.MockWslLib{}, reg, "Arch", runNoArgsOptions{}, deps)
	if err != nil {
		t.Fatalf("executeNoArgsWithOptionsAndDeps returned error: %v", err)
	}
	if repairCalled {
		t.Fatal("repairRegistry should not be called when answer is n")
	}
	if executeCalled != 1 {
		t.Fatalf("execute call count = %d, want 1", executeCalled)
	}
}

func TestExecuteNoArgsWithOptionsAndDeps_IsParentConsoleError_FallsBackToExecute(t *testing.T) {
	t.Parallel()

	executeCalled := 0
	deps := newRunNoArgsDepsForTest()
	deps.isParentConsole = func() (bool, error) { return false, errors.New("detect failed") }
	deps.execute = func(_ wsllib.WslLib, name string, args []string) error {
		executeCalled++
		if name != "Arch" {
			t.Fatalf("name = %q, want %q", name, "Arch")
		}
		if args != nil {
			t.Fatalf("args = %v, want nil", args)
		}
		return nil
	}

	err := executeNoArgsWithOptionsAndDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", runNoArgsOptions{}, deps)
	if err != nil {
		t.Fatalf("executeNoArgsWithOptionsAndDeps returned error: %v", err)
	}
	if executeCalled != 1 {
		t.Fatalf("execute call count = %d, want 1", executeCalled)
	}
}

func TestExecuteNoArgsWithOptionsAndDeps_NonConsoleWT_UsesWindowsTerminal(t *testing.T) {
	t.Parallel()

	freeCalled := 0
	execWTCalled := 0
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{WsldlTerm: wsllib.FlagWsldlTermWT}, nil
		},
	}
	deps := newRunNoArgsDepsForTest()
	deps.isParentConsole = func() (bool, error) { return false, nil }
	deps.freeConsole = func() error {
		freeCalled++
		return nil
	}
	deps.execWindowsTerminal = func(_ wsllib.WslReg, name string) error {
		execWTCalled++
		if name != "Arch" {
			t.Fatalf("name = %q, want %q", name, "Arch")
		}
		return nil
	}
	deps.execute = func(wsllib.WslLib, string, []string) error {
		t.Fatal("execute should not be called in WT branch")
		return nil
	}

	err := executeNoArgsWithOptionsAndDeps(wsllib.MockWslLib{}, reg, "Arch", runNoArgsOptions{}, deps)
	if err != nil {
		t.Fatalf("executeNoArgsWithOptionsAndDeps returned error: %v", err)
	}
	if freeCalled != 1 {
		t.Fatalf("freeConsole call count = %d, want 1", freeCalled)
	}
	if execWTCalled != 1 {
		t.Fatalf("execWindowsTerminal call count = %d, want 1", execWTCalled)
	}
}

func TestExecuteNoArgsWithOptionsAndDeps_NonConsoleFluteProcessError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("spawn failed")
	allocCalled := 0
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{WsldlTerm: wsllib.FlagWsldlTermFlute}, nil
		},
	}
	deps := newRunNoArgsDepsForTest()
	deps.isParentConsole = func() (bool, error) { return false, nil }
	deps.getenv = func(key string) string {
		if key != "LOCALAPPDATA" {
			t.Fatalf("env key = %q, want %q", key, "LOCALAPPDATA")
		}
		return `C:\Users\user\AppData\Local`
	}
	deps.createProcessAndWait = func(commandLine string) (int, error) {
		if !strings.Contains(commandLine, "flute.exe") {
			t.Fatalf("commandLine = %q, want to contain flute.exe", commandLine)
		}
		return 0, wantErr
	}
	deps.allocConsole = func() {
		allocCalled++
	}

	err := executeNoArgsWithOptionsAndDeps(wsllib.MockWslLib{}, reg, "Arch", runNoArgsOptions{}, deps)
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
	if allocCalled != 1 {
		t.Fatalf("allocConsole call count = %d, want 1", allocCalled)
	}
}

func TestExecuteNoArgsWithOptionsAndDeps_NonConsoleFluteNonZeroExit_ReturnsExitCodeError(t *testing.T) {
	t.Parallel()

	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{WsldlTerm: wsllib.FlagWsldlTermFlute}, nil
		},
	}
	deps := newRunNoArgsDepsForTest()
	deps.isParentConsole = func() (bool, error) { return false, nil }
	deps.createProcessAndWait = func(string) (int, error) { return 9, nil }

	err := executeNoArgsWithOptionsAndDeps(wsllib.MockWslLib{}, reg, "Arch", runNoArgsOptions{}, deps)
	ee := assertExitCodeError(t, err)
	if ee.Code != 9 {
		t.Fatalf("exit code = %d, want %d", ee.Code, 9)
	}
	if ee.Pause {
		t.Fatal("pause = true, want false")
	}
}

func TestExecuteNoArgsWithOptionsAndDeps_NonConsoleFluteSuccess_ReturnsNil(t *testing.T) {
	t.Parallel()

	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{WsldlTerm: wsllib.FlagWsldlTermFlute}, nil
		},
	}
	freeCalled := 0
	deps := newRunNoArgsDepsForTest()
	deps.isParentConsole = func() (bool, error) { return false, nil }
	deps.freeConsole = func() error {
		freeCalled++
		return nil
	}
	deps.createProcessAndWait = func(string) (int, error) { return 0, nil }

	err := executeNoArgsWithOptionsAndDeps(wsllib.MockWslLib{}, reg, "Arch", runNoArgsOptions{}, deps)
	if err != nil {
		t.Fatalf("executeNoArgsWithOptionsAndDeps returned error: %v", err)
	}
	if freeCalled != 1 {
		t.Fatalf("freeConsole call count = %d, want 1", freeCalled)
	}
}

func TestExecuteNoArgsWithOptionsAndDeps_NonConsoleDefault_UsesRegistryName(t *testing.T) {
	t.Parallel()

	titled := ""
	executed := ""
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{
				DistributionName: "ArchLinux",
				WsldlTerm:        wsllib.FlagWsldlTermDefault,
			}, nil
		},
	}
	deps := newRunNoArgsDepsForTest()
	deps.isParentConsole = func() (bool, error) { return false, nil }
	deps.setConsoleTitle = func(title string) { titled = title }
	deps.execute = func(_ wsllib.WslLib, name string, args []string) error {
		executed = name
		if args != nil {
			t.Fatalf("args = %v, want nil", args)
		}
		return nil
	}

	err := executeNoArgsWithOptionsAndDeps(wsllib.MockWslLib{}, reg, "arch", runNoArgsOptions{}, deps)
	if err != nil {
		t.Fatalf("executeNoArgsWithOptionsAndDeps returned error: %v", err)
	}
	if titled != "ArchLinux" {
		t.Fatalf("title = %q, want %q", titled, "ArchLinux")
	}
	if executed != "ArchLinux" {
		t.Fatalf("executed name = %q, want %q", executed, "ArchLinux")
	}
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
	if cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = true, want false in unit-test dependencies")
	}
}

func TestGetCommandP_DefaultDeps_BasicShape(t *testing.T) {
	t.Parallel()

	cmd := GetCommandP()
	if len(cmd.Names) == 0 {
		t.Fatal("Names should not be empty")
	}
	if cmd.Run == nil {
		t.Fatal("Run is nil")
	}
	if cmd.Visible == nil {
		t.Fatal("Visible is nil")
	}
	if cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = true, want false in unit-test dependencies")
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
	if cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = true, want false in unit-test dependencies")
	}
}

func TestDefaultRunNoArgsDeps_ReadInput(t *testing.T) {
	deps := defaultRunNoArgsDeps()
	if deps.readInput == nil {
		t.Fatal("readInput is nil")
	}

	origStdin := os.Stdin
	inR, inW, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdin pipe failed: %v", err)
	}
	os.Stdin = inR
	t.Cleanup(func() {
		os.Stdin = origStdin
		_ = inR.Close()
		_ = inW.Close()
	})

	if _, err := io.WriteString(inW, "repair\n"); err != nil {
		t.Fatalf("write stdin failed: %v", err)
	}
	if err := inW.Close(); err != nil {
		t.Fatalf("close stdin writer failed: %v", err)
	}

	got := deps.readInput()
	if got != "repair" {
		t.Fatalf("readInput = %q, want %q", got, "repair")
	}
}

func newRunNoArgsDepsForTest() runNoArgsDeps {
	return runNoArgsDeps{
		mustExecutable:        func() string { return "/tmp/Arch.exe" },
		stat:                  func(string) (os.FileInfo, error) { return nil, nil },
		isInstalledFilesExist: func() bool { return false },
		readInput:             func() string { return "n" },
		repairRegistry:        func(wsllib.WslReg, wsllib.Profile) error { return nil },
		isParentConsole:       func() (bool, error) { return true, nil },
		freeConsole:           func() error { return nil },
		allocConsole:          func() {},
		setConsoleTitle:       func(string) {},
		execWindowsTerminal:   func(wsllib.WslReg, string) error { return nil },
		getenv:                func(string) string { return `C:\Users\user\AppData\Local` },
		createProcessAndWait:  func(string) (int, error) { return 0, nil },
		execute:               func(wsllib.WslLib, string, []string) error { return nil },
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
