package get

import (
	"errors"
	"fmt"
	"log"
	"syscall"

	"github.com/yuk7/wsldl/lib/wslapi"
)

//Execute is default install entrypoint
func Execute(name string, args []string) {
	uid, flags := WslGetConfig(name)
	if len(args) == 1 {
		switch args[0] {
		case "--default-uid":
			print(uid)

		case "--append-path":
			print(flags&wslapi.FlagAppendNTPath == wslapi.FlagAppendNTPath)

		case "--mount-drive":
			print(flags&wslapi.FlagEnableDriveMounting == wslapi.FlagEnableDriveMounting)

		case "--wsl-version":
			if flags&wslapi.FlagEnableWsl2 == wslapi.FlagEnableWsl2 {
				print("2")
			} else {
				print("1")
			}

		case "--flags-val":
			print(flags)

		case "--flags-bits":
			fmt.Printf("%04b", flags)

		default:
			println("ERR: Invalid argument")
			err := errors.New("invalid args")
			log.Fatal(err)
		}
	} else {
		println("ERR: Invalid argument")
		err := errors.New("invalid args")
		log.Fatal(err)
	}
}

//WslGetConfig is getter of distribution configuration
func WslGetConfig(distributionName string) (uid uint64, flags uint32) {
	var err error
	_, uid, flags, err = wslapi.WslGetDistributionConfiguration(distributionName)
	if err != nil {
		fmt.Println("ERR: Failed to GetDistributionConfiguration")
		var errno syscall.Errno
		if errors.As(err, &errno) {
			fmt.Printf("Code: 0x%x\n", int(errno))
		}
		log.Fatal(err)
	}
	return
}
