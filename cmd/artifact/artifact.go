package artifact

import "github.com/spf13/cobra"

// NewCmd returns the "artifact" command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Manage artifacts",
	}
	cmd.AddCommand(
		newNamespaceCmd(),
		newListCmd(),
		newVersionsCmd(),
		newUploadCmd(),
		newGetCmd(),
		newDownloadCmd(),
		newTagCmd(),
		newDeleteCmd(),
	)
	return cmd
}
