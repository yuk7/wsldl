package backup

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    backup [contents]")
	println("      - `--tar`: Output backup.tar to the current directory.")
	println("      - `--reg`: Output settings registry file to the current directory.")
}
