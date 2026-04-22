package policy

import (
	"fmt"
	"os"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an RBAC policy",
		RunE:  runCreate,
	}
	cmd.Flags().String("role", "", "Role name (required)")
	cmd.Flags().String("resource-group", "", "Resource group name (required)")
	cmd.Flags().String("method", "", "HTTP method: GET, POST, PUT, PATCH, DELETE, HEAD, * (required)")
	_ = cmd.MarkFlagRequired("role")
	_ = cmd.MarkFlagRequired("resource-group")
	_ = cmd.MarkFlagRequired("method")
	return cmd
}

func runCreate(cmd *cobra.Command, _ []string) error {
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
	if err := c.CreatePolicy(cmd.Context(), p); err != nil {
		return fmt.Errorf("create policy: %w", err)
	}
	return printer.Print([]any{p})
}
