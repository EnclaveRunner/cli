package user

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <username>",
		Short: "Update a user",
		Args:  cobra.ExactArgs(1),
		RunE:  runUpdate,
	}
	cmd.Flags().String("display-name", "", "New display name")
	cmd.Flags().String("password", "", "New password")

	return cmd
}

func runUpdate(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.UserColumns,
		os.Stdout,
	)

	var opts []enclave.UpdateUserOption
	if v, _ := cmd.Flags().GetString("display-name"); v != "" {
		opts = append(opts, enclave.WithDisplayName(v))
	}
	if v, _ := cmd.Flags().GetString("password"); v != "" {
		opts = append(opts, enclave.WithPassword(v))
	}

	u, err := c.UpdateUser(cmd.Context(), args[0], opts...)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return printer.Print([]any{u})
}
