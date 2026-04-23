package cmd

import (
	"fmt"

	iv "cli/internal/version"

	"github.com/spf13/cobra"
)

var appVersion string

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the encl version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), appVersion)
			if err != nil {
				return err
			}

			// Check remote version (best-effort)
			remote, newer, err := iv.CheckRemote(appVersion)
			if err == nil && newer {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "New version available:", remote)
			}

			return nil
		},
	}
}
