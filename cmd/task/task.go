package task

import "github.com/spf13/cobra"

// NewCmd returns the "task" command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
	}
	cmd.AddCommand(
		newListCmd(),
		newGetCmd(),
		newCreateCmd(),
		newLogsCmd(),
	)

	return cmd
}
