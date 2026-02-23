package utils

import (
	"syscall"
)

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
