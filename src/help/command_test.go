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
	if cmd.Help == nil || cmd.Run == nil {
		t.Fatal("Help or Run is nil")
	}
	if got := cmd.Help("Arch", true); got == "" {
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
			Help: func(distroName string, isListQuery bool) string {
				return "get"
			},
		},
	}
	ShowHelpFromCommands(commands, "Arch", []string{"get"})
	ShowHelpFromCommands(commands, "Arch", nil)
}
