//go:build windows

package wsllib

import (
	"syscall"

	wsllibgo "github.com/yuk7/wsllib-go"
	wslreg "github.com/yuk7/wslreglib-go"
)

type nativeWslLib struct{}

func NewNativeWslLib() WslLib {
	return nativeWslLib{}
}

func (nativeWslLib) IsDistributionRegistered(name string) bool {
	return wsllibgo.WslIsDistributionRegistered(name)
}

func (nativeWslLib) RegisterDistribution(name, rootPath string) error {
	return wsllibgo.WslRegisterDistribution(name, rootPath)
}

func (nativeWslLib) UnregisterDistribution(name string) error {
	return wsllibgo.WslUnregisterDistribution(name)
}

func (nativeWslLib) LaunchInteractive(name, command string, inheritPath bool) (uint32, error) {
	return wsllibgo.WslLaunchInteractive(name, command, inheritPath)
}

func (nativeWslLib) Launch(name, command string, inheritPath bool, stdin, stdout, stderr Handle) (Handle, error) {
	return wsllibgo.WslLaunch(name, command, inheritPath, syscall.Handle(stdin), syscall.Handle(stdout), syscall.Handle(stderr))
}

func (nativeWslLib) GetDistributionConfiguration(name string) (uint32, uint64, uint32, error) {
	return wsllibgo.WslGetDistributionConfiguration(name)
}

func (nativeWslLib) ConfigureDistribution(name string, uid uint64, flags uint32) error {
	return wsllibgo.WslConfigureDistribution(name, uid, flags)
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
	p, err := wslreg.GetProfileFromName(name)
	return toProfile(p), err
}

func (nativeWslReg) GetProfileFromBasePath(path string) (Profile, error) {
	p, err := wslreg.GetProfileFromBasePath(path)
	return toProfile(p), err
}

func (nativeWslReg) WriteProfile(profile Profile) error {
	return wslreg.WriteProfile(fromProfile(profile))
}

func (nativeWslReg) SetWslVersion(name string, version int) error {
	return wslreg.SetWslVersion(name, version)
}

func (nativeWslReg) GenerateProfile() Profile {
	return toProfile(wslreg.GenerateProfile())
}
