package backup

import (
	"os"

	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
)

//Execute is default run entrypoint.
func Execute(name string, args []string) {
	opttar := false
	opttgz := false
	optvhdx := false
	optvhdxgz := false
	optreg := false
	switch len(args) {
	case 0:
		_, _, flags, _ := wsllib.WslGetDistributionConfiguration(name)
		if flags&wsllib.FlagEnableWsl2 == wsllib.FlagEnableWsl2 {
			optvhdxgz = true
			optreg = true
		} else {
			opttgz = true
			optreg = true
		}

	case 1:
		switch args[0] {
		case "--tar":
			opttar = true
		case "--tgz":
			opttgz = true
		case "--vhdx":
			optvhdx = true
		case "--vhdxgz":
			optvhdxgz = true
		case "--reg":
			optreg = true
		}

	default:
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}

	if optreg {
		err := backupReg(name, "backup.reg")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}
	if opttar {
		err := backupTar(name, "backup.tar")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}

	}
	if opttgz {
		err := backupTar(name, "backup.tar.gz")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}
	if optvhdx {
		err := backupExt4Vhdx(name, "backup.ext4.vhdx")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}
	if optvhdxgz {
		err := backupExt4Vhdx(name, "backup.ext4.vhdx.gz")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}
}
