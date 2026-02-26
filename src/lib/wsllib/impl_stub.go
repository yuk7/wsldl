//go:build !windows

package wsllib

import "errors"

var errUnsupportedPlatform = errors.New("wsllib native implementation is only available on Windows")

type nativeWslLib struct{}

func NewNativeWslLib() WslLib {
	return nativeWslLib{}
}

func (nativeWslLib) IsDistributionRegistered(name string) bool {
	return false
}

func (nativeWslLib) RegisterDistribution(name, rootPath string) error {
	return errUnsupportedPlatform
}

func (nativeWslLib) UnregisterDistribution(name string) error {
	return errUnsupportedPlatform
}

func (nativeWslLib) LaunchInteractive(name, command string, inheritPath bool) (uint32, error) {
	return 0, errUnsupportedPlatform
}

func (nativeWslLib) Launch(name, command string, inheritPath bool, stdin, stdout, stderr Handle) (Handle, error) {
	return Handle(0), errUnsupportedPlatform
}

func (nativeWslLib) GetDistributionConfiguration(name string) (uint32, uint64, uint32, error) {
	return 0, 0, 0, errUnsupportedPlatform
}

func (nativeWslLib) ConfigureDistribution(name string, uid uint64, flags uint32) error {
	return errUnsupportedPlatform
}

type nativeWslReg struct{}

func NewNativeWslReg() WslReg {
	return nativeWslReg{}
}

func (nativeWslReg) GetProfileFromName(name string) (Profile, error) {
	return Profile{}, errUnsupportedPlatform
}

func (nativeWslReg) GetProfileFromBasePath(path string) (Profile, error) {
	return Profile{}, errUnsupportedPlatform
}

func (nativeWslReg) WriteProfile(profile Profile) error {
	return errUnsupportedPlatform
}

func (nativeWslReg) SetWslVersion(name string, version int) error {
	return errUnsupportedPlatform
}

func (nativeWslReg) GenerateProfile() Profile {
	return Profile{}
}
