package help

import (
	"github.com/yuk7/wsldl/backup"
	"github.com/yuk7/wsldl/clean"
	"github.com/yuk7/wsldl/config"
	"github.com/yuk7/wsldl/get"
	"github.com/yuk7/wsldl/run"
)

//Execute is default install entrypoint
func Execute(args []string) {
	if len(args) == 0 {
		ShowHelpAll()
	} else {
		switch args[0] {
		case "run", "-c", "/c", "runp", "-p", "/p":
			run.ShowHelp(true)
		case "config", "set":
			config.ShowHelp(true)
		case "get":
			get.ShowHelp(true)
		case "backup":
			backup.ShowHelp(true)
		case "clean":
			clean.ShowHelp(true)
		case "help":
			ShowHelp(true)
		default:
			ShowHelpAll()
		}
	}

}

// ShowHelpAll shows all help messages
func ShowHelpAll() {
	println("Usage :")
	run.ShowHelp(false)
	println()
	config.ShowHelp(false)
	println()
	get.ShowHelp(false)
	println()
	backup.ShowHelp(false)
	println()
	clean.ShowHelp(false)
	println()
	ShowHelp(false)
}

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    help")
	println("      - Print this usage message.")
}
