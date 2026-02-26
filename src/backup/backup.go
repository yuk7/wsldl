package backup

import (
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yuk7/wsldl/lib/fileutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func backupReg(reg wsllib.WslReg, name string, destFileName string) error {
	profile, err := reg.GetProfileFromName(name)
	if err != nil {
		return err
	}

	regexe := fileutil.GetWindowsDirectory() + "\\System32\\reg.exe"
	regpath := "HKEY_CURRENT_USER\\" + wsllib.LxssBaseKey + "\\" + profile.UUID
	_, err = exec.Command(regexe, "export", regpath, destFileName, "/y").Output()
	return err
}

func backupTar(distributionName string, destFileName string) error {
	// compress and copy
	rootPathLower := strings.ToLower(destFileName)
	if strings.HasSuffix(rootPathLower, ".gz") {
		// create temporary tar
		tmpTarFn := os.TempDir()
		if tmpTarFn == "" {
			return errors.New("failed to create temp directory")
		}
		rand.NewSource(time.Now().UnixNano())
		tmpTarFn = tmpTarFn + "\\" + strconv.Itoa(rand.Intn(10000)) + ".tar"
		wslexe := fileutil.GetWindowsDirectory() + "\\System32\\wsl.exe"
		_, err := exec.Command(wslexe, "--export", distributionName, tmpTarFn).Output()
		defer os.Remove(tmpTarFn)
		if err != nil {
			return err
		}

		return fileutil.CopyFile(tmpTarFn, destFileName, true)
	} else {
		// not compressed
		wslexe := fileutil.GetWindowsDirectory() + "\\System32\\wsl.exe"
		_, err := exec.Command(wslexe, "--export", distributionName, destFileName).Output()
		return err
	}
}

func backupExt4Vhdx(reg wsllib.WslReg, name string, destFileName string) error {
	return backupExt4VhdxWithCopy(reg, name, destFileName, fileutil.CopyFile)
}

func backupExt4VhdxWithCopy(reg wsllib.WslReg, name string, destFileName string, copyFile func(srcPath, destPath string, compress bool) error) error {
	prof, err := reg.GetProfileFromName(name)
	if prof.BasePath != "" {

	} else {
		if err != nil {
			return err
		}
		return errors.New("get profile failed")
	}

	vhdxPath := filepath.Join(prof.BasePath, "ext4.vhdx")

	rootPathLower := strings.ToLower(destFileName)
	compress := strings.HasSuffix(rootPathLower, ".gz") || strings.HasSuffix(rootPathLower, ".tgz")
	return copyFile(vhdxPath, destFileName, compress)
}
