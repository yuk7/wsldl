package config

import (
	"errors"
	"log"
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
				err = errors.New("invalid args")
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
			err = errors.New("invalid args")
		}
		if err != nil {
			println("ERR: Failed to parse your argument")
			log.Fatal(err)
		}
		wslapi.WslConfigureDistribution(name, uid, flags)
	} else {
		println("ERR: Invalid argument")
		err = errors.New("invalid args")
		log.Fatal(err)
	}
}

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    config [setting [value]]")
	println("      - `--default-user <user>`: Set the default user for this distro to <user>")
	println("      - `--default-uid <uid>`: Set the default user for this distro to <uid>")
	println("      - `--append-path <true|false>`: Switch of Append Windows PATH to $PATH")
	println("      - `--mount-drive <on|off>`: Switch of Mount drives")
	println("      - `--default-term <default|wt|flute>`: Set default terminal window")
}
