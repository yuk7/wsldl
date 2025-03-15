package get

// getHelpMessage returns the help message
func getHelpMessage() string {
	return "" +
		"get [setting [value]]\n" +
		"  - `--default-uid`: Get the default user uid in this instance.\n" +
		"  - `--append-path`: Get true/false status of Append Windows PATH to $PATH.\n" +
		"  - `--mount-drive`: Get true/false status of Mount drives.\n" +
		"  - `--wsl-version`: Get the version os the WSL (1/2) of this instance.\n" +
		"  - `--default-term`: Get Default Terminal type of this instance launcher.\n" +
		"  - `--wt-profile-name`: Get Profile Name from Windows Terminal\n" +
		"  - `--lxguid`: Get WSL GUID key for this instance."
}
