package cmdline

type Command struct {
	Names     []string
	IsDefault bool
	Visible   func(distroName string) bool
	HelpText  func() string
	Run       func(distroName string, args []string) error
}
