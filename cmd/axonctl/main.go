package main

import (
	"os"

	"github.com/gjkim/axon/cmd/axonctl/commands"
)

func main() {
	if err := commands.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
