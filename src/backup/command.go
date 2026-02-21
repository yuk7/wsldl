package backup

import (
	"os"
	"strings"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
)

// GetCommand returns the backup command structure
func GetCommand() cmdline.Command {
	return cmdline.Command{
		Names: []string{"backup"},
		Help: func(distroName string, isListQuery bool) string {
			if wsllib.WslIsDistributionRegistered(distroName) || !isListQuery {
				return getHelpMessage()
			}
			return ""
		},
		Run: execute,
	}
}

// execute is default backup entrypoint
func execute(name string, args []string) error {
	opttar := ""
	optvhdx := ""
	optreg := ""
	switch len(args) {
	case 0:
		_, _, flags, _ := wsllib.WslGetDistributionConfiguration(name)
		if flags&wsllib.FlagEnableWsl2 == wsllib.FlagEnableWsl2 {
			optvhdx = "backup.ext4.vhdx.gz"
			optreg = "backup.reg"
		} else {
			opttar = "backup.tar.gz"
			optreg = "backup.reg"
		}

	case 1:
		arg0Lower := strings.ToLower(args[0])
		switch arg0Lower {
		case "--tar":
			opttar = "backup.tar"
		case "--tgz":
			opttar = "backup.tar.gz"
		case "--vhdx":
			optvhdx = "backup.ext4.vhdx"
		case "--vhdxgz":
			optvhdx = "backup.ext4.vhdx.gz"
		case "--reg":
			optreg = "backup.reg"
		default:
			if strings.HasSuffix(arg0Lower, ".tar") || strings.HasSuffix(arg0Lower, ".tar.gz") || strings.HasSuffix(arg0Lower, ".tgz") {
				opttar = args[0]
			} else if strings.HasSuffix(arg0Lower, ".ext4.vhdx") || strings.HasSuffix(arg0Lower, ".ext4.vhdx.gz") {
				optvhdx = args[0]
			} else if strings.HasSuffix(arg0Lower, ".reg") {
				optreg = args[0]
			} else {
				return utils.NewDisplayError(os.ErrInvalid, true, true, false)
			}
		}

	default:
		return utils.NewDisplayError(os.ErrInvalid, true, true, false)
	}

	if optreg != "" {
		err := backupReg(name, optreg)
		if err != nil {
			return utils.NewDisplayError(err, true, true, false)
		}
	}
	if opttar != "" {
		err := backupTar(name, opttar)
		if err != nil {
			return utils.NewDisplayError(err, true, true, false)
		}

	}
	if optvhdx != "" {
		err := backupExt4Vhdx(name, optvhdx)
		if err != nil {
			return utils.NewDisplayError(err, true, true, false)
		}
	}
	return nil
}
