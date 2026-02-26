package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/console"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/fileutil"
	"github.com/yuk7/wsldl/lib/repair"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wslexec"
	"github.com/yuk7/wsldl/lib/wsllib"
)

// GetCommandWithNoArgs returns the run command structure with no arguments
func GetCommandWithNoArgs() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithNoArgsWithDeps(deps.Wsl, deps.Reg)
}

// GetCommandWithNoArgsWithDeps returns the run command structure with no arguments and injectable dependencies.
func GetCommandWithNoArgsWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg) cmdline.Command {
	return cmdline.Command{
		Names: []string{},
		Help: func(distroName string, isListQuery bool) string {
			if wsl.IsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessageNoArgs()
			} else {
				return ""
			}
		},
		Run: func(name string, args []string) error {
			return executeNoArgs(wsl, reg, name, args)
		},
	}
}

// GetCommand returns the run command structure
func GetCommand() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithDeps(deps.Wsl)
}

// GetCommandWithDeps returns the run command structure with injectable dependencies.
func GetCommandWithDeps(wsl wsllib.WslLib) cmdline.Command {
	return cmdline.Command{
		Names: []string{"run", "-c", "/c"},
		Help: func(distroName string, isListQuery bool) string {
			if wsl.IsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			} else {
				return ""
			}
		},
		Run: func(name string, args []string) error {
			return execute(wsl, name, args)
		},
	}
}

// GetCommandP returns the runp command structure
func GetCommandP() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandPWithDeps(deps.Wsl)
}

// GetCommandPWithDeps returns the runp command structure with injectable dependencies.
func GetCommandPWithDeps(wsl wsllib.WslLib) cmdline.Command {
	return cmdline.Command{
		Names: []string{"runp", "-p", "/p"},
		Help: func(distroName string, isListQuery bool) string {
			if wsl.IsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessageP()
			} else {
				return ""
			}
		},
		Run: func(name string, args []string) error {
			return executeP(wsl, name, args)
		},
	}
}

// execute is default run entrypoint.
func execute(wsl wsllib.WslLib, name string, args []string) error {
	command := ""
	for _, s := range args {
		command = command + " " + fileutil.DQEscapeString(s)
	}
	var inheritpath = true
	if args == nil {
		inheritpath = !fileutil.IsCurrentDirSpecial()
	}
	exitCode, err := wsl.LaunchInteractive(name, command, inheritpath)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	if exitCode != 0 {
		return errutil.NewExitCodeError(int(exitCode), false)
	}
	return nil
}

// executeP runs execute function with Path Translator
func executeP(wsl wsllib.WslLib, name string, args []string) error {
	var convArgs []string
	for _, s := range args {
		if strings.Contains(s, "\\") {
			s = strings.Replace(s, "\\", "/", -1)
			s = fileutil.DQEscapeString(s)
			out, exitCode, err := wslexec.ExecRead(wsl, name, "wslpath -u "+s)
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

	return execute(wsl, name, convArgs)
}

// executeNoArgs runs distro, but use terminal settings
func executeNoArgs(wsl wsllib.WslLib, reg wsllib.WslReg, name string, args []string) error {
	efPath, _ := os.Executable()
	profile, _ := reg.GetProfileFromName(name)

	// repair when the installation is moved
	if profile.BasePath != "" {
		_, err := os.Stat(profile.BasePath)
		if os.IsNotExist(err) {
			if repair.IsInstalledFilesExist() {
				var in string
				fmt.Printf("This instance (%s) BasePath is not exist.\n", name)
				fmt.Printf("Do you want to repair the installation information?\n")
				fmt.Printf("Type y/n:")
				fmt.Scan(&in)

				if in == "y" {
					err := repairRegistry(reg, profile)
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
		case wsllib.FlagWsldlTermWT:
			console.FreeConsole()
			return ExecWindowsTerminal(reg, name)

		case wsllib.FlagWsldlTermFlute:
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
		return execute(wsl, name, nil)
	} else {
		return execute(wsl, name, nil)
	}
}
