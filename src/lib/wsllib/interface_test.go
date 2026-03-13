package wsllib

import "testing"

type testWslLib struct{}

func (testWslLib) IsDistributionRegistered(string) bool { return false }
func (testWslLib) RegisterDistribution(string, string) error {
	return nil
}
func (testWslLib) UnregisterDistribution(string) error { return nil }
func (testWslLib) LaunchInteractive(string, string, bool) (uint32, error) {
	return 0, nil
}
func (testWslLib) Launch(string, string, bool, Handle, Handle, Handle) (Handle, error) {
	return 0, nil
}
func (testWslLib) GetDistributionConfiguration(string) (uint32, uint64, uint32, error) {
	return 0, 0, 0, nil
}
func (testWslLib) ConfigureDistribution(string, uint64, uint32) error { return nil }

type testWslReg struct{}

func (testWslReg) GetProfileFromName(string) (Profile, error)     { return Profile{}, nil }
func (testWslReg) GetProfileFromBasePath(string) (Profile, error) { return Profile{}, nil }
func (testWslReg) WriteProfile(Profile) error                     { return nil }
func (testWslReg) SetWslVersion(string, int) error                { return nil }
func (testWslReg) GenerateProfile() Profile                       { return Profile{} }

func TestNewDependenciesFrom_UnitTestBranch(t *testing.T) {
	t.Parallel()

	deps := newDependencies(
		true,
		func() WslLib {
			t.Fatal("newWsl should not be called in unit test branch")
			return nil
		},
		func() WslReg {
			t.Fatal("newReg should not be called in unit test branch")
			return nil
		},
	)

	if _, ok := deps.Wsl.(MockWslLib); !ok {
		t.Fatalf("deps.Wsl type = %T, want MockWslLib", deps.Wsl)
	}
	if _, ok := deps.Reg.(MockWslReg); !ok {
		t.Fatalf("deps.Reg type = %T, want MockWslReg", deps.Reg)
	}
}

func TestNewDependenciesFrom_NativeBranch(t *testing.T) {
	t.Parallel()

	wslCalls := 0
	regCalls := 0

	deps := newDependencies(
		false,
		func() WslLib {
			wslCalls++
			return testWslLib{}
		},
		func() WslReg {
			regCalls++
			return testWslReg{}
		},
	)

	if wslCalls != 1 {
		t.Fatalf("newWsl call count = %d, want 1", wslCalls)
	}
	if regCalls != 1 {
		t.Fatalf("newReg call count = %d, want 1", regCalls)
	}
	if _, ok := deps.Wsl.(testWslLib); !ok {
		t.Fatalf("deps.Wsl type = %T, want testWslLib", deps.Wsl)
	}
	if _, ok := deps.Reg.(testWslReg); !ok {
		t.Fatalf("deps.Reg type = %T, want testWslReg", deps.Reg)
	}
}

func TestNewDependenciesForProcess_UnitTestBranch(t *testing.T) {
	t.Parallel()

	deps := newDependenciesForProcess(true)
	if _, ok := deps.Wsl.(MockWslLib); !ok {
		t.Fatalf("deps.Wsl type = %T, want MockWslLib", deps.Wsl)
	}
	if _, ok := deps.Reg.(MockWslReg); !ok {
		t.Fatalf("deps.Reg type = %T, want MockWslReg", deps.Reg)
	}
}

func TestNewDependenciesForProcess_NativeBranch(t *testing.T) {
	t.Parallel()

	deps := newDependenciesForProcess(false)
	if deps.Wsl == nil {
		t.Fatal("deps.Wsl is nil")
	}
	if deps.Reg == nil {
		t.Fatal("deps.Reg is nil")
	}
}
