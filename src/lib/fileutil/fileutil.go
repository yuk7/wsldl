package fileutil

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// SpecialDirs is define path of special dirs
	SpecialDirs = "SystemDrive:,SystemRoot:,SystemRoot:System32,USERPROFILE:"
)

// DQEscapeString is escape string with double quote
func DQEscapeString(str string) string {
	if strings.Contains(str, " ") {
		str = strings.Replace(str, "\"", "\\\"", -1)
		str = "\"" + str + "\""
	}
	return str
}

// GetWindowsDirectory gets windows direcotry path
func GetWindowsDirectory() string {
	return getWindowsDirectoryFromEnv(os.Getenv)
}

func getWindowsDirectoryFromEnv(getenv func(string) string) string {
	dir := getenv("SYSTEMROOT")
	if dir != "" {
		return dir
	}
	dir = getenv("WINDIR")
	if dir != "" {
		return dir
	}
	return "C:\\WINDOWS"
}

// IsCurrentDirSpecial gets whether the current directory is special (Windows, USEPROFILE)
func IsCurrentDirSpecial() bool {
	return isCurrentDirSpecial(filepath.Abs, os.Getenv)
}

func isCurrentDirSpecial(absPath func(string) (string, error), getenv func(string) string) bool {
	cdir, err := absPath(".")
	if err != nil {
		return true
	}
	return isSpecialDir(cdir, getenv, absPath)
}

func isSpecialDir(cdir string, getenv func(string) string, absPath func(string) (string, error)) bool {
	sdarr := strings.Split(SpecialDirs, ",")
	for _, item := range sdarr {
		splititem := strings.SplitN(item, ":", 2)
		itemdir := ""
		if splititem[0] != "" {
			itemdir = getenv(splititem[0])
		}
		suffix := ""
		if len(splititem) > 1 {
			suffix = splititem[1]
		}
		itemdir, err := absPath(itemdir + "\\" + suffix)
		if err != nil {
			return true
		}
		if strings.EqualFold(cdir, itemdir) {
			return true
		}
	}
	return false
}

// CopyFile copies a file to destination.
// If src has .gz/.tgz suffix, it is transparently decompressed.
// If compress is true, destination data is gzip-compressed.
func CopyFile(srcPath, destPath string, compress bool) error {
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

	srcReader := io.Reader(src)
	srcPathLower := strings.ToLower(srcPath)
	if strings.HasSuffix(srcPathLower, ".gz") || strings.HasSuffix(srcPathLower, ".tgz") {
		gr, err := gzip.NewReader(src)
		if err != nil {
			return err
		}
		defer gr.Close()
		srcReader = gr
	}

	destWriter := io.Writer(dest)
	if compress {
		gw := gzip.NewWriter(dest)
		defer gw.Close()
		destWriter = gw
	}

	_, err = io.Copy(destWriter, srcReader)
	return err
}

// CopyFileAndCompress copies a file to the destination and gzip-compresses when destination has .gz/.tgz suffix.
func CopyFileAndCompress(srcPath, destPath string) error {
	destPathLower := strings.ToLower(destPath)
	compress := strings.HasSuffix(destPathLower, ".gz") || strings.HasSuffix(destPathLower, ".tgz")
	return CopyFile(srcPath, destPath, compress)
}
