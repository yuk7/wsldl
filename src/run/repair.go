package run

import (
	"os"
	"path/filepath"

	"github.com/yuk7/wsldl/lib/wslreg"
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

func repairRegistry(profile wslreg.Profile) error {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)

	profile.BasePath = dir
	return wslreg.WriteProfile(profile)
}
