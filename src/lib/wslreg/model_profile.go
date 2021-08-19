package wslreg

// Profile is profile for WSL
type Profile struct {
	UUID              string
	BasePath          string
	DistributionName  string
	DefaultUid        int
	Flags             int
	State             int
	Version           int
	PackageFamilyName string
	WsldlTerm         int
}
