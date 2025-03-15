package help

import (
	"strings"
)

// getHelpMessage returns the help message
func getHelpMessage() string {
	return "" +
		"help\n" +
		"  - Print this usage message."
}

// indentString indents the message
func indentString(message string) string {
	return "    " + strings.ReplaceAll(message, "\n", "\n    ")
}
