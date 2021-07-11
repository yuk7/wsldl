package wslapi

import (
	"os"
	"syscall"
)

//sys	_WslIsDistributionRegistered(distributionName *uint16) (res bool) = wslapi.WslIsDistributionRegistered

// WslIsDistributionRegistered determines if a distribution is already registered
func WslIsDistributionRegistered(distributionName string) bool {
	pDistributionName, _ := syscall.UTF16PtrFromString(distributionName)
	return _WslIsDistributionRegistered(pDistributionName)
}

//sys	_WslRegisterDistribution(distributionName *uint16, tarGzFilename *uint16) (res error) = wslapi.WslRegisterDistribution

// WslRegisterDistribution registers a new distribution
func WslRegisterDistribution(distributionName, tarGzFilename string) error {
	pDistributionName, _ := syscall.UTF16PtrFromString(distributionName)
	pTarGzFilename, _ := syscall.UTF16PtrFromString(tarGzFilename)
	return _WslRegisterDistribution(pDistributionName, pTarGzFilename)
}

//sys	_WslLaunch(distributionName *uint16, command *uint16, useCurrentWorkingDirectory bool, stdIn syscall.Handle, stdOut syscall.Handle, stdErr syscall.Handle, process *syscall.Handle, exitCode *uintptr) (err error) = wslapi.WslLaunch

// WslLaunchInteractive launches the distribution with interactive shell
func WslLaunchInteractive(distributionName string, command string, useCurrentWorkingDirectory bool) (exitCode uintptr, err error) {
	pDistributionName, _ := syscall.UTF16PtrFromString(distributionName)
	pCommand, _ := syscall.UTF16PtrFromString(command)

	p, _ := syscall.GetCurrentProcess()
	stdin := syscall.Handle(0)
	stdout := syscall.Handle(0)
	stderr := syscall.Handle(0)

	syscall.DuplicateHandle(p, syscall.Handle(os.Stdin.Fd()), p, &stdin, 0, true, syscall.DUPLICATE_SAME_ACCESS)
	syscall.DuplicateHandle(p, syscall.Handle(os.Stdout.Fd()), p, &stdout, 0, true, syscall.DUPLICATE_SAME_ACCESS)
	syscall.DuplicateHandle(p, syscall.Handle(os.Stderr.Fd()), p, &stderr, 0, true, syscall.DUPLICATE_SAME_ACCESS)

	handle := syscall.Handle(0)
	_WslLaunch(pDistributionName, pCommand, useCurrentWorkingDirectory, stdin, stdout, stderr, &handle, &exitCode)
	syscall.WaitForSingleObject(handle, syscall.INFINITE)
	return
}
