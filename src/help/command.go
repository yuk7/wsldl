package help

import (
	"github.com/yuk7/wsldl/lib/cmdline"
)

// GetCommand returns the help command structure
func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"help", "--help", "-h", "/?"},
		Help: func(distroName string, isListQuery bool) string {
			return getHelpMessage()
		},
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
			if command.Help != nil {
				help := command.Help(distroName, false)
				if help != "" {
					helpStrs += "\n" + indentString(help) + "\n"
				}
			}
		}
	}

	if len(args) == 0 || helpStrs == "" {
		for _, c := range commands {
			if c.Help != nil {
				help := c.Help(distroName, true)
				if help != "" {
					helpStrs += "\n" + indentString(help) + "\n"
				}
			}
		}
	}
	helpStrs = "Usage:" + helpStrs
	print(helpStrs)
}
