package main

import (
	"flag"
	"os"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")

	flag.Parse()

	if *versionFlag {
		printVersion()
		os.Exit(0)
	}
}
