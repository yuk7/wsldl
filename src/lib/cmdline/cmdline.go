package cmdline

import (
	"errors"
)

var (
	ErrCommandNotFound = errors.New("command not found")
	ErrNoDefault       = errors.New("default command not found")
)

// RunSubCommand executes a subcommand
func RunSubCommand(commands []Command, distroName string, args []string) error {
	if len(args) == 0 {
		for _, c := range commands {
			if c.IsDefault {
				return c.Run(distroName, args)
			}
		}
		return ErrNoDefault
	}

	command, err := FindCommandFromName(commands, args[0])
	if err == nil {
		return command.Run(distroName, args[1:])
	}
	return err
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
	return Command{}, ErrCommandNotFound
}
