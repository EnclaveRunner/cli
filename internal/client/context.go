package client

import (
	"context"

	"cli/internal/config"

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
	return ctx.Value(clientKey).(*enclave.Client)
}

// WithConfig stores the config in the context.
func WithConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

// ConfigFromContext retrieves the config from the context.
func ConfigFromContext(ctx context.Context) *config.Config {
	return ctx.Value(configKey).(*config.Config)
}
