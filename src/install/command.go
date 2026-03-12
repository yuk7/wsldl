package install

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/fileutil"
	"github.com/yuk7/wsldl/lib/preset"
	"github.com/yuk7/wsldl/lib/repair"
	"github.com/yuk7/wsldl/lib/wsllib"
)

type installMode int

const (
	installModeAuto installMode = iota
	installModeRoot
	installModePath
)

type installArgs struct {
	mode          installMode
	inputPath     string
	fromNoArgCall bool
}

type installOptions struct {
	rootPath       string
	rootFileSHA256 string
	showProgress   bool
	pauseAfterRun  bool
	presetVersion  int
}

type installCommandDeps struct {
	isInstalledFilesExist func() bool
	readInput             func() string
	repairRegistry        func(wsllib.WslReg, string) error
	install               func(context.Context, wsllib.WslLib, wsllib.WslReg, string, string, string, bool) error
	setVersion            func(string, int) error
	waitForEnter          func()
}

var detectRootfsFilesFunc = detectRootfsFiles

func defaultInstallCommandDeps() installCommandDeps {
	return installCommandDeps{
		isInstalledFilesExist: repair.IsInstalledFilesExist,
		readInput: func() string {
			var in string
			fmt.Scan(&in)
			return in
		},
		repairRegistry: repairRegistry,
		install:        Install,
		setVersion: func(name string, version int) error {
			wslexe := filepath.Join(fileutil.GetWindowsDirectory(), "System32", "wsl.exe")
			_, err := exec.Command(wslexe, "--set-version", name, strconv.Itoa(version)).Output()
			return err
		},
		waitForEnter: func() {
			fmt.Fprintf(os.Stdout, "Press enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		},
	}
}

func GetCommandWithNoArgs() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithNoArgsWithDeps(deps.Wsl, deps.Reg)
}

func GetCommandWithNoArgsWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg) cmdline.Command {
	return cmdline.Command{
		Visible: func(distroName string) bool {
			return !wsl.IsDistributionRegistered(distroName)
		},
		HelpText: getHelpMessageNoArgs,
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
		Visible: func(distroName string) bool {
			return !wsl.IsDistributionRegistered(distroName)
		},
		HelpText: getHelpMessage,
		Run: func(name string, args []string) error {
			return execute(wsl, reg, name, args)
		},
	}
}

// execute is default install entrypoint
func execute(wsl wsllib.WslLib, reg wsllib.WslReg, name string, args []string) error {
	if wsl.IsDistributionRegistered(name) {
		errutil.ErrorRedPrintln("ERR: [" + name + "] is already installed.")
		return errutil.NewDisplayError(os.ErrInvalid, false, true, false)
	}

	parsed, err := parseArgs(args)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}

	opts := resolveOptions(parsed)
	return executeWithOptions(wsl, reg, name, opts)
}

func parseArgs(args []string) (installArgs, error) {
	parsed := installArgs{fromNoArgCall: args == nil}

	switch len(args) {
	case 0:
		parsed.mode = installModeAuto
	case 1:
		if args[0] == "--root" {
			parsed.mode = installModeRoot
		} else {
			parsed.mode = installModePath
			parsed.inputPath = args[0]
		}
	default:
		return installArgs{}, os.ErrInvalid
	}

	return parsed, nil
}

func resolveOptions(parsed installArgs) installOptions {
	jsonPreset, _ := preset.ReadParsePreset()

	opts := installOptions{
		showProgress:  parsed.mode == installModeAuto,
		pauseAfterRun: parsed.fromNoArgCall,
		presetVersion: jsonPreset.WslVersion,
	}

	switch parsed.mode {
	case installModeAuto, installModeRoot:
		rootPath := "rootfs.tar.gz"
		if detectedRootPath, err := detectRootfsFilesFunc(); err == nil {
			rootPath = detectedRootPath
		}
		if jsonPreset.InstallFile != "" {
			rootPath = jsonPreset.InstallFile
		}
		opts.rootPath = rootPath
		opts.rootFileSHA256 = jsonPreset.InstallFileSha256
	case installModePath:
		opts.rootPath = parsed.inputPath
	}

	return opts
}

func executeWithOptions(wsl wsllib.WslLib, reg wsllib.WslReg, name string, opts installOptions) error {
	return executeWithOptionsAndDeps(wsl, reg, name, opts, defaultInstallCommandDeps())
}

func executeWithOptionsAndDeps(wsl wsllib.WslLib, reg wsllib.WslReg, name string, opts installOptions, deps installCommandDeps) error {
	if opts.pauseAfterRun {
		if deps.isInstalledFilesExist() {
			fmt.Printf("An old installation files has been found.\n")
			fmt.Printf("Do you want to rewrite and repair the installation information?\n")
			fmt.Printf("Type y/n:")
			in := deps.readInput()

			if in == "y" {
				err := deps.repairRegistry(reg, name)
				if err != nil {
					return errutil.NewDisplayError(err, opts.showProgress, true, opts.showProgress)
				}
				errutil.StdoutGreenPrintln("done.")
				return nil
			}
		}
	}

	err := deps.install(context.Background(), wsl, reg, name, opts.rootPath, opts.rootFileSHA256, opts.showProgress)
	if err != nil {
		return errutil.NewDisplayError(err, opts.showProgress, true, opts.pauseAfterRun)
	}

	if opts.presetVersion == 1 || opts.presetVersion == 2 {
		err = deps.setVersion(name, opts.presetVersion)
	}

	if err == nil {
		if opts.showProgress {
			errutil.StdoutGreenPrintln("Installation complete")
		}
	} else {
		return errutil.NewDisplayError(err, opts.showProgress, true, opts.pauseAfterRun)
	}

	if opts.pauseAfterRun {
		deps.waitForEnter()
	}

	return nil
}
