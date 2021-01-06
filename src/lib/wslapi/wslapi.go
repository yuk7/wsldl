package wslapi

import (
	"syscall"
	"unsafe"
)

var (
	modwslapi = syscall.NewLazyDLL("wslapi.dll")

	procWslConfigureDistribution        = modwslapi.NewProc("procWslConfigureDistribution")
	procWslGetDistributionConfiguration = modwslapi.NewProc("WslGetDistributionConfiguration")
	procWslIsDistributionRegistered     = modwslapi.NewProc("WslIsDistributionRegistered")
	procWslLaunch                       = modwslapi.NewProc("WslLaunch")
	procWslRegisterDistribution         = modwslapi.NewProc("WslRegisterDistribution")
	procWslUnregisterDistribution       = modwslapi.NewProc("WslUnregisterDistribution")
)

//RequestError is used for return detailed error
type RequestError struct {
	ErrorCode uint
	Err       error
}

func _WslIsDistributionRegistered(distributionName *uint16) bool {
	r0, _, _ := syscall.Syscall(procWslIsDistributionRegistered.Addr(), 1, uintptr(unsafe.Pointer(distributionName)), 0, 0)

	if r0 == 0 {
		return false
	}
	return true
}

// WslIsDistributionRegistered determines if a distribution is already registered
func WslIsDistributionRegistered(distributionName string) bool {
	pDistributionName, _ := syscall.UTF16PtrFromString(distributionName)
	return _WslIsDistributionRegistered(pDistributionName)
}

func _WslRegisterDistribution(distributionName, tarGzFilename *uint16) uint {
	r0, _, _ := syscall.Syscall(procWslRegisterDistribution.Addr(), 2, uintptr(unsafe.Pointer(distributionName)), uintptr(unsafe.Pointer(tarGzFilename)), 0)

	return uint(r0)
}

// WslRegisterDistribution registers a new distribution
func WslRegisterDistribution(distributionName, tarGzFilename string) uint {
	pDistributionName, _ := syscall.UTF16PtrFromString(distributionName)
	pTarGzFilename, _ := syscall.UTF16PtrFromString(tarGzFilename)
	return _WslRegisterDistribution(pDistributionName, pTarGzFilename)
}
