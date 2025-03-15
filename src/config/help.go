package config

// getHelpMessage returns the help message
func getHelpMessage() string {
	return "" +
		"config [setting [value]]\n" +
		"  - `--default-user <user>`: Set the default user of this instance to <user>.\n" +
		"  - `--default-uid <uid>`: Set the default user uid of this instance to <uid>.\n" +
		"  - `--append-path <true|false>`: Switch of Append Windows PATH to $PATH\n" +
		"  - `--mount-drive <true|false>`: Switch of Mount drives\n" +
		"  - `--wsl-version <1|2>`: Set the WSL version of this instance to <1 or 2>\n" +
		"  - `--default-term <default|wt|flute>`: Set default type of terminal window."
}
