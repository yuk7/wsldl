package wtutils

// Config is json root Config
type Config struct {
	Profiles struct {
		ProfileList []Profile `json:"list"`
	} `json:"profiles"`
}

// Profile is profile of terminal
type Profile struct {
	Name        string `json:"name"`
	CommandLine string `json:"commandline"`
	GUID        string `json:"guid"`
	Source      string `json:"source"`
}
