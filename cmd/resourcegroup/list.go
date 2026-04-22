package resourcegroup

import (
	"fmt"
	"os"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all resource groups",
		RunE:  runList,
	}
}

func runList(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.ResourceGroupColumns, os.Stdout)

	rgs, err := enclave.Collect(c.ListResourceGroups(cmd.Context()))
	if err != nil {
		return fmt.Errorf("list resource groups: %w", err)
	}
	return printer.Print(rgs)
}
