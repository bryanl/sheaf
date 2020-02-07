package main

import (
	"os"

	"github.com/bryanl/sheaf/pkg/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
