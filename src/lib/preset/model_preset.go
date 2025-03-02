package preset

type Preset struct {
	WslVersion        int    `json:"wslversion,omitempty"`
	InstallFile       string `json:"installfile,omitempty"`
	InstallFileSha256 string `json:"installfilesha256,omitempty"`
}
