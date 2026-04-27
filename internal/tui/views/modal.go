package views

import (
	"cli/internal/styles"
	"strings"

	"charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
)

// ModalModel is a centered confirmation dialog overlay.
type ModalModel struct {
	message string
	width   int
	height  int
}

// NewModal creates a modal with the given confirmation message.
func NewModal(message string) ModalModel {
	return ModalModel{message: message}
}

// SetSize updates the terminal dimensions used for centering.
func (m *ModalModel) SetSize(w, h int) { m.width = w; m.height = h }

// Update handles modal key presses. Returns ModalConfirmedMsg or
// ModalCancelledMsg.
func (m ModalModel) Update(msg tea.Msg) (ModalModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "y", "Y":
			return m, func() tea.Msg { return ModalConfirmedMsg{} }
		case "n", "N", "esc":
			return m, func() tea.Msg { return ModalCancelledMsg{} }
		}
	}

	return m, nil
}

// View renders the modal as a standalone string (to be overlaid by the parent).
func (m ModalModel) View() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimaryGreen).
		Padding(1, 3)

	hint := styles.HelpKeyStyle.Render("[y]") +
		lipgloss.NewStyle().Foreground(styles.ColorSlateDark).Render(" Yes    ") +
		styles.HelpKeyStyle.Render("[n/Esc]") +
		lipgloss.NewStyle().Foreground(styles.ColorSlateDark).Render(" No")

	content := m.message + "\n\n" + hint
	box := boxStyle.Render(content)
	boxLines := strings.Split(box, "\n")
	boxH := len(boxLines)
	boxW := 0
	for _, l := range boxLines {
		w := len([]rune(stripANSI(l)))
		if w > boxW {
			boxW = w
		}
	}

	topPad := (m.height - boxH) / 2
	if topPad < 0 {
		topPad = 0
	}
	leftPad := (m.width - boxW) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	leftStr := strings.Repeat(" ", leftPad)

	var b strings.Builder
	for range topPad {
		b.WriteString("\n")
	}
	for _, line := range boxLines {
		b.WriteString(leftStr + line + "\n")
	}

	return b.String()
}
