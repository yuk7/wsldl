package wsllib

type MockWslLib struct {
	IsDistributionRegisteredFunc     func(name string) bool
	RegisterDistributionFunc         func(name, rootPath string) error
	UnregisterDistributionFunc       func(name string) error
	LaunchInteractiveFunc            func(name, command string, inheritPath bool) (uint32, error)
	LaunchFunc                       func(name, command string, inheritPath bool, stdin, stdout, stderr Handle) (Handle, error)
	GetDistributionConfigurationFunc func(name string) (uint32, uint64, uint32, error)
	ConfigureDistributionFunc        func(name string, uid uint64, flags uint32) error
}

func (m MockWslLib) IsDistributionRegistered(name string) bool {
	if m.IsDistributionRegisteredFunc != nil {
		return m.IsDistributionRegisteredFunc(name)
	}
	return false
}

func (m MockWslLib) RegisterDistribution(name, rootPath string) error {
	if m.RegisterDistributionFunc != nil {
		return m.RegisterDistributionFunc(name, rootPath)
	}
	return nil
}

func (m MockWslLib) UnregisterDistribution(name string) error {
	if m.UnregisterDistributionFunc != nil {
		return m.UnregisterDistributionFunc(name)
	}
	return nil
}

func (m MockWslLib) LaunchInteractive(name, command string, inheritPath bool) (uint32, error) {
	if m.LaunchInteractiveFunc != nil {
		return m.LaunchInteractiveFunc(name, command, inheritPath)
	}
	return 0, nil
}

func (m MockWslLib) Launch(name, command string, inheritPath bool, stdin, stdout, stderr Handle) (Handle, error) {
	if m.LaunchFunc != nil {
		return m.LaunchFunc(name, command, inheritPath, stdin, stdout, stderr)
	}
	return Handle(0), nil
}

func (m MockWslLib) GetDistributionConfiguration(name string) (uint32, uint64, uint32, error) {
	if m.GetDistributionConfigurationFunc != nil {
		return m.GetDistributionConfigurationFunc(name)
	}
	return 0, 0, 0, nil
}

func (m MockWslLib) ConfigureDistribution(name string, uid uint64, flags uint32) error {
	if m.ConfigureDistributionFunc != nil {
		return m.ConfigureDistributionFunc(name, uid, flags)
	}
	return nil
}

type MockWslReg struct {
	GetProfileFromNameFunc     func(name string) (Profile, error)
	GetProfileFromBasePathFunc func(path string) (Profile, error)
	WriteProfileFunc           func(profile Profile) error
	SetWslVersionFunc          func(name string, version int) error
	GenerateProfileFunc        func() Profile
}

func (m MockWslReg) GetProfileFromName(name string) (Profile, error) {
	if m.GetProfileFromNameFunc != nil {
		return m.GetProfileFromNameFunc(name)
	}
	return Profile{}, nil
}

func (m MockWslReg) GetProfileFromBasePath(path string) (Profile, error) {
	if m.GetProfileFromBasePathFunc != nil {
		return m.GetProfileFromBasePathFunc(path)
	}
	return Profile{}, nil
}

func (m MockWslReg) WriteProfile(profile Profile) error {
	if m.WriteProfileFunc != nil {
		return m.WriteProfileFunc(profile)
	}
	return nil
}

func (m MockWslReg) SetWslVersion(name string, version int) error {
	if m.SetWslVersionFunc != nil {
		return m.SetWslVersionFunc(name, version)
	}
	return nil
}

func (m MockWslReg) GenerateProfile() Profile {
	if m.GenerateProfileFunc != nil {
		return m.GenerateProfileFunc()
	}
	return Profile{}
}
