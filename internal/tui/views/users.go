package views

import (
	"cli/internal/styles"
	"context"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/EnclaveRunner/sdk-go/enclave"
	tea "github.com/charmbracelet/bubbletea"
)

// UsersLoadedMsg carries loaded users.
type UsersLoadedMsg struct {
	Users []enclave.User
	Err   error
}

// UsersModel is the user list view.
type UsersModel struct {
	Users     []enclave.User
	Cursor    int
	Loading   bool
	Err       error
	colOffset int
	width     int
	height    int
}

// Load fetches all users.
func (m UsersModel) Load(
	c *enclave.Client,
) tea.Cmd {
	return func() tea.Msg {
		users, err := enclave.Collect(c.ListUsers(context.Background()))

		return UsersLoadedMsg{Users: users, Err: err}
	}
}

// SetSize updates the rendering area.
func (m *UsersModel) SetSize(w, h int) { m.width = w; m.height = h }

// Update handles messages.
func (m UsersModel) Update(
	msg tea.Msg,
) (UsersModel, tea.Cmd) {
	switch msg := msg.(type) {
	case UsersLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.Users = msg.Users
		m.Cursor = 0
	case tea.KeyMsg:
		switch msg.String() {
		case keyUp, keyK:
			if m.Cursor > 0 {
				m.Cursor--
			}
		case keyDown, keyJ:
			if m.Cursor < len(m.Users)-1 {
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

// View renders the users table.
func (m UsersModel) View() string {
	if m.Loading {
		return styles.MutedStyle.Render("\n  Loading users…")
	}
	if m.Err != nil {
		return styles.ErrorStyle.Render("\n  Error: " + m.Err.Error())
	}
	if len(m.Users) == 0 {
		return styles.MutedStyle.Render("\n  No users found.")
	}

	headers := []string{"NAME", "DISPLAY NAME", "ROLES"}
	colWidths := []int{len(headers[0]), len(headers[1]), len(headers[2])}

	rows := make([][]string, len(m.Users))
	for i, u := range m.Users {
		roles := strings.Join(u.Roles, ", ")
		rows[i] = []string{u.Name, u.DisplayName, roles}
		updateWidth(&colWidths[0], len(u.Name))
		updateWidth(&colWidths[1], len(u.DisplayName))
		updateWidth(&colWidths[2], len(roles))
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
