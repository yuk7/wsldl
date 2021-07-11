package help

import (
	"github.com/yuk7/wsldl/config"
	"github.com/yuk7/wsldl/get"
	"github.com/yuk7/wsldl/run"
)

//Execute is default install entrypoint
func Execute() {
	println("Usage :")
	run.ShowHelp(false)
	println()
	config.ShowHelp(false)
	println()
	get.ShowHelp(false)
	println()
	ShowHelp(false)
}

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    clean")
	println("      - Uninstall the distro.")
}
