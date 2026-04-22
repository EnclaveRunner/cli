package views

import (
	"cli/internal/styles"
	"context"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/EnclaveRunner/sdk-go/enclave"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	keyUp    = "up"
	keyDown  = "down"
	keyLeft  = "left"
	keyRight = "right"
	keyK     = "k"
	keyJ     = "j"
)

// TasksLoadedMsg is returned by the Load command.
type TasksLoadedMsg struct {
	Tasks []enclave.Task
	Err   error
}

// TasksModel is the task list view (TUI home screen).
type TasksModel struct {
	Tasks     []enclave.Task
	Cursor    int
	Loading   bool
	Err       error
	colOffset int
	width     int
	height    int
}

// Load fetches all tasks asynchronously.
func (m TasksModel) Load( //nolint:gocritic // hugeParam: Bubbletea value receiver.
	c *enclave.Client,
) tea.Cmd {
	return func() tea.Msg {
		tasks, err := enclave.Collect(c.ListTasks(context.Background()))

		return TasksLoadedMsg{Tasks: tasks, Err: err}
	}
}

// SetSize sets the available rendering area.
func (m *TasksModel) SetSize(w, h int) { m.width = w; m.height = h }

// SelectedTask returns the currently highlighted task, or zero value.
func (m TasksModel) SelectedTask() (enclave.Task, bool) { //nolint:gocritic // hugeParam: Bubbletea requires value receiver.
	if len(m.Tasks) == 0 || m.Cursor >= len(m.Tasks) {
		return enclave.Task{}, false
	}

	return m.Tasks[m.Cursor], true
}

// Update handles messages for the tasks view.
func (m TasksModel) Update( //nolint:gocritic // hugeParam: Bubbletea value receiver.
	msg tea.Msg,
) (TasksModel, tea.Cmd) {
	switch msg := msg.(type) {
	case TasksLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.Tasks = msg.Tasks
		m.Cursor = 0

	case tea.KeyMsg:
		switch msg.String() {
		case keyUp, keyK:
			if m.Cursor > 0 {
				m.Cursor--
			}
		case keyDown, keyJ:
			if m.Cursor < len(m.Tasks)-1 {
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

// View renders the tasks table.
func (m TasksModel) View() string { //nolint:gocritic // hugeParam: Bubbletea requires value receiver.
	if m.Loading {
		return styles.MutedStyle.Render("\n  Loading tasks…")
	}
	if m.Err != nil {
		return styles.ErrorStyle.Render("\n  Error: " + m.Err.Error())
	}
	if len(m.Tasks) == 0 {
		return styles.MutedStyle.Render("\n  No tasks found.")
	}

	headers := []string{
		"ID",
		"SOURCE",
		"STATE",
		"RETRIES",
		"LAST ERROR",
		"NEXT PROCESS",
	}
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}

	rows := make([][]string, len(m.Tasks))
	for i := range m.Tasks {
		t := &m.Tasks[i]
		id := t.ID
		source := t.Source
		if len(source) > 30 {
			source = source[:29] + "…"
		}
		state := styles.TaskStateBadge(t.Status.State)
		statePlain := stripANSI(state)
		retries := strconv.Itoa(t.Status.Retries)
		lastErr := t.Status.LastError
		if len(lastErr) > 30 {
			lastErr = lastErr[:29] + "…"
		}
		next := "-"
		if !t.Status.NextProcessAt.IsZero() {
			next = t.Status.NextProcessAt.Format(time.RFC3339)
		}
		rows[i] = []string{id, source, state, retries, lastErr, next}

		updateWidth(&colWidths[0], len(id))
		updateWidth(&colWidths[1], len(source))
		updateWidth(&colWidths[2], len(statePlain))
		updateWidth(&colWidths[3], len(retries))
		updateWidth(&colWidths[4], len(lastErr))
		updateWidth(&colWidths[5], len(next))
	}

	var b strings.Builder

	// Apply column offset for horizontal scrolling.
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
			plain := stripANSI(cell)
			padding := maxInt(0, colWidths[j]-len([]rune(plain)))
			if i == m.Cursor {
				cells[j] = lipgloss.NewStyle().Padding(0, 1).
					Background(styles.ColorSecondaryGreen).
					Foreground(styles.ColorNearBlack).
					Render(plain + strings.Repeat(" ", padding))
			} else {
				cells[j] = " " + cell + strings.Repeat(" ", padding) + " "
			}
		}
		b.WriteString(strings.Join(cells[startCol:], "") + "\n")
	}

	return b.String()
}

func updateWidth(w *int, n int) {
	if n > *w {
		*w = n
	}
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}

	return s + strings.Repeat(" ", n-len(s))
}

func stripANSI(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true

			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}

			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}
