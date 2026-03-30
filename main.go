package main

import "github.com/mparvin/tfd/cmd"

// version is set via ldflags during build
var version string

func main() {
	cmd.Execute()
}
