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
		Run: func(distroName string, args []string) {
			Execute(distroName, args)
		},
	}
}

// Execute is default install entrypoint
func Execute(name string, args []string) {
	uid, flags := WslGetConfig(name)
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
					utils.ErrorExit(proferr, true, true, false)
				}
				utils.ErrorExit(errors.New("lxguid get failed"), true, true, false)
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
				utils.ErrorExit(err, true, true, false)
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
				utils.ErrorExit(errors.New("profile not found"), true, true, false)
			}

		case "--flags-val":
			print(flags)

		case "--flags-bits":
			fmt.Printf("%04b", flags)

		default:
			utils.ErrorExit(os.ErrInvalid, true, true, false)
		}
	} else {
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}
}
