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

func (nativeWslLib) Launch(name, command string, inheritPath bool, stdin, stdout, stderr syscall.Handle) (syscall.Handle, error) {
	return wsllibgo.WslLaunch(name, command, inheritPath, stdin, stdout, stderr)
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

func (nativeWslReg) GetProfileFromName(name string) (Profile, error) {
	return wslreg.GetProfileFromName(name)
}

func (nativeWslReg) GetProfileFromBasePath(path string) (Profile, error) {
	return wslreg.GetProfileFromBasePath(path)
}

func (nativeWslReg) WriteProfile(profile Profile) error {
	return wslreg.WriteProfile(profile)
}

func (nativeWslReg) SetWslVersion(name string, version int) error {
	return wslreg.SetWslVersion(name, version)
}

func (nativeWslReg) GenerateProfile() Profile {
	return wslreg.GenerateProfile()
}
