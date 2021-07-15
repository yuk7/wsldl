package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/fatih/color"
	ps "github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows/registry"
)

const (
	// LxssBaseRoot is LOCAL_MACHINE
	LxssBaseRoot = registry.CURRENT_USER
	// LxssBaseKey is path of lxss registry
	LxssBaseKey = "Software\\Microsoft\\Windows\\CurrentVersion\\Lxss"
	// ConsoleProcNames is console process list for detect parent console process
	ConsoleProcNames = "cmd.exe,powershell.exe,wsl.exe,WindowsTerminal.exe,winpty-agent.exe"
	// WsldlTermKey is registry key name used for wsldl terminal infomation
	WsldlTermKey = "wsldl-term"
	// FlagWsldlTermDefault is default terminal (conhost)
	FlagWsldlTermDefault = 0
	// FlagWsldlTermWT is Windows Terminal
	FlagWsldlTermWT = 1
	// FlagWsldlTermFlute is Fluent Terminal
	FlagWsldlTermFlute = 2
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

//WslGetUUID gets distro guid key
func WslGetUUID(distributionName string) (uuid string, err error) {
	uuidList, tmpErr := WslGetUUIDList()
	if tmpErr != nil {
		err = tmpErr
		return
	}

	errStr := ""
	for _, loopUUID := range uuidList {
		key, loopErr := registry.OpenKey(LxssBaseRoot, LxssBaseKey+"\\"+loopUUID, registry.READ)
		if loopErr == nil || loopErr == io.EOF {
			str, _, itemErr := key.GetStringValue("DistributionName")
			if itemErr == nil || itemErr == io.EOF {
				if strings.EqualFold(str, distributionName) {
					uuid = loopUUID
					return
				}
			} else {
				errStr += "\n" + "    " + loopUUID + ":" + itemErr.Error()
			}
		} else {
			errStr += "\n" + "    " + loopUUID + ":" + loopErr.Error()
		}
	}
	err = errors.New("Registry Key Not found\n" + errStr)

	return
}

//WslGetUUIDList gets guid key lists
func WslGetUUIDList() (uuidList []string, err error) {
	baseKey, tmpErr := registry.OpenKey(LxssBaseRoot, LxssBaseKey, registry.READ)
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	uuidList, tmpErr = baseKey.ReadSubKeyNames(1024)
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	return
}

// WsldlSetTerminalInfo sets terminal number from registry
func WsldlSetTerminalInfo(uuid string, value int) error {
	key, tmpErr := registry.OpenKey(LxssBaseRoot, LxssBaseKey+"\\"+uuid, registry.SET_VALUE)
	if tmpErr != nil && tmpErr != io.EOF {
		return tmpErr
	}
	tmpErr = key.SetDWordValue(WsldlTermKey, uint32(value))
	if tmpErr != nil && tmpErr != io.EOF {
		return tmpErr
	}
	return nil
}

// WsldlGetTerminalInfo gets terminal number from registry
func WsldlGetTerminalInfo(uuid string) (res int, err error) {
	key, tmpErr := registry.OpenKey(LxssBaseRoot, LxssBaseKey+"\\"+uuid, registry.READ)
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	num, _, tmpErr := key.GetIntegerValue(WsldlTermKey)
	if tmpErr == syscall.ERROR_FILE_NOT_FOUND {
		// not found is okay, it's just unconfigured.
		res = FlagWsldlTermDefault
		return
	}
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	res = int(num)
	return
}

// WslGetDistroName gets distro name from registry
func WslGetDistroName(uuid string) (res string, err error) {
	key, tmpErr := registry.OpenKey(LxssBaseRoot, LxssBaseKey+"\\"+uuid, registry.READ)
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	res, _, tmpErr = key.GetStringValue("DistributionName")
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	return
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

	info, err = ps.FindProcess(info.PPid())
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

	return
}

// SetConsoleTitle calls SetConsoleTitleW API in Windows kernel32
func SetConsoleTitle(title string) {
	kernel32, _ := syscall.LoadDLL("Kernel32.dll")
	proc, _ := kernel32.FindProc("SetConsoleTitleW")
	pTitle, _ := syscall.UTF16PtrFromString(title)
	syscall.Syscall(proc.Addr(), 1, uintptr(unsafe.Pointer(pTitle)), 0, 0)
	return
}

// ErrorExit shows error message and exit
func ErrorExit(err error, showmsg bool, showcolor bool, pause bool) {
	var errno syscall.Errno
	if err == nil {
		if showmsg {
			ErrorRedPrintln("ERR: unknown error")
			Exit(pause, 1)
		}
	}
	if showmsg {
		if showcolor {
			ErrorRedPrintln("ERR: " + err.Error())
		} else {
			fmt.Fprintln(os.Stderr, "ERR: "+err.Error())
		}

	}
	if errors.As(err, &errno) {
		if showmsg {
			fmt.Fprintf(os.Stderr, "HRESULT: 0x%x\n", int(errno))
		}
		Exit(pause, int(errno))
	} else if err == os.ErrInvalid {
		if showmsg {
			efPath, _ := os.Executable()
			exeName := filepath.Base(efPath)
			fmt.Fprintln(os.Stderr, "Your command may be incorrect.")
			fmt.Fprintf(os.Stderr, "You can get help with `%s help`.\n", exeName)
		}
	} else {
		if showmsg {
			fmt.Fprintf(os.Stderr, "%#v\n", err)
		}
		Exit(pause, 1)
	}
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
