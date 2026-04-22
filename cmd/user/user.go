package user

import (
	"cli/cmd/user/me"

	"github.com/spf13/cobra"
)

// NewCmd returns the "user" command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}
	cmd.AddCommand(
		newListCmd(),
		newGetCmd(),
		newCreateCmd(),
		newUpdateCmd(),
		newDeleteCmd(),
		me.NewCmd(),
	)
	return cmd
}
