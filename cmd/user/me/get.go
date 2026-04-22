package me

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get the currently authenticated user",
		RunE:  runGet,
	}
}

func runGet(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.UserColumns,
		os.Stdout,
	)

	u, err := c.GetMe(cmd.Context())
	if err != nil {
		return fmt.Errorf("get me: %w", err)
	}

	return printer.Print([]any{u})
}
