package utils

import "strings"

// DQEscapeString is escape string with double quote
func DQEscapeString(str string) string {
	if strings.Contains(str, " ") {
		str = strings.Replace(str, "\"", "\\\"", -1)
		str = "\"" + str + "\""
	}
	return str
}
