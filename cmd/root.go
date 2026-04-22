package cmd

import (
	"cli/cmd/artifact"
	"cli/cmd/policy"
	"cli/cmd/resourcegroup"
	"cli/cmd/role"
	"cli/cmd/task"
	"cli/cmd/user"
	"cli/internal/client"
	"cli/internal/config"
	"cli/internal/tui"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:          "encl",
	Short:        "Enclave CLI — manage users, roles, tasks, and artifacts",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// Skip setup for commands that don't need the SDK client.
		if cmd.Name() == "version" || cmd.Name() == "help" ||
			cmd.Name() == "completion" {
			return nil
		}

		cfg, err := config.Load(cmd.Root().PersistentFlags())
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		// Initialise zerolog with human-readable console output.
		level, err := zerolog.ParseLevel(cfg.LogLevel)
		if err != nil {
			level = zerolog.InfoLevel
		}
		zerolog.SetGlobalLevel(level)
		log.Logger = zerolog.New(
			zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"},
		).With().Timestamp().Logger()

		// Build the SDK client.
		c, err := client.New(cfg)
		if err != nil {
			return err
		}

		// Store both in the command context for subcommands.
		ctx := client.WithClient(cmd.Context(), c)
		ctx = client.WithConfig(ctx, cfg)
		cmd.SetContext(ctx)

		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		// When run with no subcommand and attached to a TTY, launch TUI.
		if term.IsTerminal(int(os.Stdout.Fd())) {
			c := client.FromContext(cmd.Context())
			cfg := client.ConfigFromContext(cmd.Context())

			return tui.RunWithConfig(
				c,
				cfg.APIURL,
				cfg.Username,
				strings.TrimSpace(versionFile),
			)
		}

		return cmd.Help()
	},
}

// Execute is the entry point called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.String(
		"api-url",
		"",
		"Enclave API URL (overrides config and ENCLAVE_API_URL)",
	)
	pf.String("username", "", "Username (overrides config and ENCLAVE_USERNAME)")
	pf.String("password", "", "Password (overrides config and ENCLAVE_PASSWORD)")
	pf.String(
		"log-level",
		"",
		"Log level: trace, debug, info, warn, error (default: info)",
	)
	pf.String("output", "table", "Output format: table, json, yaml")

	rootCmd.AddCommand(
		user.NewCmd(),
		role.NewCmd(),
		resourcegroup.NewCmd(),
		policy.NewCmd(),
		task.NewCmd(),
		artifact.NewCmd(),
		newVersionCmd(),
	)
}
