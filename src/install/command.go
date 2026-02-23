package install

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/fileutil"
	"github.com/yuk7/wsldl/lib/preset"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func GetCommandWithNoArgs() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithNoArgsWithDeps(deps.Wsl, deps.Reg)
}

func GetCommandWithNoArgsWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg) cmdline.Command {
	return cmdline.Command{
		Names: []string{},
		Help: func(distroName string, isListQuery bool) string {
			if !wsl.IsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessageNoArgs()
			}
			return ""
		},
		Run: func(name string, args []string) error {
			return execute(wsl, reg, name, args)
		},
	}
}

// GetCommand returns the install command structure
func GetCommand() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithDeps(deps.Wsl, deps.Reg)
}

// GetCommandWithDeps returns the install command structure with injectable dependencies.
func GetCommandWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg) cmdline.Command {
	return cmdline.Command{
		Names: []string{"install"},
		Help: func(distroName string, isListQuery bool) string {
			if !wsl.IsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			}
			return ""
		},
		Run: func(name string, args []string) error {
			return execute(wsl, reg, name, args)
		},
	}
}

// execute is default install entrypoint
func execute(wsl wsllib.WslLib, reg wsllib.WslReg, name string, args []string) error {
	if !wsl.IsDistributionRegistered(name) {
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
			return errutil.NewDisplayError(os.ErrInvalid, true, true, false)
		}

		if args == nil {
			if isInstalledFilesExist() {
				var in string
				fmt.Printf("An old installation files has been found.\n")
				fmt.Printf("Do you want to rewrite and repair the installation information?\n")
				fmt.Printf("Type y/n:")
				fmt.Scan(&in)

				if in == "y" {
					err := repairRegistry(reg, name)
					if err != nil {
						return errutil.NewDisplayError(err, showProgress, true, showProgress)
					}
					errutil.StdoutGreenPrintln("done.")
					return nil
				}
			}
		}

		err := Install(wsl, reg, name, rootPath, rootFileSha256, showProgress)
		if err != nil {
			return errutil.NewDisplayError(err, showProgress, true, args == nil)
		}

		if jsonPreset.WslVersion == 1 || jsonPreset.WslVersion == 2 {
			wslexe := fileutil.GetWindowsDirectory() + "\\System32\\wsl.exe"
			_, err = exec.Command(wslexe, "--set-version", name, strconv.Itoa(jsonPreset.WslVersion)).Output()
		}

		if err == nil {
			if showProgress {
				errutil.StdoutGreenPrintln("Installation complete")
			}
		} else {
			return errutil.NewDisplayError(err, showProgress, true, args == nil)
		}

		if args == nil {
			fmt.Fprintf(os.Stdout, "Press enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')
		}

	} else {
		errutil.ErrorRedPrintln("ERR: [" + name + "] is already installed.")
		return errutil.NewDisplayError(os.ErrInvalid, false, true, false)
	}
	return nil
}
