package errutil

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestDisplayError(t *testing.T) {
	t.Parallel()

	base := errors.New("base")
	err := NewDisplayError(base, true, false, true)
	if err == nil {
		t.Fatal("NewDisplayError returned nil")
	}

	de, ok := err.(*DisplayError)
	if !ok {
		t.Fatalf("error type = %T, want *DisplayError", err)
	}
	if !errors.Is(de, base) {
		t.Fatalf("DisplayError unwrap mismatch: got %v, want %v", de.Unwrap(), base)
	}
	if got := de.Error(); got != "base" {
		t.Fatalf("DisplayError.Error() = %q, want %q", got, "base")
	}
	if de.ShowMsg != true || de.ShowColor != false || de.Pause != true {
		t.Fatalf("DisplayError fields mismatch: %+v", de)
	}
}

func TestDisplayErrorNilSafety(t *testing.T) {
	t.Parallel()

	var de *DisplayError
	if got := de.Error(); got != "" {
		t.Fatalf("nil DisplayError.Error() = %q, want empty string", got)
	}
	if de.Unwrap() != nil {
		t.Fatal("nil DisplayError.Unwrap() should be nil")
	}
}

func TestNewDisplayErrorNilInput(t *testing.T) {
	t.Parallel()

	if err := NewDisplayError(nil, true, true, true); err != nil {
		t.Fatalf("NewDisplayError(nil, ...) = %v, want nil", err)
	}
}

func TestExitCodeError(t *testing.T) {
	t.Parallel()

	err := NewExitCodeError(42, true)
	ec, ok := err.(*ExitCodeError)
	if !ok {
		t.Fatalf("error type = %T, want *ExitCodeError", err)
	}
	if ec.Code != 42 || ec.Pause != true {
		t.Fatalf("ExitCodeError fields mismatch: %+v", ec)
	}
	if got := ec.Error(); got != "exit requested" {
		t.Fatalf("ExitCodeError.Error() = %q, want %q", got, "exit requested")
	}
}

func TestFormatError(t *testing.T) {
	t.Parallel()

	if got := FormatError(nil); got != "ERR: unknown error" {
		t.Fatalf("FormatError(nil) = %q, want %q", got, "ERR: unknown error")
	}
	if got := FormatError(errors.New("boom")); got != "ERR: boom" {
		t.Fatalf("FormatError(non-nil) = %q, want %q", got, "ERR: boom")
	}
}

func TestMustExecutable(t *testing.T) {
	if p := MustExecutable(); p == "" {
		t.Fatal("MustExecutable() returned empty path")
	}
}

func TestMustExecutable_PanicsOnExecutableError(t *testing.T) {
	orig := executablePathFunc
	executablePathFunc = func() (string, error) {
		return "", errors.New("boom")
	}
	t.Cleanup(func() {
		executablePathFunc = orig
	})

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("MustExecutable did not panic")
		}
		msg, ok := recovered.(string)
		if !ok {
			t.Fatalf("panic type = %T, want string", recovered)
		}
		if !strings.Contains(msg, "failed to get executable path") {
			t.Fatalf("panic message = %q, want to contain %q", msg, "failed to get executable path")
		}
	}()

	_ = MustExecutable()
}

func TestExit_ExitsWithProvidedCode(t *testing.T) {
	if os.Getenv("WSDLD_TEST_EXIT_HELPER") == "1" {
		Exit(false, 7)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExit_ExitsWithProvidedCode$")
	cmd.Env = append(os.Environ(), "WSDLD_TEST_EXIT_HELPER=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("helper process succeeded unexpectedly")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("error type = %T, want *exec.ExitError", err)
	}
	if exitErr.ExitCode() != 7 {
		t.Fatalf("exit code = %d, want %d", exitErr.ExitCode(), 7)
	}
}

func TestExit_WithPause_PrintsPromptAndExits(t *testing.T) {
	if os.Getenv("WSDLD_TEST_EXIT_PAUSE_HELPER") == "1" {
		Exit(true, 5)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExit_WithPause_PrintsPromptAndExits$")
	cmd.Env = append(os.Environ(), "WSDLD_TEST_EXIT_PAUSE_HELPER=1")
	cmd.Stdin = strings.NewReader("\n")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err == nil {
		t.Fatal("helper process succeeded unexpectedly")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("error type = %T, want *exec.ExitError", err)
	}
	if exitErr.ExitCode() != 5 {
		t.Fatalf("exit code = %d, want %d", exitErr.ExitCode(), 5)
	}
	if !strings.Contains(stdout.String(), "Press enter to exit...") {
		t.Fatalf("stdout = %q, want to contain pause prompt", stdout.String())
	}
}

func TestColorPrintFunctions_WriteToConfiguredWriters(t *testing.T) {
	t.Parallel()

	oldErr := color.Error
	oldOut := color.Output
	t.Cleanup(func() {
		color.Error = oldErr
		color.Output = oldOut
	})

	var errBuf bytes.Buffer
	var outBuf bytes.Buffer
	color.Error = &errBuf
	color.Output = &outBuf

	ErrorRedPrintln("err message")
	StdoutGreenPrintln("ok message")

	if !strings.Contains(errBuf.String(), "err message") {
		t.Fatalf("stderr output = %q, want to contain %q", errBuf.String(), "err message")
	}
	if !strings.Contains(outBuf.String(), "ok message") {
		t.Fatalf("stdout output = %q, want to contain %q", outBuf.String(), "ok message")
	}
}
