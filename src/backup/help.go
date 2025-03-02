package backup

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    backup [file name]")
	println("      - `*.tar`: Output backup tar file.")
	println("      - `*.tar.gz`: Output backup tar.gz file.")
	println("      - `*.ext4.vhdx`: Output backup ext4.vhdx file. (WSL2 only)")
	println("      - `*.ext4.vhdx.gz`: Output backup ext4.vhdx.gz file. (WSL2 only)")
	println("      - `*.reg`: Output settings registry file.")
}
