package wsllib

import (
	"errors"
	"testing"
)

func TestMockWslLib_Defaults(t *testing.T) {
	t.Parallel()

	m := MockWslLib{}
	if m.IsDistributionRegistered("x") {
		t.Fatal("IsDistributionRegistered should default to false")
	}
	if err := m.RegisterDistribution("x", "y"); err != nil {
		t.Fatalf("RegisterDistribution default error = %v, want nil", err)
	}
	if err := m.UnregisterDistribution("x"); err != nil {
		t.Fatalf("UnregisterDistribution default error = %v, want nil", err)
	}
	if _, _, _, err := m.GetDistributionConfiguration("x"); err != nil {
		t.Fatalf("GetDistributionConfiguration default error = %v, want nil", err)
	}
}

func TestMockWslLib_Callbacks(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	called := 0
	m := MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			called++
			return name == "Arch"
		},
		ConfigureDistributionFunc: func(name string, uid uint64, flags uint32) error {
			called++
			if name != "Arch" || uid != 1000 || flags != 6 {
				t.Fatalf("ConfigureDistribution args mismatch: name=%q uid=%d flags=%d", name, uid, flags)
			}
			return wantErr
		},
	}

	if !m.IsDistributionRegistered("Arch") {
		t.Fatal("IsDistributionRegistered callback result mismatch")
	}
	err := m.ConfigureDistribution("Arch", 1000, 6)
	if !errors.Is(err, wantErr) {
		t.Fatalf("ConfigureDistribution error = %v, want %v", err, wantErr)
	}
	if called != 2 {
		t.Fatalf("callback call count = %d, want 2", called)
	}
}

func TestMockWslReg_DefaultsAndCallbacks(t *testing.T) {
	t.Parallel()

	wantProfile := Profile{DistributionName: "Arch", BasePath: "C:\\WSL"}
	called := 0
	m := MockWslReg{
		GetProfileFromNameFunc: func(name string) (Profile, error) {
			called++
			return wantProfile, nil
		},
		GenerateProfileFunc: func() Profile {
			called++
			return Profile{Flags: FlagEnableWsl2}
		},
	}

	got, err := m.GetProfileFromName("Arch")
	if err != nil {
		t.Fatalf("GetProfileFromName failed: %v", err)
	}
	if got != wantProfile {
		t.Fatalf("profile = %+v, want %+v", got, wantProfile)
	}
	generated := m.GenerateProfile()
	if generated.Flags != FlagEnableWsl2 {
		t.Fatalf("generated.Flags = %d, want %d", generated.Flags, FlagEnableWsl2)
	}
	if called != 2 {
		t.Fatalf("callback call count = %d, want 2", called)
	}
}
