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

type configOption int

const (
	configOptionDefaultUID configOption = iota
	configOptionDefaultUser
	configOptionAppendPath
	configOptionMountDrive
	configOptionWslVersion
	configOptionDefaultTerm
	configOptionFlagsVal
)

type configOptions struct {
	option      configOption
	uid         uint64
	user        string
	enabled     bool
	wslVersion  int
	defaultTerm int
	flags       uint32
}

// GetCommand returns the config set command structure
func GetCommand() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithDeps(deps.Wsl, deps.Reg)
}

// GetCommandWithDeps returns the config set command structure with injectable dependencies.
func GetCommandWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg) cmdline.Command {
	return cmdline.Command{
		Names: []string{"config", "set"},
		Visible: func(distroName string) bool {
			return wsl.IsDistributionRegistered(distroName)
		},
		HelpText: getHelpMessage,
		Run: func(name string, args []string) error {
			return execute(wsl, reg, name, args)
		},
	}
}

// execute is default install entrypoint
func execute(wsl wsllib.WslLib, reg wsllib.WslReg, name string, args []string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	return executeWithOptions(wsl, reg, name, opts)
}

func parseArgs(args []string) (configOptions, error) {
	if len(args) != 2 {
		return configOptions{}, os.ErrInvalid
	}

	opts := configOptions{}
	switch args[0] {
	case "--default-uid":
		intUID, err := strconv.Atoi(args[1])
		if err != nil {
			return configOptions{}, err
		}
		opts.option = configOptionDefaultUID
		opts.uid = uint64(intUID)

	case "--default-user":
		opts.option = configOptionDefaultUser
		opts.user = args[1]

	case "--append-path":
		b, err := strconv.ParseBool(args[1])
		if err != nil {
			return configOptions{}, err
		}
		opts.option = configOptionAppendPath
		opts.enabled = b

	case "--mount-drive":
		b, err := strconv.ParseBool(args[1])
		if err != nil {
			return configOptions{}, err
		}
		opts.option = configOptionMountDrive
		opts.enabled = b

	case "--wsl-version":
		intWslVer, err := strconv.Atoi(args[1])
		if err != nil {
			return configOptions{}, err
		}
		if intWslVer != 1 && intWslVer != 2 {
			return configOptions{}, os.ErrInvalid
		}
		opts.option = configOptionWslVersion
		opts.wslVersion = intWslVer

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
			return configOptions{}, os.ErrInvalid
		}
		opts.option = configOptionDefaultTerm
		opts.defaultTerm = value

	case "--flags-val":
		intFlags, err := strconv.Atoi(args[1])
		if err != nil {
			return configOptions{}, err
		}
		opts.option = configOptionFlagsVal
		opts.flags = uint32(intFlags)

	default:
		return configOptions{}, os.ErrInvalid
	}

	return opts, nil
}

func executeWithOptions(wsl wsllib.WslLib, reg wsllib.WslReg, name string, opts configOptions) error {
	return executeWithOptionsAndExecRead(wsl, reg, name, opts, wslexec.ExecRead)
}

func executeWithOptionsAndExecRead(
	wsl wsllib.WslLib,
	reg wsllib.WslReg,
	name string,
	opts configOptions,
	execRead func(wsllib.WslLib, string, string) (string, uint32, error),
) error {
	uid, flags, err := wslapi.GetConfig(wsl, name)
	if err != nil {
		errutil.ErrorRedPrintln("ERR: Failed to GetDistributionConfiguration")
		return errutil.NewDisplayError(err, true, true, false)
	}

	switch opts.option {
	case configOptionDefaultUID:
		uid = opts.uid

	case configOptionDefaultUser:
		str, _, errtmp := execRead(wsl, name, "id -u "+fileutil.DQEscapeString(opts.user))
		err = errtmp
		if err == nil {
			intUID, convErr := strconv.Atoi(str)
			err = convErr
			uid = uint64(intUID)
			if err != nil {
				err = errors.New(str)
			}
		}

	case configOptionAppendPath:
		flags = updateFlag(flags, wsllib.FlagAppendNTPath, opts.enabled)

	case configOptionMountDrive:
		flags = updateFlag(flags, wsllib.FlagEnableDriveMounting, opts.enabled)

	case configOptionWslVersion:
		err = reg.SetWslVersion(name, opts.wslVersion)

	case configOptionDefaultTerm:
		profile, profileErr := reg.GetProfileFromName(name)
		err = profileErr
		if err != nil {
			break
		}
		profile.WsldlTerm = opts.defaultTerm
		err = reg.WriteProfile(profile)

	case configOptionFlagsVal:
		flags = opts.flags

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

	return nil
}

func updateFlag(flags uint32, mask uint32, enabled bool) uint32 {
	if enabled {
		return flags | mask
	}
	return flags &^ mask
}
