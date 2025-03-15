package run

func getHelpMessageNoArgs() string {
	return "" +
		"<no args>\n" +
		"  - Open a new shell with your default settings. \n" +
		"    Inherit current directory (with exception that %USERPROFILE% is changed to $HOME)."
}

// getHelpMessage returns the help message for the run command
func getHelpMessage() string {
	return "" +
		"run <command line>\n" +
		"  - Run the given command line in that instance. Inherit current directory."
}

// getHelpMessageP returns the help message for the runp command
func getHelpMessageP() string {
	return "" +
		"runp <command line (includes windows path)>\n" +
		"  - Run the given command line in that instance after converting its path."
}
