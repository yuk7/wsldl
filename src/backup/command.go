package backup

import (
	"os"
	"strings"

	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

type backupOptions struct {
	auto     bool
	tarPath  string
	vhdxPath string
	regPath  string
}

// GetCommand returns the backup command structure
func GetCommand() cmdline.Command {
	deps := wsllib.NewDependencies()
	return GetCommandWithDeps(deps.Wsl, deps.Reg)
}

// GetCommandWithDeps returns the backup command structure with injectable dependencies.
func GetCommandWithDeps(wsl wsllib.WslLib, reg wsllib.WslReg) cmdline.Command {
	return cmdline.Command{
		Names: []string{"backup"},
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

// execute is default backup entrypoint
func execute(wsl wsllib.WslLib, reg wsllib.WslReg, name string, args []string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	return executeWithBackupsOptions(wsl, reg, name, opts, backupReg, backupTar, backupExt4Vhdx)
}

func parseArgs(args []string) (backupOptions, error) {
	opts := backupOptions{}
	switch len(args) {
	case 0:
		opts.auto = true

	case 1:
		arg0Lower := strings.ToLower(args[0])
		switch arg0Lower {
		case "--tar":
			opts.tarPath = "backup.tar"
		case "--tgz":
			opts.tarPath = "backup.tar.gz"
		case "--vhdx":
			opts.vhdxPath = "backup.ext4.vhdx"
		case "--vhdxgz":
			opts.vhdxPath = "backup.ext4.vhdx.gz"
		case "--reg":
			opts.regPath = "backup.reg"
		default:
			if strings.HasSuffix(arg0Lower, ".tar") || strings.HasSuffix(arg0Lower, ".tar.gz") || strings.HasSuffix(arg0Lower, ".tgz") {
				opts.tarPath = args[0]
			} else if strings.HasSuffix(arg0Lower, ".ext4.vhdx") || strings.HasSuffix(arg0Lower, ".ext4.vhdx.gz") {
				opts.vhdxPath = args[0]
			} else if strings.HasSuffix(arg0Lower, ".reg") {
				opts.regPath = args[0]
			} else {
				return backupOptions{}, os.ErrInvalid
			}
		}

	default:
		return backupOptions{}, os.ErrInvalid
	}

	return opts, nil
}

func executeWithBackupsOptions(
	wsl wsllib.WslLib,
	reg wsllib.WslReg,
	name string,
	opts backupOptions,
	backupRegFn func(wsllib.WslReg, string, string) error,
	backupTarFn func(string, string) error,
	backupExt4VhdxFn func(wsllib.WslReg, string, string) error,
) error {
	opttar := opts.tarPath
	optvhdx := opts.vhdxPath
	optreg := opts.regPath

	if opts.auto {
		_, _, flags, _ := wsl.GetDistributionConfiguration(name)
		if flags&wsllib.FlagEnableWsl2 == wsllib.FlagEnableWsl2 {
			optvhdx = "backup.ext4.vhdx.gz"
			optreg = "backup.reg"
		} else {
			opttar = "backup.tar.gz"
			optreg = "backup.reg"
		}
	}

	if optreg != "" {
		err := backupRegFn(reg, name, optreg)
		if err != nil {
			return errutil.NewDisplayError(err, true, true, false)
		}
	}
	if opttar != "" {
		err := backupTarFn(name, opttar)
		if err != nil {
			return errutil.NewDisplayError(err, true, true, false)
		}
	}
	if optvhdx != "" {
		err := backupExt4VhdxFn(reg, name, optvhdx)
		if err != nil {
			return errutil.NewDisplayError(err, true, true, false)
		}
	}
	return nil
}

// executeWithBackups keeps compatibility with existing tests while using parsed options.
func executeWithBackups(
	wsl wsllib.WslLib,
	reg wsllib.WslReg,
	name string,
	args []string,
	backupRegFn func(wsllib.WslReg, string, string) error,
	backupTarFn func(string, string) error,
	backupExt4VhdxFn func(wsllib.WslReg, string, string) error,
) error {
	opts, err := parseArgs(args)
	if err != nil {
		return errutil.NewDisplayError(err, true, true, false)
	}
	return executeWithBackupsOptions(wsl, reg, name, opts, backupRegFn, backupTarFn, backupExt4VhdxFn)
}
