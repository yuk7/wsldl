package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuk7/wsldl/lib/console"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/fileutil"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wsllib"
	"github.com/yuk7/wsldl/lib/wtutils"
)

type execWindowsTerminalDeps struct {
	readParseWTConfig    func() (wtutils.Config, error)
	createProfileGUID    func(string) string
	getenv               func(string) string
	mustExecutable       func() string
	createProcessAndWait func(string) (int, error)
	allocConsole         func()
}

func defaultExecWindowsTerminalDeps() execWindowsTerminalDeps {
	return execWindowsTerminalDeps{
		readParseWTConfig:    wtutils.ReadParseWTConfig,
		createProfileGUID:    wtutils.CreateProfileGUID,
		getenv:               os.Getenv,
		mustExecutable:       errutil.MustExecutable,
		createProcessAndWait: utils.CreateProcessAndWait,
		allocConsole:         console.AllocConsole,
	}
}

// ExecWindowsTerminal executes Windows Terminal
func ExecWindowsTerminal(reg wsllib.WslReg, name string) error {
	return execWindowsTerminalWithDeps(reg, name, defaultExecWindowsTerminalDeps())
}

func execWindowsTerminalWithDeps(reg wsllib.WslReg, name string, deps execWindowsTerminalDeps) error {
	// Get the name from the registry to be case sensitive.
	profile, _ := reg.GetProfileFromName(name)
	if profile.DistributionName != "" {
		name = profile.DistributionName
	}

	profileName := ""
	conf, err := deps.readParseWTConfig()
	if err == nil {
		guid := "{" + deps.createProfileGUID(name) + "}"
		for _, profile := range conf.Profiles.ProfileList {
			if profile.GUID == guid {
				profileName = profile.Name
				break
			}
		}
		if profileName == "" {
			for _, profile := range conf.Profiles.ProfileList {
				if strings.EqualFold(profile.Name, name) {
					profileName = profile.Name
					break
				}
			}
		}
	}

	exe := deps.getenv("LOCALAPPDATA")
	exe = fileutil.DQEscapeString(exe + "\\Microsoft\\WindowsApps\\" + wtutils.WTPackageName + "\\wt.exe")
	cmd := exe

	if profileName != "" {
		cmd = cmd + " -p " + fileutil.DQEscapeString(profileName)
	} else {
		efPath := deps.mustExecutable()
		cmd = cmd + " " + fileutil.DQEscapeString(efPath) + " run"
	}

	res, err := deps.createProcessAndWait(cmd)
	if err != nil {
		deps.allocConsole()
		fmt.Fprintln(os.Stderr, "ERR: Failed to launch the terminal process")
		fmt.Fprintln(os.Stderr, exe)
		return errutil.NewDisplayError(err, true, false, true)
	}
	if res != 0 {
		return errutil.NewExitCodeError(res, false)
	}
	return nil
}
