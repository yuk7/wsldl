package install

import (
	"errors"
	"os"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
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

func assertDisplayError(t *testing.T, err error) *errutil.DisplayError {
	t.Helper()
	var de *errutil.DisplayError
	if !errors.As(err, &de) {
		t.Fatalf("error type = %T, want *errutil.DisplayError", err)
	}
	return de
}
