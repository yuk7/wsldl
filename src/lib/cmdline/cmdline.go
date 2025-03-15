package cmdline

// RunSubCommand executes a subcommand
func RunSubCommand(commands []Command, mismatch func(), distroName string, args []string) {
	if len(args) > 0 {
		for _, c := range commands {
			for _, n := range c.Names {
				if n == args[0] {
					c.Run(distroName, args[1:])
					return
				}
			}
		}
	}
	mismatch()
}
