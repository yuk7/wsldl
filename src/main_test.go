package main

import (
	"errors"
	"io"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/wsllib"
)

type exitSignal struct{}

func TestRunMain_SubCommandRun(t *testing.T) {
	const exePath = "/tmp/Arch.exe"

	origExit := exitFunc
	origExecutable := executableFunc
	defer func() {
		exitFunc = origExit
		executableFunc = origExecutable
	}()
	executableFunc = func() (string, error) { return exePath, nil }
	exitFunc = func(bool, int) {
		t.Fatal("unexpected exit")
	}

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

	runMain(
		wsllib.Dependencies{Wsl: wsl, Reg: wsllib.MockWslReg{}},
		[]string{"Arch.exe", "run", "echo", "hello world"},
		exePath,
	)

	if called != 1 {
		t.Fatalf("LaunchInteractive call count = %d, want 1", called)
	}
}

func TestRunMain_NoArgsNotRegistered_InstallPath(t *testing.T) {
	const exePath = "/tmp/Arch.exe"

	origExit := exitFunc
	origExecutable := executableFunc
	defer func() {
		exitFunc = origExit
		executableFunc = origExecutable
	}()
	executableFunc = func() (string, error) { return exePath, nil }

	wantErr := errors.New("register failed")
	registerCalled := 0
	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
		RegisterDistributionFunc: func(name, rootPath string) error {
			registerCalled++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if rootPath == "" {
				t.Fatal("rootPath should not be empty")
			}
			return wantErr
		},
	}

	assertExit(t, true, 1, func() {
		runMain(
			wsllib.Dependencies{Wsl: wsl, Reg: wsllib.MockWslReg{}},
			[]string{"Arch.exe"},
			exePath,
		)
	})

	if registerCalled != 1 {
		t.Fatalf("RegisterDistribution call count = %d, want 1", registerCalled)
	}
}

func TestRunMain_NoArgsRegistered_RunPath(t *testing.T) {
	const exePath = "/tmp/Arch.exe"

	origExit := exitFunc
	origExecutable := executableFunc
	defer func() {
		exitFunc = origExit
		executableFunc = origExecutable
	}()
	executableFunc = func() (string, error) { return exePath, nil }
	exitFunc = func(bool, int) {
		t.Fatal("unexpected exit")
	}

	called := 0
	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return true
		},
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if command != "" {
				t.Fatalf("command = %q, want empty string", command)
			}
			return 0, nil
		},
	}

	runMain(
		wsllib.Dependencies{Wsl: wsl, Reg: wsllib.MockWslReg{}},
		[]string{"Arch.exe"},
		exePath,
	)

	if called != 1 {
		t.Fatalf("LaunchInteractive call count = %d, want 1", called)
	}
}

func TestRunMain_InvalidSubcommand_ShowsHintAndExit(t *testing.T) {
	const exePath = "/tmp/Arch.exe"

	origExit := exitFunc
	origExecutable := executableFunc
	origStderr := os.Stderr
	defer func() {
		exitFunc = origExit
		executableFunc = origExecutable
		os.Stderr = origStderr
	}()
	executableFunc = func() (string, error) { return exePath, nil }

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stderr = w

	assertExit(t, false, 1, func() {
		runMain(
			wsllib.Dependencies{Wsl: wsllib.MockWslLib{}, Reg: wsllib.MockWslReg{}},
			[]string{"Arch.exe", "unknown"},
			exePath,
		)
	})

	if err := w.Close(); err != nil {
		t.Fatalf("stderr writer close failed: %v", err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("stderr read failed: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("stderr reader close failed: %v", err)
	}
	stderr := string(data)

	if !strings.Contains(stderr, "Your command may be incorrect.") {
		t.Fatalf("stderr = %q, want to contain %q", stderr, "Your command may be incorrect.")
	}
	if !strings.Contains(stderr, "`Arch.exe help`") {
		t.Fatalf("stderr = %q, want to contain %q", stderr, "`Arch.exe help`")
	}
}

func TestRunMain_HelpSubcommand_UsesShowHelpFromCommands(t *testing.T) {
	const exePath = "/tmp/Arch.exe"

	origExit := exitFunc
	origExecutable := executableFunc
	origShowHelp := showHelpFromCommandsFunc
	defer func() {
		exitFunc = origExit
		executableFunc = origExecutable
		showHelpFromCommandsFunc = origShowHelp
	}()
	executableFunc = func() (string, error) { return exePath, nil }
	exitFunc = func(bool, int) {
		t.Fatal("unexpected exit")
	}

	called := 0
	showHelpFromCommandsFunc = func(commands []cmdline.Command, distroName string, args []string) {
		called++
		if distroName != "Arch" {
			t.Fatalf("distroName = %q, want %q", distroName, "Arch")
		}
		if len(args) != 1 || args[0] != "run" {
			t.Fatalf("args = %v, want %v", args, []string{"run"})
		}
		if len(commands) == 0 {
			t.Fatal("commands should not be empty")
		}
	}

	runMain(
		wsllib.Dependencies{Wsl: wsllib.MockWslLib{}, Reg: wsllib.MockWslReg{}},
		[]string{"Arch.exe", "help", "run"},
		exePath,
	)

	if called != 1 {
		t.Fatalf("ShowHelpFromCommands call count = %d, want 1", called)
	}
}

func TestRunMain_ExitCodeErrorPath_ExitsWithCommandCode(t *testing.T) {
	const exePath = "/tmp/Arch.exe"

	origExit := exitFunc
	origExecutable := executableFunc
	defer func() {
		exitFunc = origExit
		executableFunc = origExecutable
	}()
	executableFunc = func() (string, error) { return exePath, nil }

	wsl := wsllib.MockWslLib{
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			return 23, nil
		},
	}

	assertExit(t, false, 23, func() {
		runMain(
			wsllib.Dependencies{Wsl: wsl, Reg: wsllib.MockWslReg{}},
			[]string{"Arch.exe", "run", "echo", "hello"},
			exePath,
		)
	})
}

func TestRunMain_UnexpectedError_FallsBackToDisplayError(t *testing.T) {
	const exePath = "/tmp/Arch.exe"

	origExit := exitFunc
	origExecutable := executableFunc
	origRunSubCommand := runSubCommandFunc
	defer func() {
		exitFunc = origExit
		executableFunc = origExecutable
		runSubCommandFunc = origRunSubCommand
	}()
	executableFunc = func() (string, error) { return exePath, nil }
	runSubCommandFunc = func([]cmdline.Command, string, []string) error {
		return errors.New("unexpected error")
	}

	stderr := captureStderr(t, func() {
		assertExit(t, false, 1, func() {
			runMain(
				wsllib.Dependencies{Wsl: wsllib.MockWslLib{}, Reg: wsllib.MockWslReg{}},
				[]string{"Arch.exe"},
				exePath,
			)
		})
	})
	if !strings.Contains(stderr, "unexpected error") {
		t.Fatalf("stderr = %q, want to contain %q", stderr, "unexpected error")
	}
}

func TestHandleDisplayError_ShowColorFalse_WritesToStderr(t *testing.T) {
	stderr := captureStderr(t, func() {
		assertExit(t, false, 1, func() {
			handleDisplayError(errors.New("plain failure"), true, false, false)
		})
	})

	if !strings.Contains(stderr, "ERR: plain failure") {
		t.Fatalf("stderr = %q, want to contain %q", stderr, "ERR: plain failure")
	}
}

func TestHandleDisplayError_NilError_ExitsCode1(t *testing.T) {
	assertExit(t, true, 1, func() {
		handleDisplayError(nil, false, false, true)
	})
}

func TestHandleDisplayError_SyscallErrno_UsesHRESULTExitCode(t *testing.T) {
	stderr := captureStderr(t, func() {
		assertExit(t, false, 0x1234, func() {
			handleDisplayError(syscall.Errno(0x1234), true, false, false)
		})
	})

	if !strings.Contains(stderr, "HRESULT: 0x1234") {
		t.Fatalf("stderr = %q, want to contain %q", stderr, "HRESULT: 0x1234")
	}
}

func TestMain_UsesExecutableAndArgs(t *testing.T) {
	const exePath = "/tmp/Arch.exe"

	origExit := exitFunc
	origExecutable := executableFunc
	origArgs := os.Args
	defer func() {
		exitFunc = origExit
		executableFunc = origExecutable
		os.Args = origArgs
	}()
	executableFunc = func() (string, error) { return exePath, nil }
	exitFunc = func(bool, int) {
		t.Fatal("unexpected exit")
	}
	os.Args = []string{"Arch.exe", "version"}

	main()
}

func assertExit(t *testing.T, wantPause bool, wantCode int, run func()) {
	t.Helper()

	exitFunc = func(pause bool, code int) {
		if pause != wantPause {
			t.Fatalf("pause = %v, want %v", pause, wantPause)
		}
		if code != wantCode {
			t.Fatalf("exit code = %d, want %d", code, wantCode)
		}
		panic(exitSignal{})
	}

	defer func() {
		r := recover()
		if _, ok := r.(exitSignal); ok {
			return
		}
		if r == nil {
			t.Fatalf("expected exit(pause=%v, code=%d) but it was not called", wantPause, wantCode)
		}
		panic(r)
	}()

	run()
}

func captureStderr(t *testing.T, run func()) string {
	t.Helper()

	origStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stderr = w

	run()

	if err := w.Close(); err != nil {
		t.Fatalf("stderr writer close failed: %v", err)
	}
	os.Stderr = origStderr

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("stderr read failed: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("stderr reader close failed: %v", err)
	}
	return string(data)
}
