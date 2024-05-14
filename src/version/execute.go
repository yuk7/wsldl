package version

import (
	"fmt"
	"runtime"
)

//Execute is default version entrypoint. prints version information
func Execute() {
	fmt.Printf("%s, version %s  (%s)\n", project, version, runtime.GOARCH)
	fmt.Printf("%s\n", url)
}
