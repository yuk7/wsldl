package repair

import (
	"os"
	"path/filepath"
)

// IsInstalledFilesExist checks if the executable directory has either WSL1 or WSL2 install files.
func IsInstalledFilesExist() bool {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)

	_, err := os.Stat(dir + "\\ext4.vhdx")
	if err == nil {
		return true
	}
	_, err = os.Stat(dir + "\\rootfs")
	return err == nil
}
