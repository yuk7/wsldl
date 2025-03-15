package backup

// getHelpMessage returns the help message
func getHelpMessage() string {
	return "" +
		"backup [file name]\n" +
		"  - `*.tar`: Output backup tar file.\n" +
		"  - `*.tar.gz`: Output backup tar.gz file.\n" +
		"  - `*.ext4.vhdx`: Output backup ext4.vhdx file. (WSL2 only)\n" +
		"  - `*.ext4.vhdx.gz`: Output backup ext4.vhdx.gz file. (WSL2 only)\n" +
		"  - `*.reg`: Output settings registry file."
}
