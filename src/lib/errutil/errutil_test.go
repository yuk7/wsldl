package errutil

import (
	"errors"
	"testing"
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
