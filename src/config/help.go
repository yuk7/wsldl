package config

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    config [setting [value]]")
	println("      - `--default-user <user>`: Set the default user of this instance to <user>.")
	println("      - `--default-uid <uid>`: Set the default user uid of this instance to <uid>.")
	println("      - `--append-path <true|false>`: Switch of Append Windows PATH to $PATH")
	println("      - `--mount-drive <true|false>`: Switch of Mount drives")
	println("      - `--wsl-version <1|2>`: Set the WSL version of this instance to <1 or 2>")
	println("      - `--default-term <default|wt|flute>`: Set default type of terminal window.")
}
