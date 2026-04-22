package policy

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an RBAC policy",
		RunE:  runCreate,
	}
	addPolicyFlags(cmd, "HTTP method: GET, POST, PUT, PATCH, DELETE, HEAD, *")

	return cmd
}

func runCreate(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.PolicyColumns,
		os.Stdout,
	)

	p := policyFromFlags(cmd)
	if err := c.CreatePolicy(cmd.Context(), p); err != nil {
		return fmt.Errorf("create policy: %w", err)
	}

	return printer.Print([]any{p})
}
