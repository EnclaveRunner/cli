package me

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete the currently authenticated user",
		RunE:  runDelete,
	}
}

func runDelete(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.UserColumns,
		os.Stdout,
	)

	u, err := c.DeleteMe(cmd.Context())
	if err != nil {
		return fmt.Errorf("delete me: %w", err)
	}

	return printer.Print([]any{u})
}
