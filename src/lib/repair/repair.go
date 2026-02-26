package repair

import (
	"os"
	"path/filepath"
)

// IsInstalledFilesExist checks if the executable directory has either WSL1 or WSL2 install files.
func IsInstalledFilesExist() bool {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)
	return isInstalledFilesExistInDir(dir, os.Stat)
}

func isInstalledFilesExistInDir(dir string, stat func(name string) (os.FileInfo, error)) bool {
	_, err := stat(dir + "\\ext4.vhdx")
	if err == nil {
		return true
	}
	_, err = stat(dir + "\\rootfs")
	return err == nil
}
