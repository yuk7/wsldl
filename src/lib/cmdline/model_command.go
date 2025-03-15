package cmdline

type Command struct {
	Names []string
	Help  func(distroName string, isListQuery bool) string
	Run   func(distroName string, args []string)
}
