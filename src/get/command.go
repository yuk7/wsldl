package get

import (
	"errors"
	"fmt"
	"os"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wslapi"
	"github.com/yuk7/wsldl/lib/wsllib"
	"github.com/yuk7/wsldl/lib/wtutils"
)

type getOption int

const (
	getOptionDefaultUID getOption = iota
	getOptionAppendPath
	getOptionMountDrive
	getOptionWslVersion
	getOptionLXGuid
	getOptionDefaultTerm
	getOptionWTProfileName
	getOptionFlagsVal
	getOptionFlagsBits
)

type getOptions struct {
	option getOption
}

// GetCommand returns the get command structure
func GetCommand() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithDeps(deps.Wsl, deps.Reg)
}

// GetCommandWithDeps returns the get command structure with injectable dependencies.
func GetCommandWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg) cmdline.Command {
	return cmdline.Command{
		Names: []string{"get"},
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
	return executeWithWTConfigReader(wsl, reg, name, args, wtutils.ReadParseWTConfig)
}

func parseArgs(args []string) (getOptions, error) {
	if len(args) != 1 {
		return getOptions{}, os.ErrInvalid
	}

	switch args[0] {
	case "--default-uid":
		return getOptions{option: getOptionDefaultUID}, nil
	case "--append-path":
		return getOptions{option: getOptionAppendPath}, nil
	case "--mount-drive":
		return getOptions{option: getOptionMountDrive}, nil
	case "--wsl-version":
		return getOptions{option: getOptionWslVersion}, nil
	case "--lxguid", "--lxuid":
		return getOptions{option: getOptionLXGuid}, nil
	case "--default-term", "--default-terminal":
		return getOptions{option: getOptionDefaultTerm}, nil
	case "--wt-profile-name", "--wt-profilename", "--wt-pn":
		return getOptions{option: getOptionWTProfileName}, nil
	case "--flags-val":
		return getOptions{option: getOptionFlagsVal}, nil
	case "--flags-bits":
		return getOptions{option: getOptionFlagsBits}, nil
	default:
		return getOptions{}, os.ErrInvalid
	}
}

func executeWithWTConfigReader(
	wsl wsllib.WslLib,
	reg wsllib.WslReg,
	name string,
	args []string,
	readWTConfig func() (wtutils.Config, error),
) error {
	opts, err := parseArgs(args)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	return executeWithOptions(wsl, reg, name, opts, readWTConfig)
}

func executeWithOptions(
	wsl wsllib.WslLib,
	reg wsllib.WslReg,
	name string,
	opts getOptions,
	readWTConfig func() (wtutils.Config, error),
) error {
	uid, flags, err := wslapi.GetConfig(wsl, name)
	if err != nil {
		errutil.ErrorRedPrintln("ERR: Failed to GetDistributionConfiguration")
		return errutil.NewDisplayError(err, true, true, false)
	}
	profile, proferr := reg.GetProfileFromName(name)

	switch opts.option {
	case getOptionDefaultUID:
		print(uid)

	case getOptionAppendPath:
		print(flags&wsllib.FlagAppendNTPath == wsllib.FlagAppendNTPath)

	case getOptionMountDrive:
		print(flags&wsllib.FlagEnableDriveMounting == wsllib.FlagEnableDriveMounting)

	case getOptionWslVersion:
		if flags&wsllib.FlagEnableWsl2 == wsllib.FlagEnableWsl2 {
			print("2")
		} else {
			print("1")
		}

	case getOptionLXGuid:
		if profile.UUID == "" {
			if proferr != nil {
				return errutil.NewDisplayError(proferr, true, true, false)
			}
			return errutil.NewDisplayError(errors.New("lxguid get failed"), true, true, false)
		}
		print(profile.UUID)

	case getOptionDefaultTerm:
		switch profile.WsldlTerm {
		case wsllib.FlagWsldlTermWT:
			print("wt")
		case wsllib.FlagWsldlTermFlute:
			print("flute")
		default:
			print("default")
		}

	case getOptionWTProfileName:
		if profile.DistributionName != "" {
			name = profile.DistributionName
		}

		conf, err := readWTConfig()
		if err != nil {
			return errutil.NewDisplayError(err, true, true, false)
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
			return errutil.NewDisplayError(errors.New("profile not found"), true, true, false)
		}

	case getOptionFlagsVal:
		print(flags)

	case getOptionFlagsBits:
		fmt.Printf("%04b", flags)

	default:
		return errutil.NewDisplayError(os.ErrInvalid, true, true, false)
	}

	return nil
}
