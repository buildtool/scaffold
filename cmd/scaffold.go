package main

import (
	service "github.com/buildtool/scaffold/pkg"
	ver "github.com/buildtool/scaffold/pkg/version"
	"io"
	"os"
)

var (
	version            = "dev"
	commit             = "none"
	date               = "unknown"
	exitFunc           = os.Exit
	out      io.Writer = os.Stdout
)

func main() {
	if ver.PrintVersionOnly(version, commit, date, out) {
		exitFunc(0)
	} else {
		dir, _ := os.Getwd()
		exitFunc(service.Setup(dir, os.Stdout, os.Args[1:]...))
	}
}
