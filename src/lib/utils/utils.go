package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/fatih/color"
	ps "github.com/mitchellh/go-ps"
)

const (
	// ConsoleProcNames is console process list for detect parent console process
	ConsoleProcNames = "cmd.exe,powershell.exe,wsl.exe,WindowsTerminal.exe,flute.exe,FluentTerminal.SystemTray.exe,winpty-agent.exe"
	// SpecialDirs is define path of special dirs
	SpecialDirs = "SystemDrive:,SystemRoot:,SystemRoot:System32,USERPROFILE:"
)

// DQEscapeString is escape string with double quote
func DQEscapeString(str string) string {
	if strings.Contains(str, " ") {
		str = strings.Replace(str, "\"", "\\\"", -1)
		str = "\"" + str + "\""
	}
	return str
}

// IsParentConsole gets is parent process is console or not
func IsParentConsole() (res bool, err error) {
	list := strings.Split(ConsoleProcNames, ",")
	info, err := ps.FindProcess(syscall.Getppid())
	if err != nil {
		return
	}

	parentName := info.Executable()
	for _, item := range list {
		if strings.EqualFold(parentName, item) {
			res = true
			return
		}
	}

	_, err = ps.FindProcess(info.PPid())
	if err != nil {
		return
	}
	for _, item := range list {
		if strings.EqualFold(parentName, item) {
			res = true
			return
		}
	}

	res = false
	return
}

// GetWindowsDirectory gets windows direcotry path
func GetWindowsDirectory() string {
	dir := os.Getenv("SYSTEMROOT")
	if dir != "" {
		return dir
	}
	dir = os.Getenv("WINDIR")
	if dir != "" {
		return dir
	}
	return "C:\\WINDOWS"

}

// IsCurrentDirSpecial gets whether the current directory is special (Windows, USEPROFILE)
func IsCurrentDirSpecial() bool {
	cdir, err := filepath.Abs(".")
	if err != nil {
		return true
	}
	sdarr := strings.Split(SpecialDirs, ",")
	for _, item := range sdarr {
		splititem := strings.Split(item, ":")
		itemdir := ""
		if splititem[0] != "" {
			itemdir = os.Getenv(splititem[0])
		}
		itemdir, err = filepath.Abs(itemdir + "\\" + splititem[1])
		if err != nil {
			return true
		}
		if strings.EqualFold(cdir, itemdir) {
			return true
		}
	}
	return false
}

// CreateProcessAndWait creating process and wait it
func CreateProcessAndWait(commandLine string) (res int, err error) {
	pCommandLine, _ := syscall.UTF16PtrFromString(commandLine)
	si := syscall.StartupInfo{}
	pi := syscall.ProcessInformation{}

	err = syscall.CreateProcess(nil, pCommandLine, nil, nil, false, 0, nil, nil, &si, &pi)
	if err != nil {
		return
	}
	_, err = syscall.WaitForSingleObject(pi.Process, syscall.INFINITE)
	var exitCode = uint32(0)
	syscall.GetExitCodeProcess(pi.Process, &exitCode)
	res = int(exitCode)
	return
}

// FreeConsole calls FreeConsole API in Windows kernel32
func FreeConsole() error {
	kernel32, _ := syscall.LoadDLL("Kernel32.dll")
	proc, err := kernel32.FindProc("FreeConsole")
	if err != nil {
		return err
	}
	proc.Call()
	return nil
}

// AllocConsole calls AllocConsole API in Windows kernel32
func AllocConsole() {
	kernel32, _ := syscall.LoadDLL("Kernel32.dll")
	alloc, _ := kernel32.FindProc("AllocConsole")
	alloc.Call()

	hout, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	herr, _ := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)
	hin, _ := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	os.Stdout = os.NewFile(uintptr(hout), "/dev/stdout")
	os.Stderr = os.NewFile(uintptr(herr), "/dev/stderr")
	os.Stdin = os.NewFile(uintptr(hin), "/dev/stdin")
}

// SetConsoleTitle calls SetConsoleTitleW API in Windows kernel32
func SetConsoleTitle(title string) {
	kernel32, _ := syscall.LoadDLL("Kernel32.dll")
	proc, _ := kernel32.FindProc("SetConsoleTitleW")
	pTitle, _ := syscall.UTF16PtrFromString(title)
	syscall.SyscallN(proc.Addr(), 1, uintptr(unsafe.Pointer(pTitle)))
}

// FormatError formats an error for CLI output.
func FormatError(err error) string {
	if err == nil {
		return "ERR: unknown error"
	}
	return "ERR: " + err.Error()
}

// Exit exits program
func Exit(pause bool, exitCode int) {
	if pause {
		fmt.Fprintf(os.Stdout, "Press enter to exit...")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
	os.Exit(exitCode)
}

// ErrorRedPrintln shows red string to stderr
func ErrorRedPrintln(str string) {
	color.New(color.FgRed).Fprintln(color.Error, str)
}

// StdoutGreenPrintln shows green string to stdout
func StdoutGreenPrintln(str string) {
	color.New(color.FgGreen).Fprintln(color.Output, str)
}
