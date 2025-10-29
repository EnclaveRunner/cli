package cmd

import (
	"cli/config"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	v       *viper.Viper
)

var rootCmd = &cobra.Command{
	Use:   "enclave <command>",
	Short: "Manage you enclave platform from your terminal",
	Long: `Enclave CLI is a command line interface to manage your enclave platform.
Run the cli without a command to start the interactive tui or use one of the available commands
to perform specific actions directly.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false},
		)

		// If a config file is specified, set it in viper
		if cfgFile != "" {
			v.SetConfigFile(cfgFile)
		}

		config.Init(v)
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Error().Msg("TUI is currently not implemented.")
	},
	CompletionOptions: cobra.CompletionOptions{},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	v = viper.New()

	// Persistent flags (available to all subcommands)
	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default searches ./cli.yaml, $HOME/.enclave/cli.yaml, /etc/enclave/cli.yaml)")
	rootCmd.PersistentFlags().String("api-url", "", "API server URL")
	rootCmd.PersistentFlags().
		String("auth-username", "", "Authentication username")
	rootCmd.PersistentFlags().
		String("auth-password", "", "Authentication password")

	// Bind flags to viper
	err := v.BindPFlag(
		"api_server_url",
		rootCmd.PersistentFlags().Lookup("api-url"),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind flag")
	}
	err = v.BindPFlag(
		"auth.username",
		rootCmd.PersistentFlags().Lookup("auth-username"),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind flag")
	}
	err = v.BindPFlag(
		"auth.password",
		rootCmd.PersistentFlags().Lookup("auth-password"),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind flag")
	}
}
