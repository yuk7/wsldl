package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuk7/wsldl/clean"
	"github.com/yuk7/wsldl/config"
	"github.com/yuk7/wsldl/get"
	"github.com/yuk7/wsldl/help"
	"github.com/yuk7/wsldl/install"
	"github.com/yuk7/wsldl/isregd"
	"github.com/yuk7/wsldl/lib/wslapi"
	"github.com/yuk7/wsldl/run"
	"github.com/yuk7/wsldl/version"
)

func main() {
	efPath, _ := os.Executable()
	name := filepath.Base(efPath[:len(efPath)-len(filepath.Ext(efPath))])

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
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

		case "clean":
			clean.Execute(name, os.Args[2:])

		case "help", "-h", "--help":
			help.Execute()

		default:
			fmt.Println("Invalid Arg.")
		}
	} else {
		if !wslapi.WslIsDistributionRegistered(name) {
			install.Execute(name, nil)
		} else {
			run.Execute(name, nil)
		}
	}
}
