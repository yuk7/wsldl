package wslapi

import (
	"errors"
	"testing"

	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestGetConfig(t *testing.T) {
	t.Parallel()

	called := 0
	mock := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			called++
			if name != "Arch" {
				t.Fatalf("distribution name = %q, want %q", name, "Arch")
			}
			return 2, 1000, 6, nil
		},
	}

	uid, flags, err := GetConfig(mock, "Arch")
	if err != nil {
		t.Fatalf("GetConfig returned error: %v", err)
	}
	if uid != 1000 || flags != 6 {
		t.Fatalf("GetConfig returned uid=%d flags=%d, want uid=%d flags=%d", uid, flags, 1000, 6)
	}
	if called != 1 {
		t.Fatalf("GetDistributionConfiguration call count = %d, want 1", called)
	}
}

func TestGetConfig_Error(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	mock := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 0, 0, 0, wantErr
		},
	}

	_, _, err := GetConfig(mock, "Arch")
	if !errors.Is(err, wantErr) {
		t.Fatalf("GetConfig error = %v, want %v", err, wantErr)
	}
}
