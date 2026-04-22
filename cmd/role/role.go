package role

import "github.com/spf13/cobra"

// NewCmd returns the "role" command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage roles",
	}
	cmd.AddCommand(
		newListCmd(),
		newGetCmd(),
		newCreateCmd(),
		newDeleteCmd(),
	)

	return cmd
}
