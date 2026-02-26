package run

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestRepairRegistry_UpdatesBasePathAndWritesProfile(t *testing.T) {
	t.Parallel()

	written := wsllib.Profile{}
	called := 0
	reg := wsllib.MockWslReg{
		WriteProfileFunc: func(profile wsllib.Profile) error {
			called++
			written = profile
			return nil
		},
	}

	input := wsllib.Profile{
		BasePath:         "C:\\old-path",
		DistributionName: "Arch",
		DefaultUid:       1000,
	}

	if err := repairRegistry(reg, input); err != nil {
		t.Fatalf("repairRegistry returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("WriteProfile call count = %d, want 1", called)
	}

	efPath, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable failed: %v", err)
	}
	wantBasePath := filepath.Dir(efPath)
	if written.BasePath != wantBasePath {
		t.Fatalf("BasePath = %q, want %q", written.BasePath, wantBasePath)
	}
	if written.DistributionName != input.DistributionName {
		t.Fatalf("DistributionName = %q, want %q", written.DistributionName, input.DistributionName)
	}
	if written.DefaultUid != input.DefaultUid {
		t.Fatalf("DefaultUid = %d, want %d", written.DefaultUid, input.DefaultUid)
	}
}

func TestRepairRegistry_WriteProfileError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("write failed")
	reg := wsllib.MockWslReg{
		WriteProfileFunc: func(profile wsllib.Profile) error {
			return wantErr
		},
	}

	err := repairRegistry(reg, wsllib.Profile{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
