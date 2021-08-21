package install

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    <no args>")
	println("      - Install a new instance with default settings.")
	println()
	println("    install <rootfs file>")
	println("      - Install a new instance with your given rootfs file")
	println("        You can use .tar(.gz) or .ext4.vhdx(.gz)")
}
