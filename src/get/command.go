package get

import (
	"errors"
	"fmt"
	"os"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wtutils"
	"github.com/yuk7/wsllib-go"
	wslreg "github.com/yuk7/wslreglib-go"
)

// GetCommand returns the get command structure
func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"get"},
		Help: func(distroName string, isListQuery bool) string {
			if wsllib.WslIsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			}
			return ""
		},
		Run: execute,
	}
}

// execute is default install entrypoint
func execute(name string, args []string) error {
	uid, flags, err := WslGetConfig(name)
	if err != nil {
		utils.ErrorRedPrintln("ERR: Failed to GetDistributionConfiguration")
		return utils.NewDisplayError(err, true, true, false)
	}
	profile, proferr := wslreg.GetProfileFromName(name)
	if len(args) == 1 {
		switch args[0] {
		case "--default-uid":
			print(uid)

		case "--append-path":
			print(flags&wsllib.FlagAppendNTPath == wsllib.FlagAppendNTPath)

		case "--mount-drive":
			print(flags&wsllib.FlagEnableDriveMounting == wsllib.FlagEnableDriveMounting)

		case "--wsl-version":
			if flags&wsllib.FlagEnableWsl2 == wsllib.FlagEnableWsl2 {
				print("2")
			} else {
				print("1")
			}

		case "--lxguid", "--lxuid":
			if profile.UUID == "" {
				if proferr != nil {
					return utils.NewDisplayError(proferr, true, true, false)
				}
				return utils.NewDisplayError(errors.New("lxguid get failed"), true, true, false)
			}
			print(profile.UUID)

		case "--default-term", "--default-terminal":
			switch profile.WsldlTerm {
			case wslreg.FlagWsldlTermWT:
				print("wt")
			case wslreg.FlagWsldlTermFlute:
				print("flute")
			default:
				print("default")
			}

		case "--wt-profile-name", "--wt-profilename", "--wt-pn":
			if profile.DistributionName != "" {
				name = profile.DistributionName
			}

			conf, err := wtutils.ReadParseWTConfig()
			if err != nil {
				return utils.NewDisplayError(err, true, true, false)
			}
			guid := "{" + wtutils.CreateProfileGUID(name) + "}"
			profileName := ""
			for _, profile := range conf.Profiles.ProfileList {
				if profile.GUID == guid {
					profileName = profile.Name
					break
				}
			}
			if profileName != "" {
				print(profileName)
			} else {
				return utils.NewDisplayError(errors.New("profile not found"), true, true, false)
			}

		case "--flags-val":
			print(flags)

		case "--flags-bits":
			fmt.Printf("%04b", flags)

		default:
			return utils.NewDisplayError(os.ErrInvalid, true, true, false)
		}
	} else {
		return utils.NewDisplayError(os.ErrInvalid, true, true, false)
	}
	return nil
}
