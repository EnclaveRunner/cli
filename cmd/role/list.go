package role

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all roles",
		RunE:  runList,
	}
}

func runList(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.RoleColumns,
		os.Stdout,
	)

	roles, err := enclave.Collect(c.ListRoles(cmd.Context()))
	if err != nil {
		return fmt.Errorf("list roles: %w", err)
	}

	return printer.Print(roles)
}
