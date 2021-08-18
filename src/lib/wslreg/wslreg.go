package wslreg

import (
	"errors"
	"io"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/sys/windows/registry"
)

const (
	// LxssBaseRoot is LOCAL_MACHINE
	LxssBaseRoot = registry.CURRENT_USER
	// LxssBaseKey is path of lxss registry
	LxssBaseKey = "Software\\Microsoft\\Windows\\CurrentVersion\\Lxss"
	// WsldlTermKey is registry key name used for wsldl terminal infomation
	WsldlTermKey = "wsldl-term"
	// InvalidNum is Num used for invalid
	InvalidNum = -1
)

// Profile is profile for WSL
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

// NewProfile creates empty profile
func NewProfile() Profile {
	profile := Profile{}
	profile.DefaultUid = InvalidNum
	profile.Flags = InvalidNum
	profile.State = InvalidNum
	profile.Version = InvalidNum
	profile.WsldlTerm = InvalidNum
	return profile
}

// GenerateProfile generates new profile with UUID
func GenerateProfile() Profile {
	profile := NewProfile()
	profile.UUID = "{" + uuid.NewV4().String() + "}"
	profile.DefaultUid = 0x0
	profile.Flags = 0x7
	profile.State = 0x1
	profile.Version = 0x2
	profile.WsldlTerm = 0x0
	return profile
}

func WriteProfile(profile Profile) error {
	if profile.UUID == "" {
		return errors.New("Empty UUID")
	}
	key, _, err := registry.CreateKey(LxssBaseRoot, LxssBaseKey+"\\"+profile.UUID, registry.ALL_ACCESS)
	if err != nil && err != io.EOF {
		return err
	}
	if profile.BasePath != "" {
		err = key.SetStringValue("BasePath", profile.BasePath)
		if err != nil && err != io.EOF {
			return err
		}
	}
	if profile.DistributionName != "" {
		err = key.SetStringValue("DistributionName", profile.DistributionName)
		if err != nil && err != io.EOF {
			return err
		}
	}
	if profile.DefaultUid != InvalidNum {
		err = key.SetDWordValue("DefaultUid", uint32(profile.DefaultUid))
		if err != nil && err != io.EOF {
			return err
		}
	}

	if profile.Flags != InvalidNum {
		err = key.SetDWordValue("Flags", uint32(profile.Flags))
		if err != nil && err != io.EOF {
			return err
		}
	}

	if profile.State != InvalidNum {
		err = key.SetDWordValue("State", uint32(profile.State))
		if err != nil && err != io.EOF {
			return err
		}
	}

	if profile.Version != InvalidNum {
		err = key.SetDWordValue("Version", uint32(profile.Version))
		if err != nil && err != io.EOF {
			return err
		}
	}

	if profile.PackageFamilyName != "" {
		err = key.SetStringValue("PackageFamilyName", profile.PackageFamilyName)
		if err != nil && err != io.EOF {
			return err
		}
	}
	if profile.WsldlTerm != InvalidNum {
		err = key.SetDWordValue(WsldlTermKey, uint32(profile.WsldlTerm))
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func ReadProfile(lxUuid string) (profile Profile, err error) {
	profile = NewProfile()
	profile.UUID = lxUuid
	key, err := registry.OpenKey(LxssBaseRoot, LxssBaseKey+"\\"+profile.UUID, registry.READ)
	if err != nil {
		return
	}
	basepath, _, tmperr := key.GetStringValue("BasePath")
	if tmperr == nil || tmperr == io.EOF {
		profile.BasePath = basepath
	}
	distributionName, _, tmperr := key.GetStringValue("DistributionName")
	if tmperr == nil || tmperr == io.EOF {
		profile.DistributionName = distributionName
	}
	flags, _, tmperr := key.GetIntegerValue("Flags")
	if tmperr == nil || tmperr == io.EOF {
		profile.Flags = int(flags)
	}
	state, _, tmperr := key.GetIntegerValue("State")
	if tmperr == nil || tmperr == io.EOF {
		profile.State = int(state)
	}
	version, _, tmperr := key.GetIntegerValue("Version")
	if tmperr == nil || tmperr == io.EOF {
		profile.Version = int(version)
	}
	wsldlTerm, _, tmperr := key.GetIntegerValue(WsldlTermKey)
	if tmperr == nil || tmperr == io.EOF {
		profile.WsldlTerm = int(wsldlTerm)
	}
	pkgName, _, tmperr := key.GetStringValue("PackageFamilyName")
	if tmperr == nil || tmperr == io.EOF {
		profile.PackageFamilyName = pkgName
	}
	return
}
