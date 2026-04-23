package views

import (
	"cli/internal/styles"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FormField defines a single input field in a form.
type FormField struct {
	Label       string
	Placeholder string
	Secret      bool
}

// FormModel is a full-screen multi-field input form.
type FormModel struct {
	title  string
	fields []FormField
	inputs []textinput.Model
	cursor int
	err    string
	width  int
	height int
}

// NewForm builds a form with the given title and fields.
func NewForm(title string, fields []FormField) FormModel {
	inputs := make([]textinput.Model, len(fields))
	for i, f := range fields {
		ti := textinput.New()
		ti.Placeholder = f.Placeholder
		if f.Secret {
			ti.EchoMode = textinput.EchoPassword
		}
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
	}

	return FormModel{
		title:  title,
		fields: fields,
		inputs: inputs,
	}
}

// SetSize updates the rendering dimensions.
func (m *FormModel) SetSize(w, h int) { m.width = w; m.height = h }

// SetError sets a validation/API error message shown at the bottom.
func (m *FormModel) SetError(err string) { m.err = err }

// Update handles key events for form navigation and submission.
func (m FormModel) Update(msg tea.Msg) (FormModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc":
			return m, func() tea.Msg { return FormCancelledMsg{} }

		case "tab", "down":
			m.inputs[m.cursor].Blur()
			m.cursor = (m.cursor + 1) % len(m.inputs)
			m.inputs[m.cursor].Focus()

		case "shift+tab", "up":
			m.inputs[m.cursor].Blur()
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.inputs) - 1
			}
			m.inputs[m.cursor].Focus()

		case "enter":
			if m.cursor == len(m.inputs)-1 {
				return m, m.submit()
			}
			m.inputs[m.cursor].Blur()
			m.cursor++
			m.inputs[m.cursor].Focus()

		case "ctrl+s":
			return m, m.submit()
		}
	}

	// Delegate remaining input to focused textinput.
	var cmd tea.Cmd
	m.inputs[m.cursor], cmd = m.inputs[m.cursor].Update(msg)

	return m, cmd
}

func (m FormModel) submit() tea.Cmd {
	vals := make([]string, len(m.inputs))
	for i, in := range m.inputs {
		vals[i] = in.Value()
	}

	return func() tea.Msg { return FormSubmittedMsg{Values: vals} }
}

// View renders the form centered on screen.
func (m FormModel) View() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSlateLight).
		Width(20)
	activeLabel := lipgloss.NewStyle().
		Foreground(styles.ColorPrimaryGreen).
		Bold(true).
		Width(20)

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render(m.title) + "\n\n")

	for i, f := range m.fields {
		label := f.Label + ":"
		if i == m.cursor {
			b.WriteString(activeLabel.Render(label) + "  " + m.inputs[i].View() + "\n")
		} else {
			b.WriteString(labelStyle.Render(label) + "  " + m.inputs[i].View() + "\n")
		}
	}

	b.WriteString("\n")
	hint := styles.HelpKeyStyle.Render("ctrl+s") +
		lipgloss.NewStyle().Foreground(styles.ColorSlateDark).Render(" submit   ") +
		styles.HelpKeyStyle.Render("tab") +
		lipgloss.NewStyle().Foreground(styles.ColorSlateDark).Render(" next field   ") +
		styles.HelpKeyStyle.Render("esc") +
		lipgloss.NewStyle().Foreground(styles.ColorSlateDark).Render(" cancel")
	b.WriteString(hint + "\n")

	if m.err != "" {
		b.WriteString("\n" + styles.ErrorStyle.Render("  "+m.err) + "\n")
	}

	return b.String()
}
