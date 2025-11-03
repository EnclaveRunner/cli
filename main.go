package main

import (
	"cli/cmd"
	"cli/config"
)

func main() {
	config.Version = "v0.1.3"
	cmd.Execute()
}
