package wslreg

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/sys/windows/registry"
)

const (
	// LxssBaseRoot is CURRENT_USER
	LxssBaseRoot = registry.CURRENT_USER
	// LxssBaseRootStr is CURRENT_USER string
	LxssBaseRootStr = "HKEY_CURRENT_USER"
	// LxssBaseKey is path of lxss registry
	LxssBaseKey = "Software\\Microsoft\\Windows\\CurrentVersion\\Lxss"
	// WsldlTermKey is registry key name used for wsldl terminal infomation
	WsldlTermKey = "wsldl-term"
	// FlagWsldlTermDefault is default terminal (conhost)
	FlagWsldlTermDefault = 0
	// FlagWsldlTermWT is Windows Terminal
	FlagWsldlTermWT = 1
	// FlagWsldlTermFlute is Fluent Terminal
	FlagWsldlTermFlute = 2
	// InvalidNum is Num used for invalid
	InvalidNum = -1
)

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
	if err != nil {
		key = 0
	}
	regpathStr := LxssBaseRootStr + "\\" + LxssBaseKey + "\\" + profile.UUID

	if profile.BasePath != "" {
		err = regSetStringWithCmdAndFix(key, regpathStr, "BasePath", profile.BasePath)
		if err != nil {
			return err
		}
	}
	if profile.DistributionName != "" {
		err = regSetStringWithCmdAndFix(key, regpathStr, "DistributionName", profile.DistributionName)
		if err != nil {
			return err
		}
	}
	if profile.DefaultUid != InvalidNum {
		err = regSetDWordWithCmdAndFix(key, regpathStr, "DefaultUid", uint32(profile.DefaultUid))
		if err != nil {
			return err
		}
	}

	if profile.Flags != InvalidNum {
		err = regSetDWordWithCmdAndFix(key, regpathStr, "Flags", uint32(profile.Flags))
		if err != nil {
			return err
		}
	}

	if profile.State != InvalidNum {
		err = regSetDWordWithCmdAndFix(key, regpathStr, "State", uint32(profile.State))
		if err != nil {
			return err
		}
	}

	if profile.Version != InvalidNum {
		err = regSetDWordWithCmdAndFix(key, regpathStr, "Version", uint32(profile.Version))
		if err != nil {
			return err
		}
	}

	if profile.PackageFamilyName != "" {
		err = regSetStringWithCmdAndFix(key, regpathStr, "PackageFamilyName", profile.PackageFamilyName)
		if err != nil {
			return err
		}
	}
	if profile.WsldlTerm != InvalidNum {
		err = regSetDWordWithCmdAndFix(key, regpathStr, WsldlTermKey, uint32(profile.WsldlTerm))
		if err != nil {
			return err
		}
	}
	return nil
}

func regSetStringWithCmdAndFix(regkey registry.Key, regpathStr, keyname, value string) error {
	// backup oldValue
	oldValue := ""
	if regkey != 0 {
		oldValue, _, _ = regkey.GetStringValue(value)
	}

	// write with external command
	err := regSetStringWithCmd(regpathStr, keyname, value)
	if err != nil {
		return err
	}

	// if regkey not 0, check if the value was written correctly
	if regkey != 0 && oldValue != "" && oldValue != value {
		newValue, _, _ := regkey.GetStringValue(value)
		if oldValue == newValue {
			// not changed, maybe appx virtual registry used
			// delete it and rewrite registry value
			regkey.DeleteValue(keyname)
			err := regSetStringWithCmd(regpathStr, keyname, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// regSetStringWithCmd writes DWord key with command, forcibly use the real registry
func regSetStringWithCmd(regpath, keyname, value string) error {
	regexe := os.Getenv("SystemRoot") + "\\System32\\reg.exe"

	_, err := exec.Command(regexe, "add", regpath, "/v", keyname, "/t", "REG_SZ", "/d", value, "/f").Output()
	return err
}

func regSetDWordWithCmdAndFix(regkey registry.Key, regpathStr, keyname string, value uint32) error {
	// backup oldValue
	oldValue := InvalidNum
	if regkey != 0 {
		val, _, _ := regkey.GetIntegerValue(keyname)
		oldValue = int(val)
	}

	// write with external command
	err := regSetDWordWithCmd(regpathStr, keyname, value)
	if err != nil {
		return err
	}

	// if regkey not 0, check if the value was written correctly
	if regkey != 0 && oldValue != InvalidNum && oldValue != int(value) {
		newValue, _, _ := regkey.GetIntegerValue(keyname)
		if oldValue == int(newValue) {
			// not changed, maybe appx virtual registry used
			// delete it and rewrite registry value
			regkey.DeleteValue(keyname)
			err := regSetDWordWithCmd(regpathStr, keyname, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// regSetDWordWithCmd writes DWord key with command, forcibly use the real registry
func regSetDWordWithCmd(regpath, keyname string, value uint32) error {
	regexe := os.Getenv("SystemRoot") + "\\System32\\reg.exe"

	_, err := exec.Command(regexe, "add", regpath, "/v", keyname, "/t", "REG_DWORD", "/d", strconv.Itoa(int(value)), "/f").Output()
	return err
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

// GetLxUuidList gets guid key lists
func GetLxUuidList() (uuidList []string, err error) {
	baseKey, tmpErr := registry.OpenKey(LxssBaseRoot, LxssBaseKey, registry.READ)
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	uuidList, tmpErr = baseKey.ReadSubKeyNames(1024)
	if tmpErr != nil && tmpErr != io.EOF {
		err = tmpErr
		return
	}
	return
}

// GetProfileFromName gets distro profile from name
func GetProfileFromName(distributionName string) (profile Profile, err error) {
	uuidList, tmpErr := GetLxUuidList()
	if tmpErr != nil {
		err = tmpErr
		return
	}

	errStr := ""
	for _, loopUUID := range uuidList {
		profile, _ = ReadProfile(loopUUID)
		if strings.EqualFold(profile.DistributionName, distributionName) {
			return
		}
	}
	err = errors.New("Registry Key Not found\n" + errStr)
	profile = NewProfile()
	return
}

// GetProfileFromBasePath gets distro profile from BasePath
func GetProfileFromBasePath(basePath string) (profile Profile, err error) {
	uuidList, tmpErr := GetLxUuidList()
	if tmpErr != nil {
		err = tmpErr
		return
	}

	basePathAbs, tmpErr := filepath.Abs(basePath)
	if err != nil {
		basePathAbs = basePath
	}

	errStr := ""
	for _, loopUUID := range uuidList {
		profile, _ = ReadProfile(loopUUID)
		if strings.EqualFold(profile.BasePath, basePathAbs) {
			return
		}
	}
	err = errors.New("Registry Key Not found\n" + errStr)
	profile = NewProfile()
	return
}

// SetWslVersion sets wsl version
func SetWslVersion(distributionName string, version int) error {
	wslexe := os.Getenv("SystemRoot") + "\\System32\\wsl.exe"
	_, err := exec.Command(wslexe, "--set-version", distributionName, strconv.Itoa(version)).Output()
	return err
}
