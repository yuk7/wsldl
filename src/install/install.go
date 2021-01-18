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
func Execute(name string, args []string) {
	if !wslapi.WslIsDistributionRegistered(name) {
		var rootPath string
		var showProgress bool
		switch len(args) {
		case 0:
			rootPath = detectRootfsFiles()
			showProgress = true

		case 1:
			showProgress = false
			if args[0] == "--root" {
				rootPath = detectRootfsFiles()
			} else {
				rootPath = args[0]
			}

		default:
			fmt.Println("Invalid Arg.")
			os.Exit(1)
		}

		Install(name, rootPath, showProgress)

	} else {
		fmt.Printf("ERR: [%s] is already installed.\n", name)
	}
}

//Install installs distribution with default rootfs file names
func Install(name string, rootPath string, showProgress bool) {
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
	for _, rootFile := range defaultRootFiles {
		rootPath := filepath.Join(efDir, rootFile)
		_, err := os.Stat(rootPath)
		if err == nil {
			return rootPath
		}
	}
	return "rootfs.tar.gz"
}
