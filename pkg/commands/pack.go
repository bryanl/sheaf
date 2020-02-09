package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewPackCommand creates a pack command.
func NewPackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pack",
		Short: "pack a bundle",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires path to bundle directory")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config := sheaf.PackConfig{Path: args[0]}
			return sheaf.Pack(config)
		},
	}

	return cmd
}
