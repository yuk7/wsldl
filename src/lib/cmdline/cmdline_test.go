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
	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("FindCommandFromName error = %v, want %v", err, ErrCommandNotFound)
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

	err := RunSubCommand(commands, "Arch", []string{"run", "-u", "root"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("RunSubCommand error = %v, want %v", err, wantErr)
	}
	if gotDistro != "Arch" {
		t.Fatalf("run distroName = %q, want %q", gotDistro, "Arch")
	}
	if !reflect.DeepEqual(gotArgs, []string{"-u", "root"}) {
		t.Fatalf("run args = %v, want %v", gotArgs, []string{"-u", "root"})
	}
}

func TestRunSubCommand_DefaultCommand(t *testing.T) {
	t.Parallel()

	called := false
	err := RunSubCommand([]Command{
		{
			IsDefault: true,
			Run: func(distroName string, args []string) error {
				called = true
				if distroName != "Arch" {
					t.Fatalf("distroName = %q, want %q", distroName, "Arch")
				}
				if args != nil {
					t.Fatalf("args = %v, want nil", args)
				}
				return nil
			},
		},
	}, "Arch", nil)
	if err != nil {
		t.Fatalf("RunSubCommand returned error: %v", err)
	}
	if !called {
		t.Fatal("default command was not called")
	}
}

func TestRunSubCommand_NoDefault(t *testing.T) {
	t.Parallel()

	err := RunSubCommand([]Command{{Names: []string{"run"}}}, "Arch", nil)
	if !errors.Is(err, ErrNoDefault) {
		t.Fatalf("RunSubCommand error = %v, want %v", err, ErrNoDefault)
	}
}

func TestRunSubCommand_EmptyNamesCommandIsNotDefault(t *testing.T) {
	t.Parallel()

	called := false
	err := RunSubCommand([]Command{
		{
			Run: func(distroName string, args []string) error {
				called = true
				return nil
			},
		},
	}, "Arch", nil)
	if !errors.Is(err, ErrNoDefault) {
		t.Fatalf("RunSubCommand error = %v, want %v", err, ErrNoDefault)
	}
	if called {
		t.Fatal("unnamed command should not be treated as default")
	}
}

func TestRunSubCommand_CommandNotFound(t *testing.T) {
	t.Parallel()

	err := RunSubCommand([]Command{{Names: []string{"run"}}}, "Arch", []string{"unknown"})
	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("RunSubCommand error = %v, want %v", err, ErrCommandNotFound)
	}
}
