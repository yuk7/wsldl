package install

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
	wslreg "github.com/yuk7/wslreglib-go"
)

var (
	defaultRootFiles = []string{
		"install.tar",
		"install.tar.gz",
		"install.tgz",
		"install.tar.zst",
		"install.tar.xz",
		"rootfs.tar",
		"rootfs.tar.gz",
		"rootfs.tgz",
		"rootfs.tar.zst",
		"rootfs.tar.xz",
		"install.ext4.vhdx",
		"install.ext4.vhdx.gz",
	}
)

// Install installs distribution with default rootfs file names
func Install(name string, rootPath string, showProgress bool) error {
	rootPathLower := strings.ToLower(rootPath)
	if showProgress {
		fmt.Printf("Using: %s\n", rootPath)
	}
	if strings.HasPrefix(rootPathLower, "http://") || strings.HasPrefix(rootPathLower, "https://") {
		progressBarWidth := 0
		if showProgress {
			fmt.Println("Downloading...")
			progressBarWidth = 35
		}
		tmpRootFn := os.TempDir()
		if tmpRootFn == "" {
			return errors.New("Failed to create temp directory")
		}
		rand.Seed(time.Now().UnixNano())
		tmpRootFn = tmpRootFn + "\\" + strconv.Itoa(rand.Intn(10000)) + filepath.Base(rootPath)
		defer os.Remove(tmpRootFn)
		err := utils.DownloadFile(rootPath, tmpRootFn, progressBarWidth)
		if err != nil {
			return err
		}
		rootPath = tmpRootFn
		rootPathLower = strings.ToLower(rootPath)
		fmt.Println()
	}
	if showProgress {
		fmt.Println("Installing...")
	}

	if strings.HasSuffix(rootPathLower, "ext4.vhdx") || strings.HasSuffix(rootPathLower, "ext4.vhdx.gz") {
		return InstallExt4Vhdx(name, rootPath)
	}
	return InstallTar(name, rootPath)
}

func InstallTar(name string, rootPath string) error {
	err := wsllib.WslRegisterDistribution(name, rootPath)
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
	err = wsllib.WslRegisterDistribution(name, tmptar)
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
	err = wsllib.WslUnregisterDistribution(name)
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

	// uncompress and copy
	rootPathLower := strings.ToLower(rootPath)
	if strings.HasSuffix(rootPathLower, ".gz") {
		// compressed with gzip
		gr, err := gzip.NewReader(src)
		if err != nil {
			return err
		}
		_, err = io.Copy(dest, gr)
		if err != nil {
			return err
		}
	} else {
		// not compressed
		_, err = io.Copy(dest, src)
		if err != nil {
			return err
		}
	}

	// write registry
	prof.Flags |= wsllib.FlagEnableWsl2
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
