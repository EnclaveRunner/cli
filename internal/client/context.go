package client

import (
	"cli/internal/config"
	"context"

	"github.com/EnclaveRunner/sdk-go/enclave"
)

type contextKey int

const (
	clientKey contextKey = iota
	configKey
)

// WithClient stores the SDK client in the context.
func WithClient(ctx context.Context, c *enclave.Client) context.Context {
	return context.WithValue(ctx, clientKey, c)
}

// FromContext retrieves the SDK client from the context.
// Panics if not set — callers must go through PersistentPreRunE.
func FromContext(ctx context.Context) *enclave.Client {
	c, _ := ctx.Value(clientKey).(*enclave.Client)

	return c
}

// WithConfig stores the config in the context.
func WithConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

// ConfigFromContext retrieves the config from the context.
func ConfigFromContext(ctx context.Context) *config.Config {
	cfg, _ := ctx.Value(configKey).(*config.Config)

	return cfg
}
