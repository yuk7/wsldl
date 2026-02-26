package main

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

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
