package install

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yuk7/wsldl/lib/download"
	"github.com/yuk7/wsldl/lib/fileutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

var (
	defaultRootFiles = []string{
		"install.tar",
		"install.tar.gz",
		"install.tgz",
		"install.tar.zst",
		"install.tar.xz",
		"install.wsl",
		"rootfs.tar",
		"rootfs.tar.gz",
		"rootfs.tgz",
		"rootfs.tar.zst",
		"rootfs.tar.xz",
		"rootfs.wsl",
		"install.ext4.vhdx",
		"install.ext4.vhdx.gz",
	}
)

// Install installs distribution with default rootfs file names
func Install(wsl wsllib.WslLib, reg wsllib.WslReg, name string, rootPath string, sha256Sum string, showProgress bool) error {
	rootPathLower := strings.ToLower(rootPath)
	sha256Actual := ""
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
			return errors.New("failed to create temp directory")
		}
		rand.NewSource(time.Now().UnixNano())
		tmpRootFn = tmpRootFn + "\\" + strconv.Itoa(rand.Intn(10000)) + filepath.Base(rootPath)
		defer os.Remove(tmpRootFn)
		var err error
		sha256Actual, err = download.DownloadFile(rootPath, tmpRootFn, progressBarWidth)
		if err != nil {
			return err
		}
		rootPath = tmpRootFn
		rootPathLower = strings.ToLower(rootPath)
		fmt.Println()
	} else if sha256Sum != "" {
		if showProgress {
			fmt.Println("Calculating checksum...")
		}
		f, err := os.Open(rootPath)
		if err != nil {
			return err
		}
		defer f.Close()
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		sha256Actual = hex.EncodeToString(h.Sum(nil))
	}

	if showProgress && sha256Actual != "" {
		fmt.Printf("Checksum(SHA256): %s\n", sha256Actual)
	}

	if sha256Sum != "" && sha256Actual != "" && sha256Sum != sha256Actual {
		return errors.New("checksum mismatch")
	}

	if showProgress {
		fmt.Println("Installing...")
	}

	if strings.HasSuffix(rootPathLower, "ext4.vhdx") || strings.HasSuffix(rootPathLower, "ext4.vhdx.gz") {
		return InstallExt4Vhdx(wsl, reg, name, rootPath)
	}
	return InstallTar(wsl, name, rootPath)
}

func InstallTar(wsl wsllib.WslLib, name string, rootPath string) error {
	err := wsl.RegisterDistribution(name, rootPath)
	return err
}

func InstallExt4Vhdx(wsl wsllib.WslLib, reg wsllib.WslReg, name string, rootPath string) error {
	// create empty tar
	tmptar := os.TempDir()
	if tmptar == "" {
		return errors.New("failed to create temp directory")
	}
	tmptar = tmptar + "\\em-vhdx-temp.tar"
	tmptarfp, err := os.Create(tmptar)
	if err != nil {
		return err
	}
	tmptarfp.Close()
	// initial empty instance entry
	err = wsl.RegisterDistribution(name, tmptar)
	if err != nil {
		return err
	}
	os.Remove(tmptar)
	// get profile of instance
	prof, err := reg.GetProfileFromName(name)
	if prof.BasePath == "" {
		return err
	}
	// remove instance temporary
	err = wsl.UnregisterDistribution(name)
	if err != nil {
		return err
	}
	// copy vhdx to destination directory
	err = fileutil.CopyFileAndCompress(rootPath, prof.BasePath+"\\ext4.vhdx")
	if err != nil {
		return err
	}

	// write registry
	prof.Flags |= wsllib.FlagEnableWsl2
	err = reg.WriteProfile(prof)
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
