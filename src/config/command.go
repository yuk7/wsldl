package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/fileutil"
	"github.com/yuk7/wsldl/lib/wslapi"
	"github.com/yuk7/wsldl/lib/wslexec"
	"github.com/yuk7/wsldl/lib/wsllib"
)

// GetCommand returns the config set command structure
func GetCommand() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithDeps(deps.Wsl, deps.Reg)
}

// GetCommandWithDeps returns the config set command structure with injectable dependencies.
func GetCommandWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg) cmdline.Command {
	return cmdline.Command{
		Names: []string{"config", "set"},
		Help: func(distroName string, isListQuery bool) string {
			if wsl.IsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			}
			return ""
		},
		Run: func(name string, args []string) error {
			return execute(wsl, reg, name, args)
		},
	}
}

// execute is default install entrypoint
func execute(wsl wsllib.WslLib, reg wsllib.WslReg, name string, args []string) error {
	uid, flags, err := wslapi.GetConfig(wsl, name)
	if err != nil {
		errutil.ErrorRedPrintln("ERR: Failed to GetDistributionConfiguration")
		return errutil.NewDisplayError(err, true, true, false)
	}
	if len(args) == 2 {
		switch args[0] {
		case "--default-uid":
			var intUID int
			intUID, err = strconv.Atoi(args[1])
			uid = uint64(intUID)

		case "--default-user":
			str, _, errtmp := wslexec.ExecRead(wsl, name, "id -u "+fileutil.DQEscapeString(args[1]))
			err = errtmp
			if err == nil {
				var intUID int
				intUID, err = strconv.Atoi(str)
				uid = uint64(intUID)
				if err != nil {
					err = errors.New(str)
				}
			}

		case "--append-path":
			var b bool
			b, err = strconv.ParseBool(args[1])
			if b {
				flags |= wsllib.FlagAppendNTPath
			} else {
				flags ^= wsllib.FlagAppendNTPath
			}

		case "--mount-drive":
			var b bool
			b, err = strconv.ParseBool(args[1])
			if b {
				flags |= wsllib.FlagEnableDriveMounting
			} else {
				flags ^= wsllib.FlagEnableDriveMounting
			}

		case "--wsl-version":
			var intWslVer int
			intWslVer, err = strconv.Atoi(args[1])
			if err == nil {
				if intWslVer == 1 || intWslVer == 2 {
					err = reg.SetWslVersion(name, intWslVer)
				} else {
					err = os.ErrInvalid
					break
				}
			}

		case "--default-term":
			value := 0
			switch args[1] {
			case "default", strconv.Itoa(wsllib.FlagWsldlTermDefault):
				value = wsllib.FlagWsldlTermDefault
			case "wt", strconv.Itoa(wsllib.FlagWsldlTermWT):
				value = wsllib.FlagWsldlTermWT
			case "flute", strconv.Itoa(wsllib.FlagWsldlTermFlute):
				value = wsllib.FlagWsldlTermFlute
			default:
				err = os.ErrInvalid
				break
			}
			profile, err := reg.GetProfileFromName(name)
			if err != nil {
				break
			}
			profile.WsldlTerm = value
			err = reg.WriteProfile(profile)

		case "--flags-val":
			var intFlags int
			intFlags, err = strconv.Atoi(args[1])
			flags = uint32(intFlags)

		default:
			err = os.ErrInvalid
		}
		if err != nil {
			return errutil.NewDisplayError(err, true, true, false)
		}
		err = wsl.ConfigureDistribution(name, uid, flags)
		if err != nil {
			return errutil.NewDisplayError(err, true, true, false)
		}
	} else {
		return errutil.NewDisplayError(os.ErrInvalid, true, true, false)
	}
	return nil
}
