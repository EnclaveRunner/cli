package tui

import (
	"cli/internal/styles"
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// logoSeg is a colored text segment.
type logoSeg struct {
	text  string
	style lipgloss.Style
}

var (
	styleLogoHi  = lipgloss.NewStyle().Foreground(styles.ColorPrimaryGreen)
	styleLogoLo  = lipgloss.NewStyle().Foreground(styles.ColorDarkGreen)
	styleLogoDim = lipgloss.NewStyle().Foreground(styles.ColorSlateDark)
)

// logoArt defines each line as a slice of colored segments.
// Spells "ENCL" in a compact ASCII font.
// Every line renders to exactly 23 visible characters.
//
//	 ____  _  _  ___  __
//	( ___)( \( )/ __)(  )
//	 )__)  )  (( (__  )(__
//	(____)(_)\_)\___)(____)
var logoArt = [][]logoSeg{
	// line 1: blank top spacer
	{{`                       `, styleLogoDim}},
	// line 2: blank
	{{`                       `, styleLogoDim}},
	// line 3
	{
		{` `, styleLogoDim},
		{`____  _  _  ___  __   `, styleLogoHi},
	},
	// line 4
	{
		{`( `, styleLogoLo},
		{`___)`, styleLogoHi},
		{`( \( )/ __)(  )  `, styleLogoHi},
	},
	// line 5
	{
		{` `, styleLogoDim},
		{`  )__)  )  (( (__  )`, styleLogoHi},
		{`(__  `, styleLogoLo},
	},
	// line 6
	{
		{`(`, styleLogoLo},
		{`____)`, styleLogoHi},
		{`(_)\_)\___)`, styleLogoHi},
		{`(____)`, styleLogoLo},
	},
	// line 7: blank bottom spacer
	{{`                       `, styleLogoDim}},
	// line 8: blank
	{{`                       `, styleLogoDim}},
	// line 9: blank
	{{`                       `, styleLogoDim}},
}

// logoWidth is the visible width of each logo line.
const logoWidth = 23

// renderLogoLine renders one logo line (slice of segments) as a single string.
func renderLogoLine(segs []logoSeg) string {
	var sb strings.Builder
	for i := range segs {
		sb.WriteString(segs[i].style.Render(segs[i].text))
	}

	return sb.String()
}

// headerContentRows must match len(logoArt).
const headerContentRows = 9

// headerPanel renders the three-column info panel at the top.
type headerPanel struct {
	apiURL       string
	username     string
	version      string
	updateNotice string
	width        int
}

func newHeaderPanel(apiURL, username, version string) headerPanel {
	return headerPanel{apiURL: apiURL, username: username, version: version}
}

// Height returns the fixed number of lines the header panel occupies
// (top border + content rows + bottom border).
func (h headerPanel) Height() int { return 2 + headerContentRows }

func (h headerPanel) View() string {
	// Reserve space: logo col + 2 separators + margins
	// " " + col1 + " │ " + col2 + " │ " + col3 + " "
	// margins: 1+1+1+1+1+1 = 6 chars, separators: 2 × "│" = 2 chars → total
	// overhead = 8
	logoColW := logoWidth + 2 // +2 for side margins inside column
	overhead := 8
	infoColW := 36
	bindColW := h.width - infoColW - logoColW - overhead
	if bindColW < 20 {
		bindColW = 20
	}

	// Column 1: connection info (9 lines, last 5 blank).
	infoLines := []string{
		styles.MutedStyle.Render(
			"server  ",
		) + lipgloss.NewStyle().
			Foreground(styles.ColorLogoTeal).
			Render(truncate(h.apiURL, infoColW-9)),
		styles.MutedStyle.Render(
			"user    ",
		) + lipgloss.NewStyle().
			Foreground(styles.ColorPrimaryGreen).
			Render(h.username),
		styles.MutedStyle.Render(
			"version ",
		) + lipgloss.NewStyle().
			Foreground(styles.ColorSlateLight).
			Render(h.version),
		// If an update notice is present, show it below the version line.
		styles.MutedStyle.Render("          ") + lipgloss.NewStyle().
			Foreground(styles.ColorPrimaryGreen).
			Render(h.updateNotice),
		"",
		"",
		"",
		"",
		"",
	}

	// Column 2: keybindings.
	kb := func(key, desc string) string {
		return styles.HelpKeyStyle.Render(key) +
			lipgloss.NewStyle().Foreground(styles.ColorSlateDark).Render(" "+desc)
	}
	bindLines := []string{
		kb("1-6", "switch view") + "  " + kb("↑↓/jk", "navigate"),
		kb("enter", "select    ") + "  " + kb("esc", "back"),
		kb("←→", "scroll cols") + "  " + kb("r", "refresh"),
		kb("c", "create     ") + "  " + kb("d", "delete"),
		kb("q", "quit"),
		"",
		"",
		"",
		"",
	}

	sep := lipgloss.NewStyle().Foreground(styles.ColorDarkGreen).Render("│")
	borderLine := lipgloss.NewStyle().
		Foreground(styles.ColorDarkGreen).
		Render(strings.Repeat("─", h.width))

	var b strings.Builder
	b.WriteString(borderLine + "\n")
	for i := range headerContentRows {
		l1 := padTo(safeGetStr(infoLines, i), infoColW)

		rawBind := safeGetStr(bindLines, i)
		l2 := padTo(rawBind, bindColW)

		// Logo: right-align within its column (pad left so it hugs the right
		// border).
		styledLogo := ""
		if i < len(logoArt) {
			styledLogo = renderLogoLine(logoArt[i])
		}
		plainLogo := stripANSI(styledLogo)
		leftPad := logoColW - len([]rune(plainLogo))
		if leftPad < 0 {
			leftPad = 0
		}
		l3 := strings.Repeat(" ", leftPad) + styledLogo

		b.WriteString(" " + l1 + " " + sep + " " + l2 + " " + sep + l3 + "\n")
	}
	b.WriteString(borderLine)

	return b.String()
}

func (h *headerPanel) setWidth(w int) { h.width = w }

// tabRibbon renders the tab bar for view switching.
type tabRibbon struct {
	activeView View
	width      int
}

func newTabRibbon() tabRibbon { return tabRibbon{} }

// navigableTabs is the ordered list of switchable views.
var navigableTabs = []View{
	ViewTasks,
	ViewUsers,
	ViewRoles,
	ViewResourceGroups,
	ViewPolicies,
	ViewArtifacts,
}

var tabLabels = map[View]string{
	ViewTasks:          "Tasks",
	ViewUsers:          "Users",
	ViewRoles:          "Roles",
	ViewResourceGroups: "RGroups",
	ViewPolicies:       "Policies",
	ViewArtifacts:      "Artifacts",
}

func (t tabRibbon) View() string {
	var parts []string
	for i, v := range navigableTabs {
		label := fmt.Sprintf("%d %s", i+1, tabLabels[v])
		if v == t.activeView {
			parts = append(parts, lipgloss.NewStyle().
				Foreground(styles.ColorNearBlack).
				Background(styles.ColorPrimaryGreen).
				Bold(true).
				Padding(0, 1).
				Render(label))
		} else {
			parts = append(parts, lipgloss.NewStyle().
				Foreground(styles.ColorSlateLight).
				Background(styles.ColorNearBlack).
				Padding(0, 1).
				Render(label))
		}
	}
	ribbon := strings.Join(parts, "")
	plain := stripANSI(ribbon)
	pad := t.width - len([]rune(plain))
	if pad > 0 {
		ribbon += lipgloss.NewStyle().
			Background(styles.ColorNearBlack).
			Render(strings.Repeat(" ", pad))
	}

	return ribbon
}

func (t *tabRibbon) setView(v View) { t.activeView = v }
func (t *tabRibbon) setWidth(w int) { t.width = w }

// --- helpers ---

// padTo pads a (possibly ANSI-coloured) string to exactly w visible runes.
func padTo(s string, w int) string {
	plain := stripANSI(s)
	pad := w - len([]rune(plain))
	if pad <= 0 {
		return s
	}

	return s + strings.Repeat(" ", pad)
}

func safeGetStr(lines []string, i int) string {
	if i < len(lines) {
		return lines[i]
	}

	return ""
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}

	return string(runes[:n-1]) + "…"
}
