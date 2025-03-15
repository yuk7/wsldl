package cmdline

type Command struct {
	Names []string
	Run   func(distroName string, args []string)
}
