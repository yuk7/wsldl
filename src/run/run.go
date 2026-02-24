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

// ExecWindowsTerminal executes Windows Terminal
func ExecWindowsTerminal(reg wsllib.WslReg, name string) error {
	// Get the name from the registry to be case sensitive.
	profile, _ := reg.GetProfileFromName(name)
	if profile.DistributionName != "" {
		name = profile.DistributionName
	}

	profileName := ""
	conf, err := wtutils.ReadParseWTConfig()
	if err == nil {
		guid := "{" + wtutils.CreateProfileGUID(name) + "}"
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

	exe := os.Getenv("LOCALAPPDATA")
	exe = fileutil.DQEscapeString(exe + "\\Microsoft\\WindowsApps\\" + wtutils.WTPackageName + "\\wt.exe")
	cmd := exe

	if profileName != "" {
		cmd = cmd + " -p " + fileutil.DQEscapeString(profileName)
	} else {
		efPath, _ := os.Executable()
		cmd = cmd + " " + fileutil.DQEscapeString(efPath) + " run"
	}

	res, err := utils.CreateProcessAndWait(cmd)
	if err != nil {
		console.AllocConsole()
		fmt.Fprintln(os.Stderr, "ERR: Failed to launch the terminal process")
		fmt.Fprintln(os.Stderr, exe)
		return errutil.NewDisplayError(err, true, false, true)
	}
	if res != 0 {
		return errutil.NewExitCodeError(res, false)
	}
	return nil
}
