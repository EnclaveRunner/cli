package resourcegroup

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new resource group",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreate,
	}
	cmd.Flags().
		StringSlice("endpoints", nil, "API endpoints to include in the resource group")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.ResourceGroupColumns,
		os.Stdout,
	)

	endpoints, _ := cmd.Flags().GetStringSlice("endpoints")
	rg, err := c.CreateResourceGroup(cmd.Context(), args[0], endpoints)
	if err != nil {
		return fmt.Errorf("create resource group: %w", err)
	}

	return printer.Print([]any{rg})
}
