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
	"github.com/yuk7/wsldl/lib/cmdline"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/run"
	"github.com/yuk7/wsldl/version"
	"github.com/yuk7/wsllib-go"
)

func main() {
	efPath, _ := os.Executable()
	name := filepath.Base(efPath[:len(efPath)-len(filepath.Ext(efPath))])

	var commands = []cmdline.Command{
		isregd.GetCommand(),
		version.GetCommand(),
		install.GetCommand(),
		run.GetCommandWithNoArgs(),
		run.GetCommand(),
		run.GetCommandP(),
		config.GetCommand(),
		get.GetCommand(),
		backup.GetCommand(),
		clean.GetCommand(),
		help.GetCommand(),
	}

	if len(os.Args) > 1 {
		cmdline.RunSubCommand(
			commands,
			func() {
				utils.ErrorExit(os.ErrInvalid, true, true, false)
			},
			name,
			os.Args[1:],
		)

	} else {
		if !wsllib.WslIsDistributionRegistered(name) {
			install.GetCommand().Run(name, nil)
		} else {
			run.GetCommandWithNoArgs().Run(name, nil)
		}
	}
}
