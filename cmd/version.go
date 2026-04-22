package cmd

import (
	_ "embed"
	"fmt"
	"strings"

	iv "cli/internal/version"

	"github.com/spf13/cobra"
)

//go:embed version.txt
var versionFile string

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the encl version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			local := strings.TrimSpace(versionFile)
			_, err := fmt.Fprintln(cmd.OutOrStdout(), local)
			if err != nil {
				return err
			}

			// Check remote version (best-effort)
			remote, newer, err := iv.CheckRemote(local)
			if err == nil && newer {
				fmt.Fprintln(cmd.OutOrStdout(), "New version available:", remote)
			}

			return nil
		},
	}
}
