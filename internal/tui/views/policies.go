package views

import (
	"cli/internal/styles"
	"context"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/EnclaveRunner/sdk-go/enclave"
	tea "github.com/charmbracelet/bubbletea"
)

// PoliciesLoadedMsg carries loaded policies.
type PoliciesLoadedMsg struct {
	Policies []enclave.Policy
	Err      error
}

// PoliciesModel is the policies list view.
type PoliciesModel struct {
	Policies  []enclave.Policy
	Cursor    int
	Loading   bool
	Err       error
	colOffset int
	width     int
	height    int
}

// Load fetches all policies.
func (m PoliciesModel) Load(
	c *enclave.Client,
) tea.Cmd {
	return func() tea.Msg {
		policies, err := enclave.Collect(c.ListPolicies(context.Background()))

		return PoliciesLoadedMsg{Policies: policies, Err: err}
	}
}

// SetSize updates the rendering area.
func (m *PoliciesModel) SetSize(w, h int) { m.width = w; m.height = h }

// Update handles messages.
func (m PoliciesModel) Update(
	msg tea.Msg,
) (PoliciesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case PoliciesLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.Policies = msg.Policies
		m.Cursor = 0
	case tea.KeyMsg:
		switch msg.String() {
		case keyUp, keyK:
			if m.Cursor > 0 {
				m.Cursor--
			}
		case keyDown, keyJ:
			if m.Cursor < len(m.Policies)-1 {
				m.Cursor++
			}
		case keyLeft:
			if m.colOffset > 0 {
				m.colOffset--
			}
		case keyRight:
			m.colOffset++
		}
	}

	return m, nil
}

// View renders the policies table.
func (m PoliciesModel) View() string {
	if m.Loading {
		return styles.MutedStyle.Render("\n  Loading policies…")
	}
	if m.Err != nil {
		return styles.ErrorStyle.Render("\n  Error: " + m.Err.Error())
	}
	if len(m.Policies) == 0 {
		return styles.MutedStyle.Render("\n  No policies found.")
	}

	headers := []string{"ROLE", "RESOURCE GROUP", "METHOD"}
	colWidths := []int{len(headers[0]), len(headers[1]), len(headers[2])}

	rows := make([][]string, len(m.Policies))
	for i, p := range m.Policies {
		method := string(p.Method)
		rows[i] = []string{p.Role, p.ResourceGroup, method}
		updateWidth(&colWidths[0], len(p.Role))
		updateWidth(&colWidths[1], len(p.ResourceGroup))
		updateWidth(&colWidths[2], len(method))
	}

	var b strings.Builder

	startCol := m.colOffset
	if startCol >= len(headers) {
		startCol = len(headers) - 1
	}

	headerCells := make([]string, len(headers))
	for i, h := range headers {
		headerCells[i] = styles.HeaderStyle.Render(padRight(h, colWidths[i]))
	}
	b.WriteString(strings.Join(headerCells[startCol:], "") + "\n")

	for i, row := range rows {
		cells := make([]string, len(row))
		for j, cell := range row {
			if i == m.Cursor {
				cells[j] = lipgloss.NewStyle().Padding(0, 1).
					Background(styles.ColorSecondaryGreen).
					Foreground(styles.ColorNearBlack).
					Render(padRight(cell, colWidths[j]))
			} else {
				cells[j] = " " + padRight(cell, colWidths[j]) + " "
			}
		}
		b.WriteString(strings.Join(cells[startCol:], "") + "\n")
	}

	return b.String()
}
