package install

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuk7/wsldl/lib/download"
	"github.com/yuk7/wsldl/lib/errutil"
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

type installDeps struct {
	tempDir       func() string
	createFile    func(path string) (io.Closer, error)
	removeFile    func(path string) error
	copyFile      func(srcPath, destPath string, compress bool) error
	confirmResume func() bool
}

func defaultInstallDeps() installDeps {
	return installDeps{
		tempDir: os.TempDir,
		createFile: func(path string) (io.Closer, error) {
			return os.Create(path)
		},
		removeFile: os.Remove,
		copyFile:   fileutil.CopyFile,
		confirmResume: func() bool {
			fmt.Printf("A partial download was found.\n")
			fmt.Printf("Do you want to resume the download?\n")
			fmt.Printf("Type y/n:")
			in, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			return strings.EqualFold(strings.TrimSpace(in), "y")
		},
	}
}

func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// Install installs distribution with default rootfs file names
func Install(ctx context.Context, wsl wsllib.WslLib, reg wsllib.WslReg, name string, rootPath string, sha256Sum string, showProgress bool) error {
	return installWithDeps(ctx, wsl, reg, name, rootPath, sha256Sum, showProgress, defaultInstallDeps())
}

func installWithDeps(ctx context.Context, wsl wsllib.WslLib, reg wsllib.WslReg, name string, rootPath string, sha256Sum string, showProgress bool, deps installDeps) error {
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return err
	}

	rootPathLower := strings.ToLower(rootPath)
	sha256Actual := ""
	usedCachedDownload := false
	if showProgress {
		fmt.Printf("Using: %s\n", rootPath)
	}

	if strings.HasPrefix(rootPathLower, "http://") || strings.HasPrefix(rootPathLower, "https://") {
		progressBarWidth := 0
		if showProgress {
			progressBarWidth = 35
		}
		tmpRootDir := deps.tempDir()
		if tmpRootDir == "" {
			return errors.New("failed to create temp directory")
		}
		downloadURL := rootPath
		cacheRootPath := getDownloadCachePath(tmpRootDir, downloadURL)
		cachePartialPath := cacheRootPath + ".part"
		downloadToCache := func() (string, error) {
			if showProgress {
				fmt.Println("Downloading...")
			}
			sum, err := download.DownloadFile(ctx, downloadURL, cachePartialPath, progressBarWidth)
			if err != nil {
				return "", err
			}
			if err := os.Rename(cachePartialPath, cacheRootPath); err != nil {
				return "", err
			}
			if showProgress {
				fmt.Println()
			}
			return sum, nil
		}

		if _, err := os.Stat(cacheRootPath); err == nil {
			usedCachedDownload = true
			rootPath = cacheRootPath
			if showProgress {
				fmt.Printf("Using cached download: %s\n", rootPath)
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		} else {
			if _, err := os.Stat(cachePartialPath); err == nil {
				keepPartial := true
				if showProgress {
					if deps.confirmResume != nil {
						keepPartial = deps.confirmResume()
					}
				}
				if !keepPartial {
					_ = deps.removeFile(cachePartialPath)
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				return err
			}

			var err error
			sha256Actual, err = downloadToCache()
			if err != nil {
				return err
			}
			rootPath = cacheRootPath
		}
		rootPathLower = strings.ToLower(rootPath)

		if sha256Sum != "" && sha256Actual == "" {
			if showProgress {
				fmt.Println("Calculating checksum...")
			}
			var err error
			sha256Actual, err = calculateFileSHA256(rootPath)
			if err != nil {
				return err
			}
		}
		shouldRetryOnChecksumMismatch := false
		if sha256Sum != "" && sha256Actual != "" && sha256Sum != sha256Actual {
			shouldRetryOnChecksumMismatch = usedCachedDownload || !showProgress
		}
		if shouldRetryOnChecksumMismatch {
			if showProgress {
				fmt.Println("Checksum mismatch. Re-downloading...")
			}
			_ = deps.removeFile(cacheRootPath)
			_ = deps.removeFile(cachePartialPath)
			var err error
			sha256Actual, err = downloadToCache()
			if err != nil {
				return err
			}
			rootPath = cacheRootPath
			rootPathLower = strings.ToLower(rootPath)
		}
	} else if sha256Sum != "" {
		if showProgress {
			fmt.Println("Calculating checksum...")
		}
		var err error
		sha256Actual, err = calculateFileSHA256(rootPath)
		if err != nil {
			return err
		}
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

	if err := ctx.Err(); err != nil {
		return err
	}

	if strings.HasSuffix(rootPathLower, "ext4.vhdx") || strings.HasSuffix(rootPathLower, "ext4.vhdx.gz") {
		return installExt4VhdxWithDeps(wsl, reg, name, rootPath, deps)
	}
	return InstallTar(wsl, name, rootPath)
}

func getDownloadCachePath(tempDir, rawURL string) string {
	u, err := url.Parse(rawURL)
	cacheBase := "download.bin"
	if err == nil {
		if base := filepath.Base(u.Path); base != "." && base != "/" && base != "" {
			cacheBase = base
		}
	}
	urlHash := sha256.Sum256([]byte(rawURL))
	return filepath.Join(tempDir, fmt.Sprintf("wsldl-download-%s-%s", hex.EncodeToString(urlHash[:8]), cacheBase))
}

func calculateFileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func InstallTar(wsl wsllib.WslLib, name string, rootPath string) error {
	err := wsl.RegisterDistribution(name, rootPath)
	return err
}

func InstallExt4Vhdx(wsl wsllib.WslLib, reg wsllib.WslReg, name string, rootPath string) error {
	return installExt4VhdxWithDeps(wsl, reg, name, rootPath, defaultInstallDeps())
}

func installExt4VhdxWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg, name string, rootPath string, deps installDeps) error {
	// create empty tar
	tmptar := deps.tempDir()
	if tmptar == "" {
		return errors.New("failed to create temp directory")
	}
	tmptar = filepath.Join(tmptar, "em-vhdx-temp.tar")
	tmptarfp, err := deps.createFile(tmptar)
	if err != nil {
		return err
	}
	tmptarfp.Close()
	// initial empty instance entry
	err = wsl.RegisterDistribution(name, tmptar)
	if err != nil {
		return err
	}
	deps.removeFile(tmptar)
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
	err = deps.copyFile(rootPath, filepath.Join(prof.BasePath, "ext4.vhdx"), false)
	if err != nil {
		return err
	}

	// write registry
	prof.Flags |= wsllib.FlagEnableWsl2
	err = reg.WriteProfile(prof)
	return err
}

func detectRootfsFiles() (string, error) {
	efPath := errutil.MustExecutable()
	efDir := filepath.Dir(efPath)
	rootFile, err := detectRootfsFileName(os.DirFS(efDir))
	if err != nil {
		return "", err
	}
	if rootFile == "rootfs.tar.gz" {
		return rootFile, nil
	}
	return filepath.Join(efDir, rootFile), nil
}

func detectRootfsFileName(root fs.FS) (string, error) {
	for _, rootFile := range defaultRootFiles {
		if _, err := fs.Stat(root, rootFile); err == nil {
			return rootFile, nil
		}
	}
	return "", errors.New("no rootfs file found in the directory")
}
