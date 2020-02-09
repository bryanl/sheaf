package commands

import "github.com/spf13/cobra"

// NewRootCommand creates a root command.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sheaf",
		Short:        "sheaf bundles Kubernetes applications",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		NewCreateCommand(),
		NewStageCommand())

	return cmd
}

// Execute executes the root command for sheaf.
func Execute() error {
	root := NewRootCommand()
	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}
