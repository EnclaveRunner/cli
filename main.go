package main

import (
	"cli/cmd"
	_ "embed"
	"strings"
)

//go:embed Version
var versionFile string

func main() {
	cmd.Execute(strings.TrimSpace(versionFile))
}
