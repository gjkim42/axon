package main

import (
	"os"

	"github.com/gjkim/axon/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
