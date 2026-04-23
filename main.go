package main

import (
	_ "embed"
	"strings"

	"cli/cmd"
)

//go:embed Version
var versionFile string

func main() {
	cmd.Execute(strings.TrimSpace(versionFile))
}
