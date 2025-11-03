package config

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Cfg is the global configuration instance
var Cfg Config

var Version string

type Config struct {
	APIServerURL string     `mapstructure:"api_server_url"`
	Auth         AuthConfig `mapstructure:"-"`
}

type AuthConfig interface {
	GetAuthHeader() string
}

type BasicAuth struct {
	Username string
	Password string
}

func (authCfg BasicAuth) GetAuthHeader() string {
	auth := authCfg.Username + ":" + authCfg.Password

	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func Init(v *viper.Viper) {
	// Set up environment variable support
	v.SetEnvPrefix("ENCLAVE") // will be uppercased automatically
	v.AutomaticEnv()

	// Bind specific environment variables
	err := v.BindEnv("api_server_url", "ENCLAVE_API_SERVER_URL")
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind environment variable")
	}
	err = v.BindEnv("auth.username", "ENCLAVE_AUTH_USERNAME")
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind environment variable")
	}
	err = v.BindEnv("auth.password", "ENCLAVE_AUTH_PASSWORD")
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind environment variable")
	}

	// Set defaults
	v.SetDefault("api_server_url", "https://api.enclave.io")

	// Try to read config file (only if explicitly set or if file exists)
	configFileSet := v.ConfigFileUsed() != ""
	if !configFileSet {
		// Set config name and type
		v.SetConfigName("cli.yml")
		v.SetConfigType("yaml")

		// Add config paths in order of precedence (last added has lowest
		// precedence)
		// 1. Current directory (highest precedence)
		v.AddConfigPath(".")

		// 2. Home directory
		home, err := os.UserHomeDir()
		if err == nil {
			v.AddConfigPath(filepath.Join(home, ".enclave"))
		}

		// 3. System-wide config (lowest precedence)
		v.AddConfigPath("/etc/enclave")
	}

	// Read config file (it's okay if it doesn't exist)
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Debug().Err(err).Msg("Config file not found")
			// Config file not found; ignore error and use defaults/env vars
		} else {
			// Config file was found but another error was produced
			// Only fail if config file was explicitly set via --config flag
			if configFileSet {
				log.Error().Err(err).Msg("Failed to load config file")
			} else {
				log.Debug().Err(err).Msg("Failed to load config file")
			}
			// Otherwise, ignore parse errors and use defaults/env vars
		}
	}

	// Unmarshal config into the global Cfg variable (excluding Auth field)
	if err := v.Unmarshal(&Cfg); err != nil {
		log.Error().Err(err).Msg("Failed unmarshaling config")
	}

	// Handle auth configuration manually since it's an interface
	username := v.GetString("auth.username")
	password := v.GetString("auth.password")
	if username != "" || password != "" {
		Cfg.Auth = BasicAuth{
			Username: username,
			Password: password,
		}
	}
}
