package main

import (
	"fmt"
	"os"

	"github.com/yuk7/wsldl/version"
)

func main() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			version.Print()

		default:
			fmt.Println("Invalid Arg.")
		}
	}
}
