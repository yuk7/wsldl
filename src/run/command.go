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

type runOptions struct {
	commandArgs []string
	inheritPath bool
}

type runPOptions struct {
	commandArgs []string
}

type runNoArgsOptions struct{}

type runNoArgsDeps struct {
	mustExecutable        func() string
	stat                  func(string) (os.FileInfo, error)
	isInstalledFilesExist func() bool
	readInput             func() string
	repairRegistry        func(wsllib.WslReg, wsllib.Profile) error
	isParentConsole       func() (bool, error)
	freeConsole           func() error
	allocConsole          func()
	setConsoleTitle       func(string)
	execWindowsTerminal   func(wsllib.WslReg, string) error
	getenv                func(string) string
	createProcessAndWait  func(string) (int, error)
	execute               func(wsllib.WslLib, string, []string) error
}

func defaultRunNoArgsDeps() runNoArgsDeps {
	return runNoArgsDeps{
		mustExecutable:        errutil.MustExecutable,
		stat:                  os.Stat,
		isInstalledFilesExist: repair.IsInstalledFilesExist,
		readInput: func() string {
			var in string
			fmt.Scan(&in)
			return in
		},
		repairRegistry:       repairRegistry,
		isParentConsole:      console.IsParentConsole,
		freeConsole:          console.FreeConsole,
		allocConsole:         console.AllocConsole,
		setConsoleTitle:      console.SetConsoleTitle,
		execWindowsTerminal:  ExecWindowsTerminal,
		getenv:               os.Getenv,
		createProcessAndWait: utils.CreateProcessAndWait,
		execute:              execute,
	}
}

// GetCommandWithNoArgs returns the run command structure with no arguments
func GetCommandWithNoArgs() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithNoArgsWithDeps(deps.Wsl, deps.Reg)
}

// GetCommandWithNoArgsWithDeps returns the run command structure with no arguments and injectable dependencies.
func GetCommandWithNoArgsWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg) cmdline.Command {
	return cmdline.Command{
		Visible: func(distroName string) bool {
			return wsl.IsDistributionRegistered(distroName)
		},
		HelpText: getHelpMessageNoArgs,
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
		Visible: func(distroName string) bool {
			return wsl.IsDistributionRegistered(distroName)
		},
		HelpText: getHelpMessage,
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
		Visible: func(distroName string) bool {
			return wsl.IsDistributionRegistered(distroName)
		},
		HelpText: getHelpMessageP,
		Run: func(name string, args []string) error {
			return executeP(wsl, name, args)
		},
	}
}

func parseRunArgs(args []string) (runOptions, error) {
	opts := runOptions{
		inheritPath: true,
	}
	if args == nil {
		opts.inheritPath = !fileutil.IsCurrentDirSpecial()
		return opts, nil
	}
	opts.commandArgs = append(opts.commandArgs, args...)
	return opts, nil
}

func parseRunPArgs(args []string) (runPOptions, error) {
	opts := runPOptions{}
	opts.commandArgs = append(opts.commandArgs, args...)
	return opts, nil
}

func parseRunNoArgs(args []string) (runNoArgsOptions, error) {
	if len(args) != 0 {
		return runNoArgsOptions{}, os.ErrInvalid
	}
	return runNoArgsOptions{}, nil
}

// execute is default run entrypoint.
func execute(wsl wsllib.WslLib, name string, args []string) error {
	opts, err := parseRunArgs(args)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	return executeWithOptions(wsl, name, opts)
}

func executeWithOptions(wsl wsllib.WslLib, name string, opts runOptions) error {
	command := ""
	for _, s := range opts.commandArgs {
		command = command + " " + fileutil.DQEscapeString(s)
	}

	exitCode, err := wsl.LaunchInteractive(name, command, opts.inheritPath)
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
	opts, err := parseRunPArgs(args)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	return executePWithOptions(wsl, name, opts)
}

func executePWithOptions(wsl wsllib.WslLib, name string, opts runPOptions) error {
	return executePWithOptionsWithExecRead(wsl, name, opts, wslexec.ExecRead)
}

func executePWithOptionsWithExecRead(
	wsl wsllib.WslLib,
	name string,
	opts runPOptions,
	execRead func(wsllib.WslLib, string, string) (string, uint32, error),
) error {
	var convArgs []string
	for _, s := range opts.commandArgs {
		if strings.Contains(s, "\\") {
			s = strings.Replace(s, "\\", "/", -1)
			s = fileutil.DQEscapeString(s)
			out, exitCode, err := execRead(wsl, name, "wslpath -u "+s)
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
	opts, err := parseRunNoArgs(args)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	return executeNoArgsWithOptions(wsl, reg, name, opts)
}

func executeNoArgsWithOptions(wsl wsllib.WslLib, reg wsllib.WslReg, name string, opts runNoArgsOptions) error {
	return executeNoArgsWithOptionsAndDeps(wsl, reg, name, opts, defaultRunNoArgsDeps())
}

func executeNoArgsWithOptionsAndDeps(wsl wsllib.WslLib, reg wsllib.WslReg, name string, _ runNoArgsOptions, deps runNoArgsDeps) error {
	efPath := deps.mustExecutable()
	profile, _ := reg.GetProfileFromName(name)

	// repair when the installation is moved
	if profile.BasePath != "" {
		_, err := deps.stat(profile.BasePath)
		if os.IsNotExist(err) {
			if deps.isInstalledFilesExist() {
				fmt.Printf("This instance (%s) BasePath is not exist.\n", name)
				fmt.Printf("Do you want to repair the installation information?\n")
				fmt.Printf("Type y/n:")
				in := deps.readInput()

				if in == "y" {
					err := deps.repairRegistry(reg, profile)
					if err != nil {
						return errutil.NewDisplayError(err, true, true, true)
					}
					errutil.StdoutGreenPrintln("done.")
					return errutil.NewExitCodeError(0, true)
				}
			}
		}
	}

	b, err := deps.isParentConsole()
	if err != nil {
		b = true
	}
	if !b {
		switch profile.WsldlTerm {
		case wsllib.FlagWsldlTermWT:
			_ = deps.freeConsole()
			return deps.execWindowsTerminal(reg, name)

		case wsllib.FlagWsldlTermFlute:
			_ = deps.freeConsole()
			exe := deps.getenv("LOCALAPPDATA")
			exe = fileutil.DQEscapeString(exe + "\\Microsoft\\WindowsApps\\53621FSApps.FluentTerminal_87x1pks76srcp\\flute.exe")

			cmd := exe + " run " + fileutil.DQEscapeString(efPath+" run")
			res, err := deps.createProcessAndWait(cmd)
			if err != nil {
				deps.allocConsole()
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

		deps.setConsoleTitle(name)
		return deps.execute(wsl, name, nil)
	} else {
		return deps.execute(wsl, name, nil)
	}
}
