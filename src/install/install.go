package install

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuk7/wsldl/lib/wslapi"
)

var (
	defaultRootFiles = []string{"install.tar", "install.tar.gz", "rootfs.tar", "rootfs.tar.gz"}
)

//Install installs distribution with default rootfs file names
func Install(name string, rootPath string, showProgress bool) error {
	if showProgress {
		fmt.Printf("Using: %s\n", rootPath)
		fmt.Println("Installing...")
	}
	err := wslapi.WslRegisterDistribution(name, rootPath)
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
