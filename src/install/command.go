package install

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/preset"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
)

func GetCommandWithNoArgs() cmdline.Command {
	return cmdline.Command{
		Names: []string{},
		Help: func(distroName string, isListQuery bool) string {
			if !wsllib.WslIsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessageNoArgs()
			}
			return ""
		},
		Run: execute,
	}
}

// GetCommand returns the install command structure
func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"install"},
		Help: func(distroName string, isListQuery bool) string {
			if !wsllib.WslIsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			}
			return ""
		},
		Run: execute,
	}
}

// execute is default install entrypoint
func execute(name string, args []string) {
	if !wsllib.WslIsDistributionRegistered(name) {
		var rootPath string
		var rootFileSha256 string = ""
		var showProgress bool
		jsonPreset, _ := preset.ReadParsePreset()
		switch len(args) {
		case 0:
			rootPath = detectRootfsFiles()
			if jsonPreset.InstallFile != "" {
				rootPath = jsonPreset.InstallFile
			}
			rootFileSha256 = jsonPreset.InstallFileSha256
			showProgress = true

		case 1:
			showProgress = false
			if args[0] == "--root" {
				rootPath = detectRootfsFiles()
				if jsonPreset.InstallFile != "" {
					rootPath = jsonPreset.InstallFile
				}
				rootFileSha256 = jsonPreset.InstallFileSha256
			} else {
				rootPath = args[0]
			}

		default:
			utils.ErrorExit(os.ErrInvalid, true, true, false)
		}

		if args == nil {
			if isInstalledFilesExist() {
				var in string
				fmt.Printf("An old installation files has been found.\n")
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

		err := Install(name, rootPath, rootFileSha256, showProgress)
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
