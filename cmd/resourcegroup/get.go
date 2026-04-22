package resourcegroup

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get a resource group by name",
		Args:  cobra.ExactArgs(1),
		RunE:  runGet,
	}
}

func runGet(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.ResourceGroupColumns,
		os.Stdout,
	)

	rg, err := c.GetResourceGroup(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("get resource group: %w", err)
	}

	return printer.Print([]any{rg})
}
