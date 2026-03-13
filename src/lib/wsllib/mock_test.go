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

func TestMockWslLib_Defaults_AllMethods(t *testing.T) {
	t.Parallel()

	m := MockWslLib{}

	if got := m.IsDistributionRegistered("x"); got {
		t.Fatal("IsDistributionRegistered default = true, want false")
	}
	if err := m.RegisterDistribution("x", "rootfs.tar"); err != nil {
		t.Fatalf("RegisterDistribution default error = %v, want nil", err)
	}
	if err := m.UnregisterDistribution("x"); err != nil {
		t.Fatalf("UnregisterDistribution default error = %v, want nil", err)
	}
	if code, err := m.LaunchInteractive("x", "echo hi", true); err != nil || code != 0 {
		t.Fatalf("LaunchInteractive default = (%d, %v), want (0, nil)", code, err)
	}
	if handle, err := m.Launch("x", "echo hi", true, 1, 2, 3); err != nil || handle != 0 {
		t.Fatalf("Launch default = (%d, %v), want (0, nil)", handle, err)
	}
	if v, uid, flags, err := m.GetDistributionConfiguration("x"); err != nil || v != 0 || uid != 0 || flags != 0 {
		t.Fatalf("GetDistributionConfiguration default = (%d,%d,%d,%v), want (0,0,0,nil)", v, uid, flags, err)
	}
	if err := m.ConfigureDistribution("x", 1000, 6); err != nil {
		t.Fatalf("ConfigureDistribution default error = %v, want nil", err)
	}
}

func TestMockWslLib_Callbacks_AllMethods(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	called := 0
	m := MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return true
		},
		RegisterDistributionFunc: func(name, rootPath string) error {
			called++
			if name != "Arch" || rootPath != "rootfs.tar" {
				t.Fatalf("RegisterDistribution args mismatch: name=%q rootPath=%q", name, rootPath)
			}
			return wantErr
		},
		UnregisterDistributionFunc: func(name string) error {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return wantErr
		},
		LaunchInteractiveFunc: func(name, command string, inheritPath bool) (uint32, error) {
			called++
			if name != "Arch" || command != "echo hi" || !inheritPath {
				t.Fatalf("LaunchInteractive args mismatch: name=%q command=%q inheritPath=%v", name, command, inheritPath)
			}
			return 7, wantErr
		},
		LaunchFunc: func(name, command string, inheritPath bool, stdin, stdout, stderr Handle) (Handle, error) {
			called++
			if name != "Arch" || command != "echo hi" || inheritPath || stdin != 1 || stdout != 2 || stderr != 3 {
				t.Fatalf("Launch args mismatch: name=%q command=%q inheritPath=%v stdin=%d stdout=%d stderr=%d", name, command, inheritPath, stdin, stdout, stderr)
			}
			return Handle(9), wantErr
		},
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return 2, 1000, 6, wantErr
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
	if err := m.RegisterDistribution("Arch", "rootfs.tar"); !errors.Is(err, wantErr) {
		t.Fatalf("RegisterDistribution error = %v, want %v", err, wantErr)
	}
	if err := m.UnregisterDistribution("Arch"); !errors.Is(err, wantErr) {
		t.Fatalf("UnregisterDistribution error = %v, want %v", err, wantErr)
	}
	if code, err := m.LaunchInteractive("Arch", "echo hi", true); !errors.Is(err, wantErr) || code != 7 {
		t.Fatalf("LaunchInteractive = (%d, %v), want (7, %v)", code, err, wantErr)
	}
	if handle, err := m.Launch("Arch", "echo hi", false, 1, 2, 3); !errors.Is(err, wantErr) || handle != Handle(9) {
		t.Fatalf("Launch = (%d, %v), want (9, %v)", handle, err, wantErr)
	}
	v, uid, flags, err := m.GetDistributionConfiguration("Arch")
	if !errors.Is(err, wantErr) || v != 2 || uid != 1000 || flags != 6 {
		t.Fatalf("GetDistributionConfiguration = (%d,%d,%d,%v), want (2,1000,6,%v)", v, uid, flags, err, wantErr)
	}
	if err := m.ConfigureDistribution("Arch", 1000, 6); !errors.Is(err, wantErr) {
		t.Fatalf("ConfigureDistribution error = %v, want %v", err, wantErr)
	}

	if called != 7 {
		t.Fatalf("callback call count = %d, want 7", called)
	}
}

func TestMockWslReg_Defaults_AllMethods(t *testing.T) {
	t.Parallel()

	m := MockWslReg{}
	if p, err := m.GetProfileFromName("Arch"); err != nil || p != (Profile{}) {
		t.Fatalf("GetProfileFromName default = (%+v, %v), want (zero, nil)", p, err)
	}
	if p, err := m.GetProfileFromBasePath("X:\\WSL"); err != nil || p != (Profile{}) {
		t.Fatalf("GetProfileFromBasePath default = (%+v, %v), want (zero, nil)", p, err)
	}
	if err := m.WriteProfile(Profile{DistributionName: "Arch"}); err != nil {
		t.Fatalf("WriteProfile default error = %v, want nil", err)
	}
	if err := m.SetWslVersion("Arch", 2); err != nil {
		t.Fatalf("SetWslVersion default error = %v, want nil", err)
	}
	if p := m.GenerateProfile(); p != (Profile{}) {
		t.Fatalf("GenerateProfile default = %+v, want zero value", p)
	}
}

func TestMockWslReg_Callbacks_AllMethods(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	called := 0
	fromNameProfile := Profile{DistributionName: "Arch", BasePath: "X:\\one"}
	fromBasePathProfile := Profile{DistributionName: "Debian", BasePath: "X:\\two"}
	generated := Profile{Flags: FlagEnableWsl2, WsldlTerm: FlagWsldlTermWT}

	m := MockWslReg{
		GetProfileFromNameFunc: func(name string) (Profile, error) {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return fromNameProfile, wantErr
		},
		GetProfileFromBasePathFunc: func(path string) (Profile, error) {
			called++
			if path != "X:\\base" {
				t.Fatalf("path = %q, want %q", path, "X:\\base")
			}
			return fromBasePathProfile, wantErr
		},
		WriteProfileFunc: func(profile Profile) error {
			called++
			if profile.DistributionName != "Arch" {
				t.Fatalf("profile.DistributionName = %q, want %q", profile.DistributionName, "Arch")
			}
			return wantErr
		},
		SetWslVersionFunc: func(name string, version int) error {
			called++
			if name != "Arch" || version != 2 {
				t.Fatalf("SetWslVersion args mismatch: name=%q version=%d", name, version)
			}
			return wantErr
		},
		GenerateProfileFunc: func() Profile {
			called++
			return generated
		},
	}

	if p, err := m.GetProfileFromName("Arch"); p != fromNameProfile || !errors.Is(err, wantErr) {
		t.Fatalf("GetProfileFromName = (%+v, %v), want (%+v, %v)", p, err, fromNameProfile, wantErr)
	}
	if p, err := m.GetProfileFromBasePath("X:\\base"); p != fromBasePathProfile || !errors.Is(err, wantErr) {
		t.Fatalf("GetProfileFromBasePath = (%+v, %v), want (%+v, %v)", p, err, fromBasePathProfile, wantErr)
	}
	if err := m.WriteProfile(Profile{DistributionName: "Arch"}); !errors.Is(err, wantErr) {
		t.Fatalf("WriteProfile error = %v, want %v", err, wantErr)
	}
	if err := m.SetWslVersion("Arch", 2); !errors.Is(err, wantErr) {
		t.Fatalf("SetWslVersion error = %v, want %v", err, wantErr)
	}
	if p := m.GenerateProfile(); p != generated {
		t.Fatalf("GenerateProfile = %+v, want %+v", p, generated)
	}
	if called != 5 {
		t.Fatalf("callback call count = %d, want 5", called)
	}
}
