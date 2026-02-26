package wsllib

import (
	"os"
	"path/filepath"
	"strings"
)

// Profile is registry profile information used by wsldl.
type Profile struct {
	UUID              string
	BasePath          string
	DistributionName  string
	DefaultUid        int
	Flags             int
	State             int
	Version           int
	PackageFamilyName string
	WsldlTerm         int
}

const (
	FlagAppendNTPath        = 2
	FlagEnableDriveMounting = 4
	FlagEnableWsl2          = 8
	FlagWsldlTermDefault    = 0
	FlagWsldlTermWT         = 1
	FlagWsldlTermFlute      = 2
	LxssBaseKey             = "Software\\Microsoft\\Windows\\CurrentVersion\\Lxss"
)

type WslLib interface {
	IsDistributionRegistered(name string) bool
	RegisterDistribution(name, rootPath string) error
	UnregisterDistribution(name string) error
	LaunchInteractive(name, command string, inheritPath bool) (uint32, error)
	Launch(name, command string, inheritPath bool, stdin, stdout, stderr Handle) (Handle, error)
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
	if isUnitTestProcess() {
		return Dependencies{
			Wsl: MockWslLib{},
			Reg: MockWslReg{},
		}
	}

	return Dependencies{
		Wsl: NewNativeWslLib(),
		Reg: NewNativeWslReg(),
	}
}

func isUnitTestProcess() bool {
	if strings.HasSuffix(filepath.Base(os.Args[0]), ".test") {
		return true
	}
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}
