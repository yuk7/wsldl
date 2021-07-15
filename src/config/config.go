package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/yuk7/wsldl/get"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wslapi"
	"github.com/yuk7/wsldl/run"
)

//Execute is default install entrypoint
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
				flags |= wslapi.FlagAppendNTPath
			} else {
				flags ^= wslapi.FlagAppendNTPath
			}

		case "--mount-drive":
			var b bool
			b, err = strconv.ParseBool(args[1])
			if b {
				flags |= wslapi.FlagEnableDriveMounting
			} else {
				flags ^= wslapi.FlagEnableDriveMounting
			}

		case "--default-term":
			value := 0
			switch args[1] {
			case "default", strconv.Itoa(utils.FlagWsldlTermDefault):
				value = utils.FlagWsldlTermDefault
			case "wt", strconv.Itoa(utils.FlagWsldlTermWT):
				value = utils.FlagWsldlTermWT
			case "flute", strconv.Itoa(utils.FlagWsldlTermFlute):
				value = utils.FlagWsldlTermFlute
			default:
				err = os.ErrInvalid
				break
			}
			uuid, err := utils.WslGetUUID(name)
			if err != nil {
				break
			}
			err = utils.WsldlSetTerminalInfo(uuid, value)

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
		wslapi.WslConfigureDistribution(name, uid, flags)
	} else {
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}
}

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    config [setting [value]]")
	println("      - `--default-user <user>`: Set the default user of this instance to <user>.")
	println("      - `--default-uid <uid>`: Set the default user uid of this instance to <uid>.")
	println("      - `--append-path <true|false>`: Switch of Append Windows PATH to $PATH")
	println("      - `--mount-drive <true|false>`: Switch of Mount drives")
	println("      - `--default-term <default|wt|flute>`: Set default type of terminal window.")
}
