package run

import "strings"

func getHelpMessageNoArgs() string {
	return "" +
		"<no args>\n" +
		"  - Open a new shell with your default settings. \n" +
		"    Inherit current directory (with exception that %%USERPROFILE%% is changed to $HOME)."
}

func getHelpMessage() string {
	return "" +
		"<no args>\n" +
		"  - Open a new shell with your default settings. \n" +
		"    Inherit current directory (with exception that %%USERPROFILE%% is changed to $HOME).\n" +
		"\n" +
		"run <command line>\n" +
		"  - Run the given command line in that instance. Inherit current directory.\n" +
		"\n" +
		"runp <command line (includes windows path)>\n" +
		"  - Run the given command line in that instance after converting its path."
}

func getHelpMessageP() string {
	return "" +
		"runp <command line (includes windows path)>\n" +
		"  - Run the given command line in that instance after converting its path."
}

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println(indentString(getHelpMessageNoArgs()))
	println()
	println(indentString(getHelpMessage()))
	println()
	println(indentString(getHelpMessageP()))
}

func indentString(message string) string {
	return strings.ReplaceAll(message, "\n", "\n    ")
}
