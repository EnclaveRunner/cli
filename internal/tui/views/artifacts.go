package views

import (
	"context"
	"strings"

	"cli/internal/styles"

	"github.com/EnclaveRunner/sdk-go/enclave"
	tea "github.com/charmbracelet/bubbletea"
	"charm.land/lipgloss/v2"
)

// ArtifactsLoadedMsg carries loaded artifacts.
type ArtifactsLoadedMsg struct {
	Artifacts []enclave.Artifact
	Level     int    // 0=namespaces, 1=artifacts, 2=versions
	Namespace string // set at level>=1
	Name      string // set at level==2
	Err       error
}

// ArtifactsModel is the artifacts drill-down view.
type ArtifactsModel struct {
	Items     []enclave.Artifact
	Cursor    int
	Loading   bool
	Err       error
	level     int    // 0=namespaces, 1=artifacts, 2=versions
	namespace string // active namespace
	name      string // active artifact name
	width     int
	height    int
}

// Load fetches namespaces (level 0).
func (m ArtifactsModel) Load(c *enclave.Client) tea.Cmd {
	return func() tea.Msg {
		items, err := enclave.Collect(c.ListArtifactNamespaces(context.Background()))
		return ArtifactsLoadedMsg{Artifacts: items, Level: 0, Err: err}
	}
}

// SetSize updates the rendering area.
func (m *ArtifactsModel) SetSize(w, h int) { m.width = w; m.height = h }

// Update handles messages. Requires client for drill-down navigation.
func (m ArtifactsModel) Update(msg tea.Msg, c *enclave.Client) (ArtifactsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ArtifactsLoadedMsg:
		m.Loading = false
		m.Err = msg.Err
		m.Items = msg.Artifacts
		m.level = msg.Level
		m.namespace = msg.Namespace
		m.name = msg.Name
		m.Cursor = 0

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}
		case "enter":
			if len(m.Items) == 0 || m.Cursor >= len(m.Items) {
				break
			}
			item := m.Items[m.Cursor]
			m.Loading = true
			if m.level == 0 {
				ns := item.Namespace
				return m, func() tea.Msg {
					arts, err := enclave.Collect(c.ListArtifacts(context.Background(), ns))
					return ArtifactsLoadedMsg{Artifacts: arts, Level: 1, Namespace: ns, Err: err}
				}
			} else if m.level == 1 {
				ns, name := item.Namespace, item.Name
				return m, func() tea.Msg {
					vers, err := enclave.Collect(c.ListArtifactVersions(context.Background(), ns, name))
					return ArtifactsLoadedMsg{Artifacts: vers, Level: 2, Namespace: ns, Name: name, Err: err}
				}
			}
		case "esc":
			if m.level > 0 {
				m.Loading = true
				lvl := m.level
				ns := m.namespace
				return m, func() tea.Msg {
					if lvl == 2 {
						arts, err := enclave.Collect(c.ListArtifacts(context.Background(), ns))
						return ArtifactsLoadedMsg{Artifacts: arts, Level: 1, Namespace: ns, Err: err}
					}
					items, err := enclave.Collect(c.ListArtifactNamespaces(context.Background()))
					return ArtifactsLoadedMsg{Artifacts: items, Level: 0, Err: err}
				}
			}
		}
	}
	return m, nil
}

// View renders the artifacts drill-down table.
func (m ArtifactsModel) View() string {
	title := "Namespaces"
	if m.level == 1 {
		title = "Artifacts › " + m.namespace
	} else if m.level == 2 {
		title = "Versions › " + m.namespace + "/" + m.name
	}

	if m.Loading {
		return styles.MutedStyle.Render("\n  Loading " + strings.ToLower(title) + "…")
	}
	if m.Err != nil {
		return styles.ErrorStyle.Render("\n  Error: " + m.Err.Error())
	}

	var b strings.Builder
	b.WriteString(styles.TitleStyle.Render(title) + "\n")

	if len(m.Items) == 0 {
		b.WriteString(styles.MutedStyle.Render("\n  No items found."))
		return b.String()
	}

	if m.level == 0 {
		header := styles.HeaderStyle.Render(padRight("NAMESPACE", 30))
		b.WriteString(header + "\n")
		seen := map[string]bool{}
		idx := 0
		for _, a := range m.Items {
			if seen[a.Namespace] {
				continue
			}
			seen[a.Namespace] = true
			style := lipgloss.NewStyle().Padding(0, 1)
			if idx == m.Cursor {
				style = style.Background(styles.ColorSecondaryGreen).Foreground(styles.ColorNearBlack)
			}
			b.WriteString(style.Render(padRight(a.Namespace, 30)) + "\n")
			idx++
		}
	} else if m.level == 1 {
		header := styles.HeaderStyle.Render(padRight("NAME", 30))
		b.WriteString(header + "\n")
		for i, a := range m.Items {
			style := lipgloss.NewStyle().Padding(0, 1)
			if i == m.Cursor {
				style = style.Background(styles.ColorSecondaryGreen).Foreground(styles.ColorNearBlack)
			}
			b.WriteString(style.Render(padRight(a.Name, 30)) + "\n")
		}
	} else {
		headers := []string{"HASH", "TAGS", "CREATED", "PULLS"}
		b.WriteString(strings.Join([]string{
			styles.HeaderStyle.Render(padRight(headers[0], 16)),
			styles.HeaderStyle.Render(padRight(headers[1], 20)),
			styles.HeaderStyle.Render(padRight(headers[2], 12)),
			styles.HeaderStyle.Render(padRight(headers[3], 5)),
		}, " ") + "\n")
		for i, a := range m.Items {
			h := a.VersionHash
			if len(h) > 16 {
				h = h[:16]
			}
			tags := strings.Join(a.Tags, ", ")
			created := a.CreatedAt.Format("2006-01-02")
			style := lipgloss.NewStyle().Padding(0, 1)
			if i == m.Cursor {
				style = style.Background(styles.ColorSecondaryGreen).Foreground(styles.ColorNearBlack)
			}
			b.WriteString(strings.Join([]string{
				style.Render(padRight(h, 16)),
				style.Render(padRight(tags, 20)),
				style.Render(padRight(created, 12)),
				style.Render(padRight(strings.TrimSpace(""), 5)),
			}, " ") + "\n")
		}
	}

	return b.String()
}
