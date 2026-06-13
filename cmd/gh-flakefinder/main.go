package main

import (
	"os"

	appcmd "github.com/Frank-Li-Yixuan/gh-flakefinder/internal/cmd"
)

func main() {
	os.Exit(appcmd.Execute(os.Args[1:], os.Stdout, os.Stderr))
}
