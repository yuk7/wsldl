package backup

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    backup [contents]")
	println("      - `--tar`: Output backup.tar to the current directory.")
	println("      - `--tgz`: Output backup.tar.gz to the current directory.")
	println("      - `--vhdx`: Output backup.ext4.vhdx to the current directory. (WSL2 only)")
	println("      - `--vhdxgz`: Output backup.ext4.vhdx.gz to the current directory. (WSL2 only)")
	println("      - `--reg`: Output settings registry file to the current directory.")
}
