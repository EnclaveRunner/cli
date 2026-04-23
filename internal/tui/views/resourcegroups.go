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

// ResourceGroupsLoadedMsg carries loaded resource groups.
type ResourceGroupsLoadedMsg struct {
	RGs []enclave.ResourceGroup
	Err error
}

type rgMode int

const (
	rgModeList     rgMode = iota
	rgModeDescribe        // enter: full detail with endpoints
	rgModeModal           // d: confirm delete
	rgModeForm            // c: create form
)

// ResourceGroupsModel is the resource groups view.
type ResourceGroupsModel struct {
	RGs       []enclave.ResourceGroup
	Cursor    int
	Loading   bool
	Err       error
	colOffset int
	width     int
	height    int

	mode  rgMode
	modal ModalModel
	form  FormModel
}

// Load fetches all resource groups.
func (m ResourceGroupsModel) Load(
	c *enclave.Client,
) tea.Cmd {
	return func() tea.Msg {
		rgs, err := enclave.Collect(c.ListResourceGroups(context.Background()))

		return ResourceGroupsLoadedMsg{RGs: rgs, Err: err}
	}
}

// SetSize updates the rendering area.
func (m *ResourceGroupsModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.modal.SetSize(w, h)
	m.form.SetSize(w, h)
}

// IsCapturing reports whether the view is in a mode that owns the keyboard.
func (m ResourceGroupsModel) IsCapturing() bool {
	return m.mode != rgModeList
}

func (m ResourceGroupsModel) selectedRG() (enclave.ResourceGroup, bool) {
	if len(m.RGs) == 0 || m.Cursor >= len(m.RGs) {
		return enclave.ResourceGroup{}, false
	}

	return m.RGs[m.Cursor], true
}

// Update handles messages.
func (m ResourceGroupsModel) Update(
	msg tea.Msg,
) (ResourceGroupsModel, tea.Cmd) {
	switch m.mode {
	case rgModeModal:
		return m.updateModal(msg)
	case rgModeForm:
		return m.updateForm(msg)
	case rgModeDescribe:
		return m.updateDescribe(msg)
	}

	return m.updateList(msg)
}

func (m ResourceGroupsModel) updateList(msg tea.Msg) (ResourceGroupsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ResourceGroupsLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.RGs = msg.RGs
		m.Cursor = 0

	case ResourceGroupDeletedMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
		}

	case ResourceGroupCreatedMsg:
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
			if m.Cursor < len(m.RGs)-1 {
				m.Cursor++
			}
		case keyLeft:
			if m.colOffset > 0 {
				m.colOffset--
			}
		case keyRight:
			m.colOffset++

		case "enter":
			if _, ok := m.selectedRG(); ok {
				m.mode = rgModeDescribe
			}

		case "d":
			if rg, ok := m.selectedRG(); ok {
				m.modal = NewModal("Delete resource group \"" + rg.Name + "\"?")
				m.modal.SetSize(m.width, m.height)
				m.mode = rgModeModal
			}

		case "c":
			m.form = NewForm("Create Resource Group", []FormField{
				{Label: "Name", Placeholder: "my-api"},
				{Label: "Endpoints (comma-sep)", Placeholder: "/api/v1/*, /health"},
			})
			m.form.SetSize(m.width, m.height)
			m.mode = rgModeForm
		}
	}

	return m, nil
}

func (m ResourceGroupsModel) updateModal(msg tea.Msg) (ResourceGroupsModel, tea.Cmd) {
	switch msg.(type) {
	case ModalConfirmedMsg:
		if rg, ok := m.selectedRG(); ok {
			m.mode = rgModeList
			m.Loading = true
			name := rg.Name

			return m, func() tea.Msg { return FormDeleteRGMsg{Name: name} }
		}
		m.mode = rgModeList

	case ModalCancelledMsg:
		m.mode = rgModeList

	default:
		var cmd tea.Cmd
		m.modal, cmd = m.modal.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m ResourceGroupsModel) updateForm(msg tea.Msg) (ResourceGroupsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case FormSubmittedMsg:
		if len(msg.Values) >= 1 {
			name := msg.Values[0]
			if name == "" {
				m.form.SetError("name is required")

				return m, nil
			}
			endpointsRaw := ""
			if len(msg.Values) >= 2 {
				endpointsRaw = msg.Values[1]
			}
			m.mode = rgModeList

			return m, func() tea.Msg { return FormCreateRGMsg{Name: name, EndpointsRaw: endpointsRaw} }
		}
		m.mode = rgModeList

	case FormCancelledMsg:
		m.mode = rgModeList

	default:
		var cmd tea.Cmd
		m.form, cmd = m.form.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m ResourceGroupsModel) updateDescribe(msg tea.Msg) (ResourceGroupsModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc", "q":
			m.mode = rgModeList
		}
	}

	return m, nil
}

// View renders the resource groups table or the active overlay.
func (m ResourceGroupsModel) View() string {
	switch m.mode {
	case rgModeDescribe:
		return m.renderDescribe()
	case rgModeModal:
		return m.renderList() + m.modal.View()
	case rgModeForm:
		return m.form.View()
	}

	return m.renderList()
}

func (m ResourceGroupsModel) renderList() string {
	if m.Loading {
		return styles.MutedStyle.Render("\n  Loading resource groups…")
	}
	if m.Err != nil {
		return styles.ErrorStyle.Render("\n  Error: " + m.Err.Error())
	}
	if len(m.RGs) == 0 {
		return styles.MutedStyle.Render("\n  No resource groups found.")
	}

	headers := []string{"NAME", "ENDPOINTS"}
	colWidths := []int{len(headers[0]), len(headers[1])}

	rows := make([][]string, len(m.RGs))
	for i, rg := range m.RGs {
		count := strconv.Itoa(len(rg.Endpoints))
		rows[i] = []string{rg.Name, count}
		updateWidth(&colWidths[0], len(rg.Name))
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

func (m ResourceGroupsModel) renderDescribe() string {
	rg, ok := m.selectedRG()
	if !ok {
		return styles.MutedStyle.Render("\n  No resource group selected.")
	}

	field := func(label, value string) string {
		return styles.MutedStyle.Render(label+": ") + value
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render("Resource Group "+rg.Name) + "\n\n")
	b.WriteString(field("Name", rg.Name) + "\n")
	b.WriteString(field("Endpoint count", strconv.Itoa(len(rg.Endpoints))) + "\n")

	if len(rg.Endpoints) > 0 {
		b.WriteString("\n" + styles.TitleStyle.Render("Endpoints") + "\n")
		for _, ep := range rg.Endpoints {
			b.WriteString("  " + ep + "\n")
		}
	}

	b.WriteString("\n" + styles.HelpKeyStyle.Render("esc") +
		lipgloss.NewStyle().Foreground(styles.ColorSlateDark).Render(" back"))

	return b.String()
}

// FormDeleteRGMsg requests an async resource group delete in app.go.
type FormDeleteRGMsg struct{ Name string }

// FormCreateRGMsg requests an async resource group create in app.go.
type FormCreateRGMsg struct {
	Name         string
	EndpointsRaw string // comma-separated endpoint list
}
