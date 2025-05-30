package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
	wslreg "github.com/yuk7/wslreglib-go"
)

// GetCommandWithNoArgs returns the run command structure with no arguments
func GetCommandWithNoArgs() cmdline.Command {
	return cmdline.Command{
		Names: []string{},
		Help: func(distroName string, isListQuery bool) string {
			if wsllib.WslIsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessageNoArgs()
			} else {
				return ""
			}
		},
		Run: executeNoArgs,
	}
}

// GetCommand returns the run command structure
func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"run", "-c", "/c"},
		Help: func(distroName string, isListQuery bool) string {
			if wsllib.WslIsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			} else {
				return ""
			}
		},
		Run: execute,
	}
}

// GetCommandP returns the runp command structure
func GetCommandP() cmdline.Command {
	return cmdline.Command{
		Names: []string{"runp", "-p", "/p"},
		Help: func(distroName string, isListQuery bool) string {
			if wsllib.WslIsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessageP()
			} else {
				return ""
			}
		},
		Run: executeP,
	}
}

// execute is default run entrypoint.
func execute(name string, args []string) {
	command := ""
	for _, s := range args {
		command = command + " " + utils.DQEscapeString(s)
	}
	var inheritpath = true
	if args == nil {
		inheritpath = !utils.IsCurrentDirSpecial()
	}
	exitCode, err := wsllib.WslLaunchInteractive(name, command, inheritpath)
	if err != nil {
		utils.ErrorExit(err, true, true, false)
	} else {
		os.Exit(int(exitCode))
	}
}

// executeP runs execute function with Path Translator
func executeP(name string, args []string) {
	var convArgs []string
	for _, s := range args {
		if strings.Contains(s, "\\") {
			s = strings.Replace(s, "\\", "/", -1)
			s = utils.DQEscapeString(s)
			out, exitCode, err := ExecRead(name, "wslpath -u "+s)
			if err != nil || exitCode != 0 {
				utils.ErrorRedPrintln("ERR: Failed to Path Translation")
				fmt.Fprintf(os.Stderr, "ExitCode: 0x%x\n", int(exitCode))
				if err != nil {
					utils.ErrorExit(err, true, true, false)
				}
				os.Exit(int(exitCode))
			}
			convArgs = append(convArgs, out)
		} else {
			convArgs = append(convArgs, s)
		}
	}

	execute(name, convArgs)
}

// executeNoArgs runs distro, but use terminal settings
func executeNoArgs(name string, args []string) {
	efPath, _ := os.Executable()
	profile, _ := wslreg.GetProfileFromName(name)

	// repair when the installation is moved
	if profile.BasePath != "" {
		_, err := os.Stat(profile.BasePath)
		if os.IsNotExist(err) {
			if isInstalledFilesExist() {
				var in string
				fmt.Printf("This instance (%s) BasePath is not exist.\n", name)
				fmt.Printf("Do you want to repair the installation information?\n")
				fmt.Printf("Type y/n:")
				fmt.Scan(&in)

				if in == "y" {
					err := repairRegistry(profile)
					if err != nil {
						utils.ErrorExit(err, true, true, true)
					}
					utils.StdoutGreenPrintln("done.")
					utils.Exit(true, 0)
				}
			}
		}
	}

	b, err := utils.IsParentConsole()
	if err != nil {
		b = true
	}
	if !b {
		switch profile.WsldlTerm {
		case wslreg.FlagWsldlTermWT:
			utils.FreeConsole()
			ExecWindowsTerminal(name)
			os.Exit(0)

		case wslreg.FlagWsldlTermFlute:
			utils.FreeConsole()
			exe := os.Getenv("LOCALAPPDATA")
			exe = utils.DQEscapeString(exe + "\\Microsoft\\WindowsApps\\53621FSApps.FluentTerminal_87x1pks76srcp\\flute.exe")

			cmd := exe + " run " + utils.DQEscapeString(efPath+" run")
			res, err := utils.CreateProcessAndWait(cmd)
			if err != nil {
				utils.AllocConsole()
				fmt.Fprintln(os.Stderr, "ERR: Failed to launch the terminal process")
				fmt.Fprintf(os.Stderr, "%s\n", exe)
				utils.ErrorExit(err, true, false, true)
			}
			os.Exit(res)
		}

		// Parent isn't console, launch instance with default conhost
		// Get the name from the registry to be case sensitive.
		if profile.DistributionName != "" {
			name = profile.DistributionName
		}

		utils.SetConsoleTitle(name)
		execute(name, nil)
	} else {
		execute(name, nil)
	}
}
