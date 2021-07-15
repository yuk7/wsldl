package main

import (
	"os"
	"path/filepath"

	"github.com/yuk7/wsldl/backup"
	"github.com/yuk7/wsldl/clean"
	"github.com/yuk7/wsldl/config"
	"github.com/yuk7/wsldl/get"
	"github.com/yuk7/wsldl/help"
	"github.com/yuk7/wsldl/install"
	"github.com/yuk7/wsldl/isregd"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wslapi"
	"github.com/yuk7/wsldl/run"
	"github.com/yuk7/wsldl/version"
)

func main() {
	efPath, _ := os.Executable()
	name := filepath.Base(efPath[:len(efPath)-len(filepath.Ext(efPath))])

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "-v", "--version":
			version.Execute()

		case "isregd":
			isregd.Execute(name)

		case "install":
			install.Execute(name, os.Args[2:])

		case "run", "-c", "/c":
			run.Execute(name, os.Args[2:])

		case "runp", "-p", "/p":
			run.ExecuteP(name, os.Args[2:])

		case "config", "set":
			config.Execute(name, os.Args[2:])

		case "get":
			get.Execute(name, os.Args[2:])

		case "backup":
			backup.Execute(name, os.Args[2:])

		case "clean":
			clean.Execute(name, os.Args[2:])

		case "help", "-h", "--help", "/?":
			help.Execute(os.Args[2:])

		default:
			utils.ErrorExit(os.ErrInvalid, true, true, false)
		}
	} else {
		if !wslapi.WslIsDistributionRegistered(name) {
			install.Execute(name, nil)
		} else {
			run.ExecuteNoArgs(name)
		}
	}
}
