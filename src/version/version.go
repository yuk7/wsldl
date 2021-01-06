package version

import (
	"fmt"
)

var (
	project string = "wsldl2"
	version string = "Unknown"
	url     string = "https://git.io/wsldl"
)

// Print version infomation
func Print() {
	fmt.Printf("%s, version %s\n", project, version)
	fmt.Printf("%s\n", url)
}
