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

type policiesMode int

const (
	policiesModeList  policiesMode = iota
	policiesModeModal              // d: confirm delete
	policiesModeForm               // c: create form
)

// PoliciesModel is the policies list view.
type PoliciesModel struct {
	Policies  []enclave.Policy
	Cursor    int
	Loading   bool
	Err       error
	colOffset int
	width     int
	height    int

	mode  policiesMode
	modal ModalModel
	form  FormModel
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
func (m *PoliciesModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.modal.SetSize(w, h)
	m.form.SetSize(w, h)
}

// IsCapturing reports whether the view is in a mode that owns the keyboard.
func (m PoliciesModel) IsCapturing() bool {
	return m.mode != policiesModeList
}

// Update handles messages.
func (m PoliciesModel) Update(
	msg tea.Msg,
) (PoliciesModel, tea.Cmd) {
	switch m.mode {
	case policiesModeModal:
		return m.updateModal(msg)
	case policiesModeForm:
		return m.updateForm(msg)
	case policiesModeList:
		return m.updateList(msg)
	}

	return m.updateList(msg)
}

// View renders the policies table or the active overlay.
func (m PoliciesModel) View() string {
	switch m.mode {
	case policiesModeModal:
		return m.renderList() + m.modal.View()
	case policiesModeForm:
		return m.form.View()
	case policiesModeList:
		return m.renderList()
	}

	return m.renderList()
}

func (m PoliciesModel) selectedPolicy() (enclave.Policy, bool) {
	if len(m.Policies) == 0 || m.Cursor >= len(m.Policies) {
		return enclave.Policy{}, false
	}

	return m.Policies[m.Cursor], true
}

func (m PoliciesModel) updateList(msg tea.Msg) (PoliciesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case PoliciesLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.Policies = msg.Policies
		m.Cursor = 0

	case PolicyDeletedMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
		}

	case PolicyCreatedMsg:
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
			if m.Cursor < len(m.Policies)-1 {
				m.Cursor++
			}
		case keyLeft:
			if m.colOffset > 0 {
				m.colOffset--
			}
		case keyRight:
			m.colOffset++

		case "d":
			if p, ok := m.selectedPolicy(); ok {
				label := p.Role + " / " + p.ResourceGroup + " / " + string(p.Method)
				m.modal = NewModal("Delete policy \"" + label + "\"?")
				m.modal.SetSize(m.width, m.height)
				m.mode = policiesModeModal
			}

		case "c":
			m.form = NewForm("Create Policy", []FormField{
				{Label: "Role", Placeholder: "admin"},
				{Label: "Resource Group", Placeholder: "my-api"},
				{Label: "Method", Placeholder: "GET, POST, PUT, DELETE, *"},
			})
			m.form.SetSize(m.width, m.height)
			m.mode = policiesModeForm
		}
	}

	return m, nil
}

func (m PoliciesModel) updateModal(msg tea.Msg) (PoliciesModel, tea.Cmd) {
	switch msg.(type) {
	case ModalConfirmedMsg:
		if p, ok := m.selectedPolicy(); ok {
			m.mode = policiesModeList
			m.Loading = true
			policy := p

			return m, func() tea.Msg { return FormDeletePolicyMsg{Policy: policy} }
		}
		m.mode = policiesModeList

	case ModalCancelledMsg:
		m.mode = policiesModeList

	default:
		var cmd tea.Cmd
		m.modal, cmd = m.modal.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m PoliciesModel) updateForm(msg tea.Msg) (PoliciesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case FormSubmittedMsg:
		if len(msg.Values) >= 3 {
			role, rg, method := msg.Values[0], msg.Values[1], msg.Values[2]
			if role == "" || rg == "" || method == "" {
				m.form.SetError("role, resource group, and method are required")

				return m, nil
			}
			m.mode = policiesModeList

			return m, func() tea.Msg {
				return FormCreatePolicyMsg{
					Role:          role,
					ResourceGroup: rg,
					Method:        method,
				}
			}
		}
		m.mode = policiesModeList

	case FormCancelledMsg:
		m.mode = policiesModeList

	default:
		var cmd tea.Cmd
		m.form, cmd = m.form.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m PoliciesModel) renderList() string {
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

// FormDeletePolicyMsg requests an async policy delete in app.go.
type FormDeletePolicyMsg struct{ Policy enclave.Policy }

// FormCreatePolicyMsg requests an async policy create in app.go.
type FormCreatePolicyMsg struct {
	Role          string
	ResourceGroup string
	Method        string
}
