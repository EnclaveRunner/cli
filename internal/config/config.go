package config

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all resolved configuration values.
type Config struct {
	APIURL   string `mapstructure:"api_url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	LogLevel string `mapstructure:"log_level"`
	Output   string `mapstructure:"output"`
}

// OutputFormat returns the output format as a string (table, json, yaml).
func (c *Config) OutputFormat() string {
	if c.Output == "" {
		return "table"
	}
	return c.Output
}

// Load initialises Viper, binds pflags, reads config file(s), and returns
// a populated Config. flags may be nil.
func Load(flags *pflag.FlagSet) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME/.enclave")
	v.AddConfigPath("./.enclave")
	v.AddConfigPath("/etc/enclave")

	v.SetEnvPrefix("ENCLAVE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("log_level", "info")
	v.SetDefault("output", "table")

	if flags != nil {
		if f := flags.Lookup("api-url"); f != nil {
			_ = v.BindPFlag("api_url", f)
		}
		if f := flags.Lookup("username"); f != nil {
			_ = v.BindPFlag("username", f)
		}
		if f := flags.Lookup("password"); f != nil {
			_ = v.BindPFlag("password", f)
		}
		if f := flags.Lookup("log-level"); f != nil {
			_ = v.BindPFlag("log_level", f)
		}
		if f := flags.Lookup("output"); f != nil {
			_ = v.BindPFlag("output", f)
		}
	}

	// Ignore config file not found; all settings may come from env/flags.
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
