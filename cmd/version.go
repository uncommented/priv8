package main

import (
	"fmt"
	"runtime"
)

var (
	Executable = "<unset>"
	VCSCommit  = "<unset>"
	VCSTag     = "<unset>"
)

func printVersion() {
	fmt.Println("priv8 version:")
	fmt.Printf("  Executable: %s\n", Executable)
	fmt.Printf("  VCS Commit: %s\n", VCSCommit)
	fmt.Printf("  VCS Tag: %s\n", VCSTag)
	fmt.Printf("  Go Version: %s (%s/%s) \n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
