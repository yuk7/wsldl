package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/yuk7/wsldl/get"
	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/run"
	"github.com/yuk7/wsllib-go"
	wslreg "github.com/yuk7/wslreglib-go"
)

// GetCommand returns the config set command structure
func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"config", "set"},
		Help: func(distroName string, isListQuery bool) string {
			if wsllib.WslIsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			}
			return ""
		},
		Run: func(distroName string, args []string) {
			Execute(distroName, args)
		},
	}
}

// Execute is default install entrypoint
func Execute(name string, args []string) {
	var err error
	uid, flags := get.WslGetConfig(name)
	if len(args) == 2 {
		switch args[0] {
		case "--default-uid":
			var intUID int
			intUID, err = strconv.Atoi(args[1])
			uid = uint64(intUID)

		case "--default-user":
			str, _, errtmp := run.ExecRead(name, "id -u "+utils.DQEscapeString(args[1]))
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
					err = wslreg.SetWslVersion(name, intWslVer)
				} else {
					err = os.ErrInvalid
					break
				}
			}

		case "--default-term":
			value := 0
			switch args[1] {
			case "default", strconv.Itoa(wslreg.FlagWsldlTermDefault):
				value = wslreg.FlagWsldlTermDefault
			case "wt", strconv.Itoa(wslreg.FlagWsldlTermWT):
				value = wslreg.FlagWsldlTermWT
			case "flute", strconv.Itoa(wslreg.FlagWsldlTermFlute):
				value = wslreg.FlagWsldlTermFlute
			default:
				err = os.ErrInvalid
				break
			}
			profile, err := wslreg.GetProfileFromName(name)
			if err != nil {
				break
			}
			profile.WsldlTerm = value
			err = wslreg.WriteProfile(profile)

		case "--flags-val":
			var intFlags int
			intFlags, err = strconv.Atoi(args[1])
			flags = uint32(intFlags)

		default:
			err = os.ErrInvalid
		}
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
		wsllib.WslConfigureDistribution(name, uid, flags)
	} else {
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}
}
