package install

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/yuk7/wsldl/lib/wslreg"
	"github.com/yuk7/wsllib-go"
)

func isInstalledFilesExist() bool {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)

	_, err := os.Stat(dir + "\\ext4.vhdx")
	if err == nil {
		return true
	}
	_, err = os.Stat(dir + "\\rootfs")
	if err == nil {
		return true
	}
	return false
}

func repairRegistry(name string) error {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)

	// for rename instance
	prof, _ := wslreg.GetProfileFromBasePath(dir)
	if prof.BasePath != "" {
		// profile found, maybe executable renamed
		// write the new name to the registry
		prof.DistributionName = name
		return wslreg.WriteProfile(prof)
	}

	// for write new WSL2 configuration
	_, err := os.Stat(dir + "\\ext4.vhdx")
	if err == nil {
		prof := wslreg.GenerateProfile()
		prof.DistributionName = name
		prof.BasePath = dir

		prof.Flags |= wsllib.FlagEnableWsl2
		return wslreg.WriteProfile(prof)
	}
	// for write new WSL1 configuration
	_, err = os.Stat(dir + "\\rootfs")
	if err == nil {
		prof := wslreg.GenerateProfile()
		prof.DistributionName = name
		prof.BasePath = dir
		prof.Flags ^= wsllib.FlagEnableWsl2
		return wslreg.WriteProfile(prof)
	}

	return errors.New("Repair failed")
}
