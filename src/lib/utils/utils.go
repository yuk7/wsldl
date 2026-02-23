package utils

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
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
