package help

import (
	"github.com/yuk7/wsldl/lib/cmdline"
)

// GetCommand returns the help command structure
func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names:    []string{"help", "--help", "-h", "/?"},
		HelpText: getHelpMessage,
		Run: func(distroName string, args []string) error {
			println("Usage:")
			println(indentString(getHelpMessage()))
			return nil
		},
	}
}

// ShowHelp prints help message
func ShowHelpFromCommands(commands []cmdline.Command, distroName string, args []string) {
	helpStrs := ""
	if len(args) > 0 {
		command, err := cmdline.FindCommandFromName(commands, args[0])
		if err == nil {
			help := commandHelpText(command)
			if help != "" {
				helpStrs += "\n" + indentString(help) + "\n"
			}
		}
	}

	if len(args) == 0 || helpStrs == "" {
		for _, c := range commands {
			if !commandVisible(c, distroName) {
				continue
			}
			help := commandHelpText(c)
			if help != "" {
				helpStrs += "\n" + indentString(help) + "\n"
			}
		}
	}
	helpStrs = "Usage:" + helpStrs
	print(helpStrs)
}

func commandVisible(command cmdline.Command, distroName string) bool {
	if command.Visible == nil {
		return true
	}
	return command.Visible(distroName)
}

func commandHelpText(command cmdline.Command) string {
	if command.HelpText == nil {
		return ""
	}
	return command.HelpText()
}
