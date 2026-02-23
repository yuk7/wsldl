package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/console"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/fileutil"
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
func execute(name string, args []string) error {
	command := ""
	for _, s := range args {
		command = command + " " + fileutil.DQEscapeString(s)
	}
	var inheritpath = true
	if args == nil {
		inheritpath = !fileutil.IsCurrentDirSpecial()
	}
	exitCode, err := wsllib.WslLaunchInteractive(name, command, inheritpath)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	if exitCode != 0 {
		return errutil.NewExitCodeError(int(exitCode), false)
	}
	return nil
}

// executeP runs execute function with Path Translator
func executeP(name string, args []string) error {
	var convArgs []string
	for _, s := range args {
		if strings.Contains(s, "\\") {
			s = strings.Replace(s, "\\", "/", -1)
			s = fileutil.DQEscapeString(s)
			out, exitCode, err := ExecRead(name, "wslpath -u "+s)
			if err != nil || exitCode != 0 {
				errutil.ErrorRedPrintln("ERR: Failed to Path Translation")
				fmt.Fprintf(os.Stderr, "ExitCode: 0x%x\n", int(exitCode))
				if err != nil {
					return errutil.NewDisplayError(err, true, true, false)
				}
				return errutil.NewExitCodeError(int(exitCode), false)
			}
			convArgs = append(convArgs, out)
		} else {
			convArgs = append(convArgs, s)
		}
	}

	return execute(name, convArgs)
}

// executeNoArgs runs distro, but use terminal settings
func executeNoArgs(name string, args []string) error {
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
						return errutil.NewDisplayError(err, true, true, true)
					}
					errutil.StdoutGreenPrintln("done.")
					return errutil.NewExitCodeError(0, true)
				}
			}
		}
	}

	b, err := console.IsParentConsole()
	if err != nil {
		b = true
	}
	if !b {
		switch profile.WsldlTerm {
		case wslreg.FlagWsldlTermWT:
			console.FreeConsole()
			return ExecWindowsTerminal(name)

		case wslreg.FlagWsldlTermFlute:
			console.FreeConsole()
			exe := os.Getenv("LOCALAPPDATA")
			exe = fileutil.DQEscapeString(exe + "\\Microsoft\\WindowsApps\\53621FSApps.FluentTerminal_87x1pks76srcp\\flute.exe")

			cmd := exe + " run " + fileutil.DQEscapeString(efPath+" run")
			res, err := utils.CreateProcessAndWait(cmd)
			if err != nil {
				console.AllocConsole()
				fmt.Fprintln(os.Stderr, "ERR: Failed to launch the terminal process")
				fmt.Fprintf(os.Stderr, "%s\n", exe)
				return errutil.NewDisplayError(err, true, false, true)
			}
			if res != 0 {
				return errutil.NewExitCodeError(res, false)
			}
			return nil
		}

		// Parent isn't console, launch instance with default conhost
		// Get the name from the registry to be case sensitive.
		if profile.DistributionName != "" {
			name = profile.DistributionName
		}

		console.SetConsoleTitle(name)
		return execute(name, nil)
	} else {
		return execute(name, nil)
	}
}
