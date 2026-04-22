package me

import (
	"fmt"
	"os"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the currently authenticated user",
		RunE:  runUpdate,
	}
	cmd.Flags().String("display-name", "", "New display name")
	cmd.Flags().String("password", "", "New password")
	return cmd
}

func runUpdate(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.UserColumns, os.Stdout)

	var opts []enclave.UpdateUserOption
	if v, _ := cmd.Flags().GetString("display-name"); v != "" {
		opts = append(opts, enclave.WithDisplayName(v))
	}
	if v, _ := cmd.Flags().GetString("password"); v != "" {
		opts = append(opts, enclave.WithPassword(v))
	}

	u, err := c.UpdateMe(cmd.Context(), opts...)
	if err != nil {
		return fmt.Errorf("update me: %w", err)
	}
	return printer.Print([]any{u})
}
