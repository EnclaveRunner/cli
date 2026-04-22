package user

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
		Short: "List all users",
		RunE:  runList,
	}
}

func runList(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.UserColumns,
		os.Stdout,
	)

	users, err := enclave.Collect(c.ListUsers(cmd.Context()))
	if err != nil {
		return fmt.Errorf("list users: %w", err)
	}

	return printer.Print(users)
}
