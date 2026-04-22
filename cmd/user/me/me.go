package me

import (
	"github.com/spf13/cobra"
)

// NewCmd returns the "user me" command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Manage the currently authenticated user",
	}
	cmd.AddCommand(
		newGetCmd(),
		newUpdateCmd(),
		newDeleteCmd(),
	)

	return cmd
}
