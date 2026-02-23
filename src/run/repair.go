package run

import (
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

func repairRegistry(reg wsllib.WslReg, profile wsllib.Profile) error {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)

	profile.BasePath = dir
	return reg.WriteProfile(profile)
}
