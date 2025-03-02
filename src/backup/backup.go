package backup

import (
	"compress/gzip"
	"errors"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/yuk7/wsldl/lib/utils"
	wslreg "github.com/yuk7/wslreglib-go"
)

func backupReg(name string, destFileName string) error {
	profile, err := wslreg.GetProfileFromName(name)
	if err != nil {
		utils.ErrorExit(err, true, true, false)
	}

	regexe := utils.GetWindowsDirectory() + "\\System32\\reg.exe"
	regpath := "HKEY_CURRENT_USER\\" + wslreg.LxssBaseKey + "\\" + profile.UUID
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
		wslexe := utils.GetWindowsDirectory() + "\\System32\\wsl.exe"
		_, err := exec.Command(wslexe, "--export", distributionName, tmpTarFn).Output()
		defer os.Remove(tmpTarFn)
		if err != nil {
			return err
		}

		return copyFileAndCompress(tmpTarFn, destFileName)
	} else {
		// not compressed
		wslexe := utils.GetWindowsDirectory() + "\\System32\\wsl.exe"
		_, err := exec.Command(wslexe, "--export", distributionName, destFileName).Output()
		return err
	}
}

func backupExt4Vhdx(name string, destFileName string) error {
	prof, err := wslreg.GetProfileFromName(name)
	if prof.BasePath != "" {

	} else {
		if err != nil {
			return err
		}
		return errors.New("Get profile failed")
	}

	vhdxPath := prof.BasePath + "\\ext4.vhdx"

	return copyFileAndCompress(vhdxPath, destFileName)
}

func copyFileAndCompress(srcPath, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	// compress and copy
	destPathLower := strings.ToLower(destPath)
	if strings.HasSuffix(destPathLower, ".gz") || strings.HasSuffix(destPathLower, ".tgz") {
		// compressed with gzip
		gw := gzip.NewWriter(dest)
		defer gw.Close()
		_, err = io.Copy(gw, src)
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
	return nil
}
