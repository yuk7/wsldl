package wslapi

import (
	"os"
	"syscall"
)

const (
	//FlagEnableInterop is flag of interop feature
	FlagEnableInterop = 1
	//FlagAppendNTPath is flag of appending windows path
	FlagAppendNTPath = 2
	//FlagEnableDriveMounting is flag of mounting windows drive
	FlagEnableDriveMounting = 4
	//FlagEnableWsl2 is flag of enabled wsl2, read only
	FlagEnableWsl2 = 8
)

//sys	_WslGetDistributionConfiguration(distributionName *uint16, distributionVersion *uint32, defaultUID *uint64,  wslDistributionFlags *uint32, defaultEnv ***uint16, defaultEnvCnt *uint64) (res error) = wslapi.WslGetDistributionConfiguration

// WslGetDistributionConfiguration gets distribution configuration
func WslGetDistributionConfiguration(distributionName string) (distributionVersion uint32, defaultUID uint64, flags uint32, err error) {
	pDistributionName, _ := syscall.UTF16PtrFromString(distributionName)

	var pEnv **uint16
	var envCnt uint64

	err = _WslGetDistributionConfiguration(pDistributionName, &distributionVersion, &defaultUID, &flags, &pEnv, &envCnt)

	return
}

//sys	_WslConfigureDistribution(distributionName *uint16, defaultUID uint64, wslDistributionFlags uint32) (res error) = wslapi.WslConfigureDistribution

// WslConfigureDistribution configures distribution configuration
func WslConfigureDistribution(distributionName string, defaultUID uint64, wslDistributionFlags uint32) (err error) {
	pDistributionName, _ := syscall.UTF16PtrFromString(distributionName)

	return _WslConfigureDistribution(pDistributionName, defaultUID, wslDistributionFlags)
}

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

//sys	_WslLaunch(distributionName *uint16, command *uint16, useCurrentWorkingDirectory bool, stdIn syscall.Handle, stdOut syscall.Handle, stdErr syscall.Handle, process *syscall.Handle) (err error) = wslapi.WslLaunch

// WslLaunch launches the distribution with handle
func WslLaunch(distributionName string, command string, useCurrentWorkingDirectory bool, stdIn syscall.Handle, stdOut syscall.Handle, stdErr syscall.Handle) (process syscall.Handle, err error) {
	pDistributionName, _ := syscall.UTF16PtrFromString(distributionName)
	pCommand, _ := syscall.UTF16PtrFromString(command)

	_WslLaunch(pDistributionName, pCommand, useCurrentWorkingDirectory, stdIn, stdOut, stdErr, &process)
	return
}

// WslLaunchInteractive launches the distribution with interactive shell
func WslLaunchInteractive(distributionName string, command string, useCurrentWorkingDirectory bool) (exitCode uint32, err error) {
	p, _ := syscall.GetCurrentProcess()
	stdin := syscall.Handle(0)
	stdout := syscall.Handle(0)
	stderr := syscall.Handle(0)

	syscall.DuplicateHandle(p, syscall.Handle(os.Stdin.Fd()), p, &stdin, 0, true, syscall.DUPLICATE_SAME_ACCESS)
	syscall.DuplicateHandle(p, syscall.Handle(os.Stdout.Fd()), p, &stdout, 0, true, syscall.DUPLICATE_SAME_ACCESS)
	syscall.DuplicateHandle(p, syscall.Handle(os.Stderr.Fd()), p, &stderr, 0, true, syscall.DUPLICATE_SAME_ACCESS)

	handle, err := WslLaunch(distributionName, command, useCurrentWorkingDirectory, stdin, stdout, stderr)
	syscall.WaitForSingleObject(handle, syscall.INFINITE)
	syscall.GetExitCodeProcess(handle, &exitCode)
	return
}

//sys	_WslUnregisterDistribution(distributionName *uint16) (res error) = wslapi.WslUnregisterDistribution

// WslUnregisterDistribution unregisters the specified distribution
func WslUnregisterDistribution(distributionName string) error {
	pDistributionName, _ := syscall.UTF16PtrFromString(distributionName)
	return _WslUnregisterDistribution(pDistributionName)
}
