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

func TestShowHelpFromCommands_SkipsInvisibleCommands(t *testing.T) {
	t.Parallel()

	invisibleCalled := 0
	commands := []cmdline.Command{
		{
			Names: []string{"hidden"},
			Visible: func(string) bool {
				invisibleCalled++
				return false
			},
			HelpText: func() string {
				t.Fatal("hidden HelpText should not be called")
				return ""
			},
		},
		{
			Names: []string{"shown"},
			HelpText: func() string {
				return "shown"
			},
		},
	}

	ShowHelpFromCommands(commands, "Arch", nil)

	if invisibleCalled != 1 {
		t.Fatalf("invisible Visible call count = %d, want 1", invisibleCalled)
	}
}

func TestCommandVisible_NilIsTrue(t *testing.T) {
	t.Parallel()

	got := commandVisible(cmdline.Command{}, "Arch")
	if !got {
		t.Fatal("commandVisible(nil) = false, want true")
	}
}

func TestCommandVisible_UsesVisibleFunc(t *testing.T) {
	t.Parallel()

	cmd := cmdline.Command{
		Visible: func(distroName string) bool {
			return distroName == "Arch"
		},
	}
	if !commandVisible(cmd, "Arch") {
		t.Fatal("commandVisible(Arch) = false, want true")
	}
	if commandVisible(cmd, "Ubuntu") {
		t.Fatal("commandVisible(Ubuntu) = true, want false")
	}
}

func TestCommandHelpText_NilIsEmpty(t *testing.T) {
	t.Parallel()

	got := commandHelpText(cmdline.Command{})
	if got != "" {
		t.Fatalf("commandHelpText(nil) = %q, want empty", got)
	}
}
