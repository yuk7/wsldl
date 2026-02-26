//go:build !windows

package console

const (
	// ConsoleProcNames is console process list for detect parent console process
	ConsoleProcNames = "cmd.exe,powershell.exe,wsl.exe,WindowsTerminal.exe,flute.exe,FluentTerminal.SystemTray.exe,winpty-agent.exe"
)

// IsParentConsole gets is parent process is console or not.
func IsParentConsole() (bool, error) {
	return true, nil
}

// FreeConsole is a no-op outside Windows.
func FreeConsole() error {
	return nil
}

// AllocConsole is a no-op outside Windows.
func AllocConsole() {}

// SetConsoleTitle is a no-op outside Windows.
func SetConsoleTitle(title string) {}
