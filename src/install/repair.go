package install

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/yuk7/wsldl/lib/wsllib"
)

func isInstalledFilesExist() bool {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)

	_, err := os.Stat(dir + "\\ext4.vhdx")
	if err == nil {
		return true
	}
	_, err = os.Stat(dir + "\\rootfs")
	return err == nil
}

func repairRegistry(reg wsllib.WslReg, name string) error {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)

	// for rename instance
	prof, _ := reg.GetProfileFromBasePath(dir)
	if prof.BasePath != "" {
		// profile found, maybe executable renamed
		// write the new name to the registry
		prof.DistributionName = name
		return reg.WriteProfile(prof)
	}

	// for write new WSL2 configuration
	_, err := os.Stat(dir + "\\ext4.vhdx")
	if err == nil {
		prof := reg.GenerateProfile()
		prof.DistributionName = name
		prof.BasePath = dir

		prof.Flags |= wsllib.FlagEnableWsl2
		return reg.WriteProfile(prof)
	}
	// for write new WSL1 configuration
	_, err = os.Stat(dir + "\\rootfs")
	if err == nil {
		prof := reg.GenerateProfile()
		prof.DistributionName = name
		prof.BasePath = dir
		prof.Flags ^= wsllib.FlagEnableWsl2
		return reg.WriteProfile(prof)
	}

	return errors.New("repair failed")
}
