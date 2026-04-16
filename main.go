package main

import (
	"os"

	"github.com/driftctl/driftctl-diff/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
