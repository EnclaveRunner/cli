package main

import (
	"cli/cmd"
	"cli/config"
)

func main() {
	config.Version = "v0.2.0"
	cmd.Execute()
}
