package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/yuk7/wsldl/backup"
	"github.com/yuk7/wsldl/clean"
	"github.com/yuk7/wsldl/config"
	"github.com/yuk7/wsldl/get"
	"github.com/yuk7/wsldl/help"
	"github.com/yuk7/wsldl/install"
	"github.com/yuk7/wsldl/isregd"
	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
	"github.com/yuk7/wsldl/run"
	"github.com/yuk7/wsldl/version"
)

var (
	executableFunc = os.Executable
	exitFunc       = errutil.Exit
)

// main is the entry point of the application
func main() {
	deps := wsllib.NewDependencies()
	efPath, _ := executableFunc()
	runMain(deps, os.Args, efPath)
}

func runMain(deps wsllib.Dependencies, argv []string, executablePath string) {
	name := filepath.Base(executablePath[:len(executablePath)-len(filepath.Ext(executablePath))])

	var commands = []cmdline.Command{
		isregd.GetCommandWithDeps(deps.Wsl),
		version.GetCommand(),
		install.GetCommandWithNoArgsWithDeps(deps.Wsl, deps.Reg),
		install.GetCommandWithDeps(deps.Wsl, deps.Reg),
		run.GetCommandWithNoArgsWithDeps(deps.Wsl, deps.Reg),
		run.GetCommandWithDeps(deps.Wsl),
		run.GetCommandPWithDeps(deps.Wsl),
		config.GetCommandWithDeps(deps.Wsl, deps.Reg),
		get.GetCommandWithDeps(deps.Wsl, deps.Reg),
		backup.GetCommandWithDeps(deps.Wsl, deps.Reg),
		clean.GetCommandWithDeps(deps.Wsl),
	}

	var helpCommand = help.GetCommand()
	var commandsWithHelp = append(commands, cmdline.Command{
		Names: helpCommand.Names,
		Help:  helpCommand.Help,
		Run: func(distroName string, _ []string) error {
			help.ShowHelpFromCommands(
				append(commands, helpCommand), distroName, argv[2:],
			)
			return nil
		},
	})

	handleCommandError := func(err error) {
		if err == nil {
			return
		}
		var displayErr *errutil.DisplayError
		if errors.As(err, &displayErr) {
			handleDisplayError(displayErr.Err, displayErr.ShowMsg, displayErr.ShowColor, displayErr.Pause)
			return
		}
		var exitCodeErr *errutil.ExitCodeError
		if errors.As(err, &exitCodeErr) {
			exitFunc(exitCodeErr.Pause, exitCodeErr.Code)
			return
		}
		handleDisplayError(err, true, true, false)
	}

	if len(argv) > 1 {
		err := cmdline.RunSubCommand(
			commandsWithHelp,
			func() error {
				return errutil.NewDisplayError(os.ErrInvalid, true, true, false)
			},
			name,
			argv[1:],
		)
		handleCommandError(err)

	} else {
		if !deps.Wsl.IsDistributionRegistered(name) {
			handleCommandError(install.GetCommandWithDeps(deps.Wsl, deps.Reg).Run(name, nil))
		} else {
			handleCommandError(run.GetCommandWithNoArgsWithDeps(deps.Wsl, deps.Reg).Run(name, nil))
		}
	}
}

func handleDisplayError(err error, showMsg bool, showColor bool, pause bool) {
	var errno syscall.Errno

	if showMsg {
		formatted := errutil.FormatError(err)
		if showColor {
			errutil.ErrorRedPrintln(formatted)
		} else {
			fmt.Fprintln(os.Stderr, formatted)
		}
	}

	if err == nil {
		exitFunc(pause, 1)
	}
	if errors.As(err, &errno) {
		if showMsg {
			fmt.Fprintf(os.Stderr, "HRESULT: 0x%x\n", int(errno))
		}
		exitFunc(pause, int(errno))
	} else if err == os.ErrInvalid {
		if showMsg {
			efPath, _ := executableFunc()
			exeName := filepath.Base(efPath)
			fmt.Fprintln(os.Stderr, "Your command may be incorrect.")
			fmt.Fprintf(os.Stderr, "You can get help with `%s help`.\n", exeName)
		}
	} else if strings.HasPrefix(fmt.Sprintf("%#v", err), "&errors.errorString{") {
		if showMsg {
			fmt.Fprintf(os.Stderr, "%#v\n", err)
		}
	}
	exitFunc(pause, 1)
}
