package run

import (
	"os"
	"path/filepath"

	"github.com/yuk7/wsldl/lib/wsllib"
)

func repairRegistry(reg wsllib.WslReg, profile wsllib.Profile) error {
	efPath, _ := os.Executable()
	dir := filepath.Dir(efPath)

	profile.BasePath = dir
	return reg.WriteProfile(profile)
}
