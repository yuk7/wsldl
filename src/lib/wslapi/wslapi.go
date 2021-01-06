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
