package client

import (
	"fmt"

	"cli/internal/config"

	"github.com/EnclaveRunner/sdk-go/enclave"
)

// New constructs an authenticated Enclave SDK client from cfg.
// Returns an error if api_url, username, or password are unset.
func New(cfg *config.Config) (*enclave.Client, error) {
	if cfg.APIURL == "" {
		return nil, fmt.Errorf("api_url is required (set --api-url, ENCLAVE_API_URL, or api_url in config)")
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("username is required (set --username, ENCLAVE_USERNAME, or username in config)")
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("password is required (set --password, ENCLAVE_PASSWORD, or password in config)")
	}
	c, err := enclave.New(cfg.APIURL, cfg.Username, cfg.Password)
	if err != nil {
		return nil, err
	}
	return c, nil
}
