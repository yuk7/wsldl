package preset

type Preset struct {
	WslVersion  int    `json:"wslversion,omitempty"`
	InstallFile string `json:"installfile,omitempty"`
}
