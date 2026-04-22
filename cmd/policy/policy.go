package policy

import "github.com/spf13/cobra"

// NewCmd returns the "policy" command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage RBAC policies",
	}
	cmd.AddCommand(
		newListCmd(),
		newCreateCmd(),
		newDeleteCmd(),
	)
	return cmd
}
