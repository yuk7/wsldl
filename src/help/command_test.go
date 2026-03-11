package help

import (
	"testing"

	"github.com/yuk7/wsldl/lib/cmdline"
)

func TestGetCommand(t *testing.T) {
	t.Parallel()

	cmd := GetCommand()
	if len(cmd.Names) == 0 {
		t.Fatal("Names is empty")
	}
	if cmd.HelpText == nil || cmd.Run == nil {
		t.Fatal("HelpText or Run is nil")
	}
	if got := cmd.HelpText(); got == "" {
		t.Fatal("help message is empty")
	}
	if err := cmd.Run("Arch", nil); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
}

func TestIndentString(t *testing.T) {
	t.Parallel()

	got := indentString("a\nb")
	want := "    a\n    b"
	if got != want {
		t.Fatalf("indentString = %q, want %q", got, want)
	}
}

func TestShowHelpFromCommands_DoesNotPanic(t *testing.T) {
	t.Parallel()

	commands := []cmdline.Command{
		{
			Names: []string{"get"},
			HelpText: func() string {
				return "get"
			},
		},
	}
	ShowHelpFromCommands(commands, "Arch", []string{"get"})
	ShowHelpFromCommands(commands, "Arch", nil)
}

func TestCommandVisible_NilIsTrue(t *testing.T) {
	t.Parallel()

	got := commandVisible(cmdline.Command{}, "Arch")
	if !got {
		t.Fatal("commandVisible(nil) = false, want true")
	}
}

func TestCommandHelpText_NilIsEmpty(t *testing.T) {
	t.Parallel()

	got := commandHelpText(cmdline.Command{})
	if got != "" {
		t.Fatalf("commandHelpText(nil) = %q, want empty", got)
	}
}
