package install

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/yuk7/wsldl/lib/preset"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
)

//Execute is default install entrypoint
func Execute(name string, args []string) {
	if !wsllib.WslIsDistributionRegistered(name) {
		var rootPath string
		var showProgress bool
		jsonPreset, _ := preset.ReadParsePreset()
		switch len(args) {
		case 0:
			rootPath = detectRootfsFiles()
			if jsonPreset.InstallFile != "" {
				rootPath = jsonPreset.InstallFile
			}
			showProgress = true

		case 1:
			showProgress = false
			if args[0] == "--root" {
				rootPath = detectRootfsFiles()
				if jsonPreset.InstallFile != "" {
					rootPath = jsonPreset.InstallFile
				}
			} else {
				rootPath = args[0]
			}

		default:
			utils.ErrorExit(os.ErrInvalid, true, true, false)
		}

		if args == nil {
			if isInstalledFilesExist() {
				var in string
				fmt.Printf("An old installation file was found.\n")
				fmt.Printf("Do you want to rewrite and repair the installation information?\n")
				fmt.Printf("Type y/n:")
				fmt.Scan(&in)

				if in == "y" {
					err := repairRegistry(name)
					if err != nil {
						utils.ErrorExit(err, showProgress, true, showProgress)
					}
					utils.StdoutGreenPrintln("done.")
					return
				}
			}
		}

		err := Install(name, rootPath, showProgress)
		if err != nil {
			utils.ErrorExit(err, showProgress, true, args == nil)
		}

		if jsonPreset.WslVersion == 1 || jsonPreset.WslVersion == 2 {
			wslexe := utils.GetWindowsDirectory() + "\\System32\\wsl.exe"
			_, err = exec.Command(wslexe, "--set-version", name, strconv.Itoa(jsonPreset.WslVersion)).Output()
		}

		if err == nil {
			if showProgress {
				utils.StdoutGreenPrintln("Installation complete")
			}
		} else {
			utils.ErrorExit(err, showProgress, true, args == nil)
		}

		if args == nil {
			fmt.Fprintf(os.Stdout, "Press enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')
		}

	} else {
		utils.ErrorRedPrintln("ERR: [" + name + "] is already installed.")
		utils.ErrorExit(os.ErrInvalid, false, true, false)
	}
}
