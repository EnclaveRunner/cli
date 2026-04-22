package views //nolint:dupl // Bubbletea view models follow identical structure by design.

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

// ResourceGroupsModel is the resource groups view.
//

type ResourceGroupsModel struct {
	RGs       []enclave.ResourceGroup
	Cursor    int
	Loading   bool
	Err       error
	colOffset int
	width     int
	height    int
}

// Load fetches all resource groups.
func (m ResourceGroupsModel) Load( //nolint:gocritic // hugeParam: Bubbletea requires value receiver.
	c *enclave.Client,
) tea.Cmd {
	return func() tea.Msg {
		rgs, err := enclave.Collect(c.ListResourceGroups(context.Background()))

		return ResourceGroupsLoadedMsg{RGs: rgs, Err: err}
	}
}

// SetSize updates the rendering area.
func (m *ResourceGroupsModel) SetSize(w, h int) { m.width = w; m.height = h }

// Update handles messages.
func (m ResourceGroupsModel) Update( //nolint:gocritic // hugeParam: Bubbletea requires value receiver.
	msg tea.Msg,
) (ResourceGroupsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ResourceGroupsLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.RGs = msg.RGs
		m.Cursor = 0
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
		}
	}

	return m, nil
}

// View renders the resource groups table.
func (m ResourceGroupsModel) View() string { //nolint:gocritic // hugeParam: Bubbletea requires value receiver.
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
