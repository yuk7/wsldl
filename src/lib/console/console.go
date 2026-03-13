//go:build windows

package console

import (
	"os"
	"strings"
	"syscall"
	"unsafe"

	ps "github.com/mitchellh/go-ps"
)

type kernel32Proc interface {
	Call(...uintptr) (uintptr, uintptr, error)
	Addr() uintptr
}

type kernel32DLL interface {
	FindProc(name string) (kernel32Proc, error)
}

type syscallKernel32DLL struct {
	dll *syscall.DLL
}

func (d syscallKernel32DLL) FindProc(name string) (kernel32Proc, error) {
	proc, err := d.dll.FindProc(name)
	if err != nil {
		return nil, err
	}
	return syscallKernel32Proc{proc: proc}, nil
}

type syscallKernel32Proc struct {
	proc *syscall.Proc
}

func (p syscallKernel32Proc) Call(args ...uintptr) (uintptr, uintptr, error) {
	return p.proc.Call(args...)
}

func (p syscallKernel32Proc) Addr() uintptr {
	return p.proc.Addr()
}

var (
	findProcessFunc  = ps.FindProcess
	getPPIDFunc      = syscall.Getppid
	loadKernel32Func = func() (kernel32DLL, error) {
		dll, err := syscall.LoadDLL("Kernel32.dll")
		if err != nil {
			return nil, err
		}
		return syscallKernel32DLL{dll: dll}, nil
	}
	getStdHandleFunc       = syscall.GetStdHandle
	utf16PtrFromStringFunc = syscall.UTF16PtrFromString
	syscallNFunc           = syscall.SyscallN
)

const (
	// ConsoleProcNames is console process list for detect parent console process
	ConsoleProcNames = "cmd.exe,powershell.exe,wsl.exe,WindowsTerminal.exe,flute.exe,FluentTerminal.SystemTray.exe,winpty-agent.exe"
)

// IsParentConsole gets is parent process is console or not
func IsParentConsole() (res bool, err error) {
	list := strings.Split(ConsoleProcNames, ",")
	info, err := findProcessFunc(getPPIDFunc())
	if err != nil {
		return
	}
	if info == nil {
		return false, nil
	}

	parentName := info.Executable()
	for _, item := range list {
		if strings.EqualFold(parentName, item) {
			res = true
			return
		}
	}

	grandParent, err := findProcessFunc(info.PPid())
	if err != nil {
		return
	}
	if grandParent == nil {
		return false, nil
	}
	grandParentName := grandParent.Executable()
	for _, item := range list {
		if strings.EqualFold(grandParentName, item) {
			res = true
			return
		}
	}

	res = false
	return
}

// FreeConsole calls FreeConsole API in Windows kernel32
func FreeConsole() error {
	kernel32, err := loadKernel32Func()
	if err != nil {
		return err
	}
	proc, err := kernel32.FindProc("FreeConsole")
	if err != nil {
		return err
	}
	proc.Call()
	return nil
}

// AllocConsole calls AllocConsole API in Windows kernel32
func AllocConsole() {
	kernel32, err := loadKernel32Func()
	if err != nil {
		return
	}
	alloc, err := kernel32.FindProc("AllocConsole")
	if err != nil {
		return
	}
	alloc.Call()

	hout, _ := getStdHandleFunc(syscall.STD_OUTPUT_HANDLE)
	herr, _ := getStdHandleFunc(syscall.STD_ERROR_HANDLE)
	hin, _ := getStdHandleFunc(syscall.STD_INPUT_HANDLE)
	os.Stdout = os.NewFile(uintptr(hout), "/dev/stdout")
	os.Stderr = os.NewFile(uintptr(herr), "/dev/stderr")
	os.Stdin = os.NewFile(uintptr(hin), "/dev/stdin")
}

// SetConsoleTitle calls SetConsoleTitleW API in Windows kernel32
func SetConsoleTitle(title string) {
	kernel32, err := loadKernel32Func()
	if err != nil {
		return
	}
	proc, err := kernel32.FindProc("SetConsoleTitleW")
	if err != nil {
		return
	}
	pTitle, err := utf16PtrFromStringFunc(title)
	if err != nil {
		return
	}
	syscallNFunc(proc.Addr(), uintptr(unsafe.Pointer(pTitle)))
}
