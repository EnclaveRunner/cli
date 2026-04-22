package resourcegroup

import "github.com/spf13/cobra"

// NewCmd returns the "resource-group" command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resource-group",
		Aliases: []string{"rg"},
		Short:   "Manage resource groups",
	}
	cmd.AddCommand(
		newListCmd(),
		newGetCmd(),
		newCreateCmd(),
		newDeleteCmd(),
	)

	return cmd
}
