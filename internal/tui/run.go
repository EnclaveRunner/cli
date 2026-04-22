package tui

import (
	"fmt"

	"github.com/EnclaveRunner/sdk-go/enclave"
	tea "github.com/charmbracelet/bubbletea"
)

// Run launches the Bubbletea TUI program.
func Run(c *enclave.Client) error {
	return RunWithConfig(c, "", "", "")
}

// RunWithConfig launches the TUI with config info for the header panel.
func RunWithConfig(c *enclave.Client, apiURL, username, version string) error {
	m := New(c, apiURL, username, version)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tui: %w", err)
	}
	return nil
}
