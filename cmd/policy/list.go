package policy

import (
	"fmt"
	"os"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List RBAC policies",
		RunE:  runList,
	}
	cmd.Flags().String("role", "", "Filter by role")
	cmd.Flags().String("resource-group", "", "Filter by resource group")
	cmd.Flags().String("method", "", "Filter by HTTP method")
	return cmd
}

func runList(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.PolicyColumns, os.Stdout)

	var opts []enclave.ListPoliciesOption
	if v, _ := cmd.Flags().GetString("role"); v != "" {
		opts = append(opts, enclave.FilterPolicyByRole(v))
	}
	if v, _ := cmd.Flags().GetString("resource-group"); v != "" {
		opts = append(opts, enclave.FilterPolicyByResourceGroup(v))
	}
	if v, _ := cmd.Flags().GetString("method"); v != "" {
		opts = append(opts, enclave.FilterPolicyByMethod(v))
	}

	policies, err := enclave.Collect(c.ListPolicies(cmd.Context(), opts...))
	if err != nil {
		return fmt.Errorf("list policies: %w", err)
	}
	return printer.Print(policies)
}
