package get

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    get [setting [value]]")
	println("      - `--default-uid`: Get the default user uid in this instance.")
	println("      - `--append-path`: Get true/false status of Append Windows PATH to $PATH.")
	println("      - `--mount-drive`: Get true/false status of Mount drives.")
	println("      - `--wsl-version`: Get the version os the WSL (1/2) of this instance.")
	println("      - `--default-term`: Get Default Terminal type of this instance launcher.")
	println("      - `--wt-profile-name`: Get Profile Name from Windows Terminal")
	println("      - `--lxguid`: Get WSL GUID key for this instance.")
}
