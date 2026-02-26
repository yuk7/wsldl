package cmdline

import (
	"errors"
	"reflect"
	"testing"
)

func TestFindCommandFromName(t *testing.T) {
	t.Parallel()

	commands := []Command{
		{Names: []string{"install", "i"}},
		{Names: []string{"run", "r"}},
	}

	got, err := FindCommandFromName(commands, "r")
	if err != nil {
		t.Fatalf("FindCommandFromName returned error: %v", err)
	}
	if !reflect.DeepEqual(got.Names, []string{"run", "r"}) {
		t.Fatalf("FindCommandFromName names = %v, want %v", got.Names, []string{"run", "r"})
	}
}

func TestFindCommandFromName_NotFound(t *testing.T) {
	t.Parallel()

	_, err := FindCommandFromName([]Command{{Names: []string{"install"}}}, "run")
	if err == nil {
		t.Fatal("FindCommandFromName succeeded unexpectedly")
	}
}

func TestRunSubCommand_RunsMatchedCommand(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("run error")
	var gotDistro string
	var gotArgs []string

	commands := []Command{
		{
			Names: []string{"run"},
			Run: func(distroName string, args []string) error {
				gotDistro = distroName
				gotArgs = append([]string(nil), args...)
				return wantErr
			},
		},
	}

	mismatchCalled := false
	err := RunSubCommand(commands, func() error {
		mismatchCalled = true
		return nil
	}, "Arch", []string{"run", "-u", "root"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("RunSubCommand error = %v, want %v", err, wantErr)
	}
	if mismatchCalled {
		t.Fatal("mismatch callback called unexpectedly")
	}
	if gotDistro != "Arch" {
		t.Fatalf("run distroName = %q, want %q", gotDistro, "Arch")
	}
	if !reflect.DeepEqual(gotArgs, []string{"-u", "root"}) {
		t.Fatalf("run args = %v, want %v", gotArgs, []string{"-u", "root"})
	}
}

func TestRunSubCommand_FallsBackToMismatch(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("mismatch")
	called := 0

	err := RunSubCommand(nil, func() error {
		called++
		return wantErr
	}, "Arch", nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("RunSubCommand error = %v, want %v", err, wantErr)
	}
	if called != 1 {
		t.Fatalf("mismatch callback count = %d, want 1", called)
	}
}
