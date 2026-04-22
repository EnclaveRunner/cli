package policy

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an RBAC policy",
		RunE:  runDelete,
	}
	addPolicyFlags(cmd, "HTTP method")

	return cmd
}

func runDelete(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.PolicyColumns,
		os.Stdout,
	)

	p := policyFromFlags(cmd)
	if err := c.DeletePolicy(cmd.Context(), p); err != nil {
		return fmt.Errorf("delete policy: %w", err)
	}

	return printer.Print([]any{p})
}
