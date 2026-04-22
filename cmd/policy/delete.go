package policy

import (
	"fmt"
	"os"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an RBAC policy",
		RunE:  runDelete,
	}
	cmd.Flags().String("role", "", "Role name (required)")
	cmd.Flags().String("resource-group", "", "Resource group name (required)")
	cmd.Flags().String("method", "", "HTTP method (required)")
	_ = cmd.MarkFlagRequired("role")
	_ = cmd.MarkFlagRequired("resource-group")
	_ = cmd.MarkFlagRequired("method")
	return cmd
}

func runDelete(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.PolicyColumns, os.Stdout)

	role, _ := cmd.Flags().GetString("role")
	rg, _ := cmd.Flags().GetString("resource-group")
	method, _ := cmd.Flags().GetString("method")

	p := enclave.Policy{
		Role:          role,
		ResourceGroup: rg,
		Method:        enclave.PolicyMethod(method),
	}
	if err := c.DeletePolicy(cmd.Context(), p); err != nil {
		return fmt.Errorf("delete policy: %w", err)
	}
	return printer.Print([]any{p})
}
