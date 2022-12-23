package main

import (
	"fmt"
	"runtime"
)

var (
	// Version carries mtr-exporter version
	Version = "0.0.0-dev-build"
	// GitHash carries git-revision if set by compile chain
	GitHash = ""
	// BuildDate carries the date of the build if set by compile chain
	BuildDate = ""
)

func printVersion() {
	fmt.Println("mtr-exporter:\t" + Version)
	fmt.Printf("compiled:\t%v on %v/%v\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)

	if GitHash != "" {
		fmt.Println("git:\t" + GitHash)
	}

	if BuildDate != "" {
		fmt.Println("build:\t" + BuildDate)
	}

	fmt.Println()
}
