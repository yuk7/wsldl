package cmdline

import (
	"errors"
)

// RunSubCommand executes a subcommand
func RunSubCommand(commands []Command, mismatch func() error, distroName string, args []string) error {
	if len(args) > 0 {
		command, err := FindCommandFromName(commands, args[0])
		if err == nil {
			return command.Run(distroName, args[1:])
		}
	}
	return mismatch()
}

// FindCommandFromName finds a command by name
func FindCommandFromName(commands []Command, commandName string) (Command, error) {
	for _, c := range commands {
		for _, n := range c.Names {
			if n == commandName {
				return c, nil
			}
		}
	}
	return Command{}, errors.New("command not found")
}
