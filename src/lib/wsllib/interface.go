package wsllib

import (
	"syscall"

	wsllibgo "github.com/yuk7/wsllib-go"
	wslreg "github.com/yuk7/wslreglib-go"
)

type Profile = wslreg.Profile

const (
	FlagAppendNTPath        = wsllibgo.FlagAppendNTPath
	FlagEnableDriveMounting = wsllibgo.FlagEnableDriveMounting
	FlagEnableWsl2          = wsllibgo.FlagEnableWsl2
	FlagWsldlTermDefault    = wslreg.FlagWsldlTermDefault
	FlagWsldlTermWT         = wslreg.FlagWsldlTermWT
	FlagWsldlTermFlute      = wslreg.FlagWsldlTermFlute
	LxssBaseKey             = wslreg.LxssBaseKey
)

type WslLib interface {
	IsDistributionRegistered(name string) bool
	RegisterDistribution(name, rootPath string) error
	UnregisterDistribution(name string) error
	LaunchInteractive(name, command string, inheritPath bool) (uint32, error)
	Launch(name, command string, inheritPath bool, stdin, stdout, stderr syscall.Handle) (syscall.Handle, error)
	GetDistributionConfiguration(name string) (uint32, uint64, uint32, error)
	ConfigureDistribution(name string, uid uint64, flags uint32) error
}

type WslReg interface {
	GetProfileFromName(name string) (Profile, error)
	GetProfileFromBasePath(path string) (Profile, error)
	WriteProfile(profile Profile) error
	SetWslVersion(name string, version int) error
	GenerateProfile() Profile
}

type Dependencies struct {
	Wsl WslLib
	Reg WslReg
}

func NewDependencies() Dependencies {
	return Dependencies{
		Wsl: NewNativeWslLib(),
		Reg: NewNativeWslReg(),
	}
}
