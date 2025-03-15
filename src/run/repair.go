package run

import (
	"os"
	"path/filepath"

	wslreg "github.com/yuk7/wslreglib-go"
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

func repairRegistry(profile wslreg.Profile) error {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)

	profile.BasePath = dir
	return wslreg.WriteProfile(profile)
}
