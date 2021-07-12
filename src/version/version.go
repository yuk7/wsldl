package version

import (
	"fmt"
	"runtime"
)

var (
	project string = "wsldl2"
	version string = "Unknown"
	url     string = "https://git.io/wsldl"
)

//Execute is default version entrypoint. prints version infomation
func Execute() {
	fmt.Printf("%s, version %s  (%s)\n", project, version, runtime.GOARCH)
	fmt.Printf("%s\n", url)
}
