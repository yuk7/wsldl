package run

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wslreg"
	"github.com/yuk7/wsldl/lib/wtutils"
	"github.com/yuk7/wsllib-go"
)

//ExecRead execs command and read output
func ExecRead(name, command string) (out string, exitCode uint32, err error) {
	stdin := syscall.Handle(0)
	stdout := syscall.Handle(0)
	stdintmp := syscall.Handle(0)
	stdouttmp := syscall.Handle(0)
	sa := syscall.SecurityAttributes{InheritHandle: 1, SecurityDescriptor: 0}

	syscall.CreatePipe(&stdin, &stdintmp, &sa, 0)
	syscall.CreatePipe(&stdout, &stdouttmp, &sa, 0)

	handle, err := wsllib.WslLaunch(name, command, true, stdintmp, stdouttmp, stdouttmp)
	syscall.WaitForSingleObject(handle, syscall.INFINITE)
	syscall.GetExitCodeProcess(handle, &exitCode)
	buf := make([]byte, syscall.MAX_LONG_PATH)
	var length uint32

	syscall.ReadFile(stdout, buf, &length, nil)

	//[]byte -> string and cut to fit the length
	out = string(buf)[:length]
	if out[len(out)-1:] == "\n" {
		out = out[:len(out)-1]
	}
	return
}

// ExecWindowsTerminal executes Windows Terminal
func ExecWindowsTerminal(name string) {
	// Get the name from the registry to be case sensitive.
	profile, _ := wslreg.GetProfileFromName(name)
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
	exe = utils.DQEscapeString(exe + "\\Microsoft\\WindowsApps\\" + wtutils.WTPackageName + "\\wt.exe")
	cmd := exe

	if profileName != "" {
		cmd = cmd + " -p " + utils.DQEscapeString(profileName)
	} else {
		efPath, _ := os.Executable()
		cmd = cmd + " " + utils.DQEscapeString(efPath) + " run"
	}

	res, err := utils.CreateProcessAndWait(cmd)
	if err != nil {
		utils.AllocConsole()
		fmt.Fprintln(os.Stderr, "ERR: Failed to launch Terminal Process")
		fmt.Fprintln(os.Stderr, exe)
		utils.ErrorExit(err, true, false, true)
	}
	os.Exit(res)
}
