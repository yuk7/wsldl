package clean

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    clean")
	println("      - Uninstall that instance.")
}
