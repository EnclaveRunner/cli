package views

import (
	"cli/internal/styles"
	"context"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/EnclaveRunner/sdk-go/enclave"
	tea "github.com/charmbracelet/bubbletea"
)

// RolesLoadedMsg carries loaded roles.
type RolesLoadedMsg struct {
	Roles []enclave.Role
	Err   error
}

// RolesModel is the roles list view.
//

type RolesModel struct {
	Roles     []enclave.Role
	Cursor    int
	Loading   bool
	Err       error
	colOffset int
	width     int
	height    int
}

// Load fetches all roles.
func (m RolesModel) Load(
	c *enclave.Client,
) tea.Cmd {
	return func() tea.Msg {
		roles, err := enclave.Collect(c.ListRoles(context.Background()))

		return RolesLoadedMsg{Roles: roles, Err: err}
	}
}

// SetSize updates the rendering area.
func (m *RolesModel) SetSize(w, h int) { m.width = w; m.height = h }

// Update handles messages.
func (m RolesModel) Update(
	msg tea.Msg,
) (RolesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case RolesLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.Roles = msg.Roles
		m.Cursor = 0
	case tea.KeyMsg:
		switch msg.String() {
		case keyUp, keyK:
			if m.Cursor > 0 {
				m.Cursor--
			}
		case keyDown, keyJ:
			if m.Cursor < len(m.Roles)-1 {
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

// View renders the roles table.
func (m RolesModel) View() string {
	if m.Loading {
		return styles.MutedStyle.Render("\n  Loading roles…")
	}
	if m.Err != nil {
		return styles.ErrorStyle.Render("\n  Error: " + m.Err.Error())
	}
	if len(m.Roles) == 0 {
		return styles.MutedStyle.Render("\n  No roles found.")
	}

	headers := []string{"NAME", "USERS"}
	colWidths := []int{len(headers[0]), len(headers[1])}

	rows := make([][]string, len(m.Roles))
	for i, r := range m.Roles {
		count := strconv.Itoa(len(r.Users))
		rows[i] = []string{r.Name, count}
		updateWidth(&colWidths[0], len(r.Name))
		updateWidth(&colWidths[1], len(count))
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
