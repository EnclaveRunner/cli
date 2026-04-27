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

type usersMode int

const (
	usersModeList     usersMode = iota
	usersModeDescribe           // enter key: full detail view
	usersModeModal              // d key: confirm delete
	usersModeForm               // c key: create form
)

// UsersModel is the user list view.
type UsersModel struct {
	Users     []enclave.User
	Cursor    int
	Loading   bool
	Err       error
	colOffset int
	width     int
	height    int

	mode  usersMode
	modal ModalModel
	form  FormModel
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
func (m *UsersModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.modal.SetSize(w, h)
	m.form.SetSize(w, h)
}

// IsCapturing reports whether the view is in a mode that owns the keyboard
// (form, modal, or describe), so the parent can suppress global hotkeys.
func (m UsersModel) IsCapturing() bool {
	return m.mode != usersModeList
}

// selectedUser returns the user at the current cursor, if any.
func (m UsersModel) selectedUser() (enclave.User, bool) {
	if len(m.Users) == 0 || m.Cursor >= len(m.Users) {
		return enclave.User{}, false
	}

	return m.Users[m.Cursor], true
}

// Update handles messages.
func (m UsersModel) Update(
	msg tea.Msg,
) (UsersModel, tea.Cmd) {
	switch m.mode {
	case usersModeModal:
		return m.updateModal(msg)
	case usersModeForm:
		return m.updateForm(msg)
	case usersModeDescribe:
		return m.updateDescribe(msg)
	}

	return m.updateList(msg)
}

func (m UsersModel) updateList(msg tea.Msg) (UsersModel, tea.Cmd) {
	switch msg := msg.(type) {
	case UsersLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.Users = msg.Users
		m.Cursor = 0

	case UserDeletedMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
		}

	case UserCreatedMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
		}

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

		case "enter":
			if _, ok := m.selectedUser(); ok {
				m.mode = usersModeDescribe
			}

		case "d":
			if u, ok := m.selectedUser(); ok {
				m.modal = NewModal("Delete user \"" + u.Name + "\"?")
				m.modal.SetSize(m.width, m.height)
				m.mode = usersModeModal
			}

		case "c":
			m.form = NewForm("Create User", []FormField{
				{Label: "Username", Placeholder: "alice"},
				{Label: "Display Name", Placeholder: "Alice Smith"},
				{Label: "Password", Placeholder: "••••••••", Secret: true},
			})
			m.form.SetSize(m.width, m.height)
			m.mode = usersModeForm
		}
	}

	return m, nil
}

func (m UsersModel) updateModal(msg tea.Msg) (UsersModel, tea.Cmd) {
	switch msg.(type) {
	case ModalConfirmedMsg:
		if u, ok := m.selectedUser(); ok {
			m.mode = usersModeList
			m.Loading = true
			name := u.Name

			return m, func() tea.Msg { return FormDeleteUserMsg{Name: name} }
		}
		m.mode = usersModeList

	case ModalCancelledMsg:
		m.mode = usersModeList

	default:
		var cmd tea.Cmd
		m.modal, cmd = m.modal.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m UsersModel) updateForm(msg tea.Msg) (UsersModel, tea.Cmd) {
	switch msg := msg.(type) {
	case FormSubmittedMsg:
		if len(msg.Values) >= 3 {
			name, display, pass := msg.Values[0], msg.Values[1], msg.Values[2]
			if name == "" {
				m.form.SetError("username is required")

				return m, nil
			}
			m.mode = usersModeList

			return m, func() tea.Msg { return FormCreateUserMsg{Name: name, Display: display, Pass: pass} }
		}
		m.mode = usersModeList

	case FormCancelledMsg:
		m.mode = usersModeList

	default:
		var cmd tea.Cmd
		m.form, cmd = m.form.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m UsersModel) updateDescribe(msg tea.Msg) (UsersModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc", "q":
			m.mode = usersModeList
		}
	}

	return m, nil
}

// View renders the users table or the active overlay.
func (m UsersModel) View() string {
	switch m.mode {
	case usersModeDescribe:
		return m.renderDescribe()
	case usersModeModal:
		return m.renderList() + m.modal.View()
	case usersModeForm:
		return m.form.View()
	}

	return m.renderList()
}

func (m UsersModel) renderList() string {
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

func (m UsersModel) renderDescribe() string {
	u, ok := m.selectedUser()
	if !ok {
		return styles.MutedStyle.Render("\n  No user selected.")
	}

	field := func(label, value string) string {
		return styles.MutedStyle.Render(label+": ") + value
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render("User "+u.Name) + "\n\n")
	b.WriteString(field("Name", u.Name) + "\n")
	b.WriteString(field("Display Name", u.DisplayName) + "\n")

	if len(u.Roles) > 0 {
		b.WriteString(field("Roles", strings.Join(u.Roles, ", ")) + "\n")
	} else {
		b.WriteString(field("Roles", styles.MutedStyle.Render("none")) + "\n")
	}

	b.WriteString("\n" + styles.HelpKeyStyle.Render("esc") +
		lipgloss.NewStyle().Foreground(styles.ColorSlateDark).Render(" back"))

	return b.String()
}

// FormCreateUserMsg carries the values for an async user create operation.
// It is handled in app.go.
type FormCreateUserMsg struct {
	Name    string
	Display string
	Pass    string
}
