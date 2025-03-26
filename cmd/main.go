package main

import (
	"flag"
	"os"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")

	if *versionFlag {
		printVersion()
		os.Exit(0)
	}
}
