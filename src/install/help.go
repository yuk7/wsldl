package install

// getHelpMessageNoArgs returns the help message
func getHelpMessageNoArgs() string {
	return "" +
		"<no args>\n" +
		"  - Install a new instance with default settings."
}

// getHelpMessage returns the help message
func getHelpMessage() string {
	return "" +
		"install [rootfs file]\n" +
		"  - Install a new instance with your given rootfs file\n" +
		"    You can use .tar(.gz) or .ext4.vhdx(.gz)"
}
