//go:build !windows

package wsllib

import (
	"errors"
	"testing"
)

func TestNewNativeWslLibStub(t *testing.T) {
	t.Parallel()

	wsl := NewNativeWslLib()

	if wsl.IsDistributionRegistered("Arch") {
		t.Fatal("IsDistributionRegistered = true, want false")
	}

	if err := wsl.RegisterDistribution("Arch", "rootfs.tar"); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("RegisterDistribution err = %v, want %v", err, errUnsupportedPlatform)
	}
	if err := wsl.UnregisterDistribution("Arch"); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("UnregisterDistribution err = %v, want %v", err, errUnsupportedPlatform)
	}
	if _, err := wsl.LaunchInteractive("Arch", "echo hi", true); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("LaunchInteractive err = %v, want %v", err, errUnsupportedPlatform)
	}
	if _, err := wsl.Launch("Arch", "echo hi", true, 0, 0, 0); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("Launch err = %v, want %v", err, errUnsupportedPlatform)
	}
	if _, _, _, err := wsl.GetDistributionConfiguration("Arch"); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("GetDistributionConfiguration err = %v, want %v", err, errUnsupportedPlatform)
	}
	if err := wsl.ConfigureDistribution("Arch", 1000, 0); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("ConfigureDistribution err = %v, want %v", err, errUnsupportedPlatform)
	}
}

func TestNewNativeWslRegStub(t *testing.T) {
	t.Parallel()

	reg := NewNativeWslReg()

	if _, err := reg.GetProfileFromName("Arch"); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("GetProfileFromName err = %v, want %v", err, errUnsupportedPlatform)
	}
	if _, err := reg.GetProfileFromBasePath("C:\\WSL"); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("GetProfileFromBasePath err = %v, want %v", err, errUnsupportedPlatform)
	}
	if err := reg.WriteProfile(Profile{DistributionName: "Arch"}); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("WriteProfile err = %v, want %v", err, errUnsupportedPlatform)
	}
	if err := reg.SetWslVersion("Arch", 2); !errors.Is(err, errUnsupportedPlatform) {
		t.Fatalf("SetWslVersion err = %v, want %v", err, errUnsupportedPlatform)
	}
	if got := reg.GenerateProfile(); got != (Profile{}) {
		t.Fatalf("GenerateProfile = %+v, want zero value", got)
	}
}

func TestNewDependencies_UsesMocksInUnitTests(t *testing.T) {
	t.Parallel()

	deps := NewDependencies()

	if _, ok := deps.Wsl.(MockWslLib); !ok {
		t.Fatalf("deps.Wsl type = %T, want MockWslLib", deps.Wsl)
	}
	if _, ok := deps.Reg.(MockWslReg); !ok {
		t.Fatalf("deps.Reg type = %T, want MockWslReg", deps.Reg)
	}
}
