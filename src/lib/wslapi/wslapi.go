package wslapi

import (
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
