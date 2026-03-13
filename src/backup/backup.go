package backup

import (
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yuk7/wsldl/lib/fileutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

type backupTarDeps struct {
	tempDir  func() string
	export   func(distributionName, destFileName string) error
	copyFile func(srcPath, destPath string, compress bool) error
	remove   func(path string) error
	randIntn func(int) int
}

func defaultBackupTarDeps() backupTarDeps {
	return backupTarDeps{
		tempDir: os.TempDir,
		export: func(distributionName, destFileName string) error {
			wslexe := filepath.Join(fileutil.GetWindowsDirectory(), "System32", "wsl.exe")
			_, err := exec.Command(wslexe, "--export", distributionName, destFileName).Output()
			return err
		},
		copyFile: fileutil.CopyFile,
		remove:   os.Remove,
		randIntn: rand.Intn,
	}
}

func backupReg(reg wsllib.WslReg, name string, destFileName string) error {
	profile, err := reg.GetProfileFromName(name)
	if err != nil {
		return err
	}

	regexe := filepath.Join(fileutil.GetWindowsDirectory(), "System32", "reg.exe")
	regpath := "HKEY_CURRENT_USER\\" + wsllib.LxssBaseKey + "\\" + profile.UUID
	_, err = exec.Command(regexe, "export", regpath, destFileName, "/y").Output()
	return err
}

func backupTar(distributionName string, destFileName string) error {
	return backupTarWithDeps(distributionName, destFileName, defaultBackupTarDeps())
}

func backupTarWithDeps(distributionName string, destFileName string, deps backupTarDeps) error {
	// compress and copy
	rootPathLower := strings.ToLower(destFileName)
	if strings.HasSuffix(rootPathLower, ".gz") {
		// create temporary tar
		tmpDir := deps.tempDir()
		if tmpDir == "" {
			return errors.New("failed to create temp directory")
		}
		tmpTarFn := filepath.Join(tmpDir, strconv.Itoa(deps.randIntn(10000))+".tar")
		err := deps.export(distributionName, tmpTarFn)
		defer deps.remove(tmpTarFn)
		if err != nil {
			return err
		}

		return deps.copyFile(tmpTarFn, destFileName, true)
	} else {
		// not compressed
		return deps.export(distributionName, destFileName)
	}
}

func backupExt4Vhdx(reg wsllib.WslReg, name string, destFileName string) error {
	return backupExt4VhdxWithCopy(reg, name, destFileName, fileutil.CopyFile)
}

func backupExt4VhdxWithCopy(reg wsllib.WslReg, name string, destFileName string, copyFile func(srcPath, destPath string, compress bool) error) error {
	prof, err := reg.GetProfileFromName(name)
	if prof.BasePath == "" {
		if err != nil {
			return err
		}
		return errors.New("get profile failed")
	}
	if err != nil {
		return err
	}

	vhdxPath := filepath.Join(prof.BasePath, "ext4.vhdx")

	rootPathLower := strings.ToLower(destFileName)
	compress := strings.HasSuffix(rootPathLower, ".gz") || strings.HasSuffix(rootPathLower, ".tgz")
	return copyFile(vhdxPath, destFileName, compress)
}
