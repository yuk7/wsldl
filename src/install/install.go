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

//Execute is default install entrypoint
func Execute(name string, arg []string) {
	Install("", true)
}

//Install installs distribution with default rootfs file names
func Install(name string, showProgress bool) {
	rootPath := detectRootfsFiles()
	if showProgress {
		fmt.Printf("Using: %s\n", rootPath)
		fmt.Println("Installing...")
	}
	res := wslapi.WslRegisterDistribution(name, rootPath)
	if showProgress {
		if res != 0 {
			fmt.Println("ERR: Failed to Install")
			fmt.Printf("Error code 0x%x", res)
		} else {
			fmt.Println("Installation complete")
		}
	}
}

func detectRootfsFiles() string {
	efPath, _ := os.Executable()
	efDir := filepath.Dir(efPath)
	fmt.Println(efDir)
	for _, rootFile := range defaultRootFiles {
		rootPath := filepath.Join(efDir, rootFile)
		_, err := os.Stat(rootPath)
		if err == nil {
			return rootPath
		}
	}
	return "rootfs.tar.gz"
}
