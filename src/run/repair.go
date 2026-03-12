package run

import (
	"path/filepath"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func repairRegistry(reg wsllib.WslReg, profile wsllib.Profile) error {
	efPath := errutil.MustExecutable()
	dir := filepath.Dir(efPath)

	profile.BasePath = dir
	return reg.WriteProfile(profile)
}
