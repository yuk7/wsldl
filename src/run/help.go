package run

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    <no args>")
	println("      - Open a new shell with your default settings. ")
	println("        Inherit current directory (with exception that %%USERPROFILE%% is changed to $HOME).")
	println()
	println("    run <command line>")
	println("      - Run the given command line in that instance. Inherit current directory.")
	println()
	println("    runp <command line (includes windows path)>")
	println("      - Run the given command line in that instance after converting its path.")
}
