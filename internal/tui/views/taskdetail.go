package views

import (
	"context"
	"fmt"
	"strings"

	"cli/internal/styles"

	"github.com/EnclaveRunner/sdk-go/enclave"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"charm.land/lipgloss/v2"
)

// TaskLogsLoadedMsg carries loaded task logs.
type TaskLogsLoadedMsg struct {
	Logs []enclave.TaskLog
	Err  error
}

// TaskDetailModel shows a single task's details and logs.
type TaskDetailModel struct {
	task    enclave.Task
	logs    []enclave.TaskLog
	loading bool
	err     error
	vp      viewport.Model
	width   int
	height  int
}

// SetTask sets the task to display and starts loading logs.
func (m TaskDetailModel) SetTask(t enclave.Task, c *enclave.Client) (TaskDetailModel, tea.Cmd) {
	m.task = t
	m.loading = true
	m.logs = nil
	m.err = nil
	return m, func() tea.Msg {
		logs, err := c.GetTaskLogs(context.Background(), t.ID)
		return TaskLogsLoadedMsg{Logs: logs, Err: err}
	}
}

// SetSize updates the viewport size.
func (m *TaskDetailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.vp = viewport.New(w, maxInt(1, h-8))
	m.vp.SetContent(m.renderLogs())
}

// Update handles messages for the task detail view.
func (m TaskDetailModel) Update(msg tea.Msg) (TaskDetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case TaskLogsLoadedMsg:
		m.loading = false
		m.err = msg.Err
		m.logs = msg.Logs
		m.vp.SetContent(m.renderLogs())
		return m, nil
	default:
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd
	}
}

// View renders the task detail pane.
func (m TaskDetailModel) View() string {
	t := m.task

	field := func(label, value string) string {
		return styles.MutedStyle.Render(label+": ") + value
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render("Task "+t.ID) + "\n\n")
	b.WriteString(field("Source", t.Source) + "\n")
	b.WriteString(field("State", styles.TaskStateBadge(t.Status.State)) + "\n")
	b.WriteString(field("Retries", fmt.Sprintf("%d / %d", t.Status.Retries, t.Retries)) + "\n")
	if t.Status.LastError != "" {
		b.WriteString(field("Last Error", styles.ErrorStyle.Render(t.Status.LastError)) + "\n")
	}
	if !t.Status.CompletedAt.IsZero() {
		b.WriteString(field("Completed", t.Status.CompletedAt.String()) + "\n")
	}
	b.WriteString("\n" + styles.TitleStyle.Render("Logs") + "\n")

	if m.loading {
		b.WriteString(styles.MutedStyle.Render("  Loading logs…\n"))
	} else {
		b.WriteString(m.vp.View())
	}

	return b.String()
}

func (m TaskDetailModel) renderLogs() string {
	if m.err != nil {
		return styles.ErrorStyle.Render("Error loading logs: " + m.err.Error())
	}
	if len(m.logs) == 0 {
		return styles.MutedStyle.Render("No logs.")
	}
	var b strings.Builder
	for _, l := range m.logs {
		ts := l.Timestamp.Format("15:04:05.000")
		level := padRight(l.Level, 5)
		b.WriteString(
			styles.MutedStyle.Render(ts+" ") +
				logLevelStyle(l.Level).Render(level+" ") +
				styles.MutedStyle.Render("["+l.Issuer+"] ") +
				l.Message + "\n",
		)
	}
	return b.String()
}

func logLevelStyle(level string) lipgloss.Style {
	switch strings.ToLower(level) {
	case "error", "fatal":
		return styles.ErrorStyle
	case "warn", "warning":
		return lipgloss.NewStyle().Foreground(styles.ColorWarmHighlight)
	default:
		return styles.MutedStyle
	}
}
