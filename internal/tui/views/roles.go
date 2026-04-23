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

type rolesMode int

const (
	rolesModeList     rolesMode = iota
	rolesModeDescribe           // enter: full detail
	rolesModeModal              // d: confirm delete
	rolesModeForm               // c: create form
)

// RolesModel is the roles list view.
type RolesModel struct {
	Roles     []enclave.Role
	Cursor    int
	Loading   bool
	Err       error
	colOffset int
	width     int
	height    int

	mode  rolesMode
	modal ModalModel
	form  FormModel
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
func (m *RolesModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.modal.SetSize(w, h)
	m.form.SetSize(w, h)
}

// IsCapturing reports whether the view is in a mode that owns the keyboard.
func (m RolesModel) IsCapturing() bool {
	return m.mode != rolesModeList
}

func (m RolesModel) selectedRole() (enclave.Role, bool) {
	if len(m.Roles) == 0 || m.Cursor >= len(m.Roles) {
		return enclave.Role{}, false
	}

	return m.Roles[m.Cursor], true
}

// Update handles messages.
func (m RolesModel) Update(
	msg tea.Msg,
) (RolesModel, tea.Cmd) {
	switch m.mode {
	case rolesModeModal:
		return m.updateModal(msg)
	case rolesModeForm:
		return m.updateForm(msg)
	case rolesModeDescribe:
		return m.updateDescribe(msg)
	}

	return m.updateList(msg)
}

func (m RolesModel) updateList(msg tea.Msg) (RolesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case RolesLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.Roles = msg.Roles
		m.Cursor = 0

	case RoleDeletedMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
		}

	case RoleCreatedMsg:
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
			if m.Cursor < len(m.Roles)-1 {
				m.Cursor++
			}
		case keyLeft:
			if m.colOffset > 0 {
				m.colOffset--
			}
		case keyRight:
			m.colOffset++

		case "enter":
			if _, ok := m.selectedRole(); ok {
				m.mode = rolesModeDescribe
			}

		case "d":
			if r, ok := m.selectedRole(); ok {
				m.modal = NewModal("Delete role \"" + r.Name + "\"?")
				m.modal.SetSize(m.width, m.height)
				m.mode = rolesModeModal
			}

		case "c":
			m.form = NewForm("Create Role", []FormField{
				{Label: "Name", Placeholder: "admin"},
				{Label: "Users (comma-sep)", Placeholder: "alice, bob"},
			})
			m.form.SetSize(m.width, m.height)
			m.mode = rolesModeForm
		}
	}

	return m, nil
}

func (m RolesModel) updateModal(msg tea.Msg) (RolesModel, tea.Cmd) {
	switch msg.(type) {
	case ModalConfirmedMsg:
		if r, ok := m.selectedRole(); ok {
			m.mode = rolesModeList
			m.Loading = true
			name := r.Name

			return m, func() tea.Msg { return FormDeleteRoleMsg{Name: name} }
		}
		m.mode = rolesModeList

	case ModalCancelledMsg:
		m.mode = rolesModeList

	default:
		var cmd tea.Cmd
		m.modal, cmd = m.modal.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m RolesModel) updateForm(msg tea.Msg) (RolesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case FormSubmittedMsg:
		if len(msg.Values) >= 1 {
			name := msg.Values[0]
			if name == "" {
				m.form.SetError("name is required")

				return m, nil
			}
			usersRaw := ""
			if len(msg.Values) >= 2 {
				usersRaw = msg.Values[1]
			}
			m.mode = rolesModeList

			return m, func() tea.Msg { return FormCreateRoleMsg{Name: name, UsersRaw: usersRaw} }
		}
		m.mode = rolesModeList

	case FormCancelledMsg:
		m.mode = rolesModeList

	default:
		var cmd tea.Cmd
		m.form, cmd = m.form.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m RolesModel) updateDescribe(msg tea.Msg) (RolesModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc", "q":
			m.mode = rolesModeList
		}
	}

	return m, nil
}

// View renders the roles table or the active overlay.
func (m RolesModel) View() string {
	switch m.mode {
	case rolesModeDescribe:
		return m.renderDescribe()
	case rolesModeModal:
		return m.renderList() + m.modal.View()
	case rolesModeForm:
		return m.form.View()
	}

	return m.renderList()
}

func (m RolesModel) renderList() string {
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

func (m RolesModel) renderDescribe() string {
	r, ok := m.selectedRole()
	if !ok {
		return styles.MutedStyle.Render("\n  No role selected.")
	}

	field := func(label, value string) string {
		return styles.MutedStyle.Render(label+": ") + value
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render("Role "+r.Name) + "\n\n")
	b.WriteString(field("Name", r.Name) + "\n")
	b.WriteString(field("User count", strconv.Itoa(len(r.Users))) + "\n")

	if len(r.Users) > 0 {
		b.WriteString(field("Users", strings.Join(r.Users, ", ")) + "\n")
	}

	b.WriteString("\n" + styles.HelpKeyStyle.Render("esc") +
		lipgloss.NewStyle().Foreground(styles.ColorSlateDark).Render(" back"))

	return b.String()
}

// FormDeleteRoleMsg requests an async role delete in app.go.
type FormDeleteRoleMsg struct{ Name string }

// FormCreateRoleMsg requests an async role create in app.go.
type FormCreateRoleMsg struct {
	Name     string
	UsersRaw string // comma-separated user list
}
