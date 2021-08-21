package install

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuk7/wsldl/lib/wslapi"
	"github.com/yuk7/wsldl/lib/wslreg"
)

var (
	defaultRootFiles = []string{"install.tar", "install.tar.gz", "rootfs.tar", "rootfs.tar.gz", "install.ext4.vhdx"}
)

//Install installs distribution with default rootfs file names
func Install(name string, rootPath string, showProgress bool) error {
	if showProgress {
		fmt.Printf("Using: %s\n", rootPath)
		fmt.Println("Installing...")
	}
	rootPathLower := strings.ToLower(rootPath)
	if strings.HasSuffix(rootPathLower, "ext4.vhdx") {
		return InstallExt4Vhdx(name, rootPath)
	}
	return InstallTar(name, rootPath)
}

func InstallTar(name string, rootPath string) error {
	err := wslapi.WslRegisterDistribution(name, rootPath)
	return err
}

func InstallExt4Vhdx(name string, rootPath string) error {
	// create empty tar
	tmptar := os.TempDir()
	if tmptar == "" {
		return errors.New("Failed to create temp directory")
	}
	tmptar = tmptar + "\\em-vhdx-temp.tar"
	tmptarfp, err := os.Create(tmptar)
	if err != nil {
		return err
	}
	tmptarfp.Close()
	// initial empty instance entry
	err = wslapi.WslRegisterDistribution(name, tmptar)
	if err != nil {
		return err
	}
	os.Remove(tmptar)
	// get profile of instance
	prof, err := wslreg.GetProfileFromName(name)
	if prof.BasePath == "" {
		return err
	}
	// remove instance temporary
	err = wslapi.WslUnregisterDistribution(name)
	if err != nil {
		return err
	}
	// copy vhdx to destination directory
	src, err := os.Open(rootPath)
	if err != nil {
		return err
	}
	defer src.Close()
	dest, err := os.Create(prof.BasePath + "\\ext4.vhdx")
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	// write registry
	prof.Flags |= wslapi.FlagEnableWsl2
	err = wslreg.WriteProfile(prof)
	return err
}

func detectRootfsFiles() string {
	efPath, _ := os.Executable()
	efDir := filepath.Dir(efPath)
	for _, rootFile := range defaultRootFiles {
		rootPath := filepath.Join(efDir, rootFile)
		_, err := os.Stat(rootPath)
		if err == nil {
			return rootPath
		}
	}
	return "rootfs.tar.gz"
}
