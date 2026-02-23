package console

import (
	"os"
	"strings"
	"syscall"
	"unsafe"

	ps "github.com/mitchellh/go-ps"
)

const (
	// ConsoleProcNames is console process list for detect parent console process
	ConsoleProcNames = "cmd.exe,powershell.exe,wsl.exe,WindowsTerminal.exe,flute.exe,FluentTerminal.SystemTray.exe,winpty-agent.exe"
)

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
