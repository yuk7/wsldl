package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wslapi"
	"github.com/yuk7/wsldl/lib/wslreg"
)

//Execute is default run entrypoint.
func Execute(name string, args []string) {
	command := ""
	for _, s := range args {
		command = command + " " + utils.DQEscapeString(s)
	}
	var inheritpath = true
	if args == nil {
		inheritpath = !utils.IsCurrentDirSpecial()
	}
	exitCode, err := wslapi.WslLaunchInteractive(name, command, inheritpath)
	if err != nil {
		utils.ErrorExit(err, true, true, false)
	} else {
		os.Exit(int(exitCode))
	}
}

//ExecuteP runs Execute function with Path Translator
func ExecuteP(name string, args []string) {
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

	Execute(name, convArgs)
}

// ExecuteNoArgs runs distro, but use terminal settings
func ExecuteNoArgs(name string) {
	efPath, _ := os.Executable()
	b, err := utils.IsParentConsole()
	if err != nil {
		b = true
	}
	if !b {
		profile, _ := wslreg.GetProfileFromName(name)
		switch profile.WsldlTerm {
		case wslreg.FlagWsldlTermWT:
			utils.FreeConsole()
			ExecWindowsTerminal(name)
			os.Exit(0)

		case wslreg.FlagWsldlTermFlute:
			utils.FreeConsole()
			exe := os.Getenv("LOCALAPPDATA")
			exe = utils.DQEscapeString(exe + "\\Microsoft\\WindowsApps\\53621FSApps.FluentTerminal_87x1pks76srcp\\flute.exe")

			cmd := exe + " run " + utils.DQEscapeString(efPath)
			res, err := utils.CreateProcessAndWait(cmd)
			if err != nil {
				utils.AllocConsole()
				fmt.Fprintln(os.Stderr, "ERR: Failed to launch Terminal Process")
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
		Execute(name, nil)
	} else {
		Execute(name, nil)
	}
}
