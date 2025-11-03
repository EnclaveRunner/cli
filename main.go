package main

import (
	"cli/cmd"
	"cli/config"
)

func main() {
	config.Version = "v0.1.4"
	cmd.Execute()
}
