package backup

import (
	"os"
	"os/exec"

	"github.com/yuk7/wsldl/lib/utils"
)

//Execute is default run entrypoint.
func Execute(name string, args []string) {
	opttar := false
	optreg := true
	switch len(args) {
	case 0:
		opttar = true
		optreg = true

	case 1:
		switch args[0] {
		case "--tar":
			opttar = true
		case "--reg":
			optreg = true
		}

	default:
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}

	if optreg {
		lxguid, err := utils.WslGetUUID(name)
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
		regexe := os.Getenv("SystemRoot") + "\\System32\\reg.exe"
		regpath := "HKEY_CURRENT_USER\\" + utils.LxssBaseKey + "\\" + lxguid
		_, err = exec.Command(regexe, "export", regpath, "backup.reg", "/y").Output()
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}
	if opttar {
		wslexe := os.Getenv("SystemRoot") + "\\System32\\wsl.exe"
		_, err := exec.Command(wslexe, "--export", name, "backup.tar").Output()
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}
}

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    backup [contents]")
	println("      - `--tar`: Output backup.tar to the current directory.")
	println("      - `--reg`: Output settings registry file to the current directory.")
}
