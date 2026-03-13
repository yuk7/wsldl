//go:build windows

package wsllib

import (
	"syscall"

	wsllibgo "github.com/yuk7/wsllib-go"
	wslreg "github.com/yuk7/wslreglib-go"
)

var (
	wslIsDistributionRegisteredFunc     = wsllibgo.WslIsDistributionRegistered
	wslRegisterDistributionFunc         = wsllibgo.WslRegisterDistribution
	wslUnregisterDistributionFunc       = wsllibgo.WslUnregisterDistribution
	wslLaunchInteractiveFunc            = wsllibgo.WslLaunchInteractive
	wslLaunchFunc                       = wsllibgo.WslLaunch
	wslGetDistributionConfigurationFunc = wsllibgo.WslGetDistributionConfiguration
	wslConfigureDistributionFunc        = wsllibgo.WslConfigureDistribution
	regGetProfileFromNameFunc           = wslreg.GetProfileFromName
	regGetProfileFromBasePathFunc       = wslreg.GetProfileFromBasePath
	regWriteProfileFunc                 = wslreg.WriteProfile
	regSetWslVersionFunc                = wslreg.SetWslVersion
	regGenerateProfileFunc              = wslreg.GenerateProfile
)

type nativeWslLib struct{}

func NewNativeWslLib() WslLib {
	return nativeWslLib{}
}

func (nativeWslLib) IsDistributionRegistered(name string) bool {
	return wslIsDistributionRegisteredFunc(name)
}

func (nativeWslLib) RegisterDistribution(name, rootPath string) error {
	return wslRegisterDistributionFunc(name, rootPath)
}

func (nativeWslLib) UnregisterDistribution(name string) error {
	return wslUnregisterDistributionFunc(name)
}

func (nativeWslLib) LaunchInteractive(name, command string, inheritPath bool) (uint32, error) {
	return wslLaunchInteractiveFunc(name, command, inheritPath)
}

func (nativeWslLib) Launch(name, command string, inheritPath bool, stdin, stdout, stderr Handle) (Handle, error) {
	return wslLaunchFunc(name, command, inheritPath, syscall.Handle(stdin), syscall.Handle(stdout), syscall.Handle(stderr))
}

func (nativeWslLib) GetDistributionConfiguration(name string) (uint32, uint64, uint32, error) {
	return wslGetDistributionConfigurationFunc(name)
}

func (nativeWslLib) ConfigureDistribution(name string, uid uint64, flags uint32) error {
	return wslConfigureDistributionFunc(name, uid, flags)
}

type nativeWslReg struct{}

func NewNativeWslReg() WslReg {
	return nativeWslReg{}
}

func toProfile(p wslreg.Profile) Profile {
	return Profile{
		UUID:              p.UUID,
		BasePath:          p.BasePath,
		DistributionName:  p.DistributionName,
		DefaultUid:        p.DefaultUid,
		Flags:             p.Flags,
		State:             p.State,
		Version:           p.Version,
		PackageFamilyName: p.PackageFamilyName,
		WsldlTerm:         p.WsldlTerm,
	}
}

func fromProfile(p Profile) wslreg.Profile {
	return wslreg.Profile{
		UUID:              p.UUID,
		BasePath:          p.BasePath,
		DistributionName:  p.DistributionName,
		DefaultUid:        p.DefaultUid,
		Flags:             p.Flags,
		State:             p.State,
		Version:           p.Version,
		PackageFamilyName: p.PackageFamilyName,
		WsldlTerm:         p.WsldlTerm,
	}
}

func (nativeWslReg) GetProfileFromName(name string) (Profile, error) {
	p, err := regGetProfileFromNameFunc(name)
	return toProfile(p), err
}

func (nativeWslReg) GetProfileFromBasePath(path string) (Profile, error) {
	p, err := regGetProfileFromBasePathFunc(path)
	return toProfile(p), err
}

func (nativeWslReg) WriteProfile(profile Profile) error {
	return regWriteProfileFunc(fromProfile(profile))
}

func (nativeWslReg) SetWslVersion(name string, version int) error {
	return regSetWslVersionFunc(name, version)
}

func (nativeWslReg) GenerateProfile() Profile {
	return toProfile(regGenerateProfileFunc())
}
