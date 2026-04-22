package styles

import "charm.land/lipgloss/v2"

var (
	// HeaderStyle is used for table column headers.
	HeaderStyle = lipgloss.NewStyle().
			Foreground(ColorNearBlack).
			Background(ColorPrimaryGreen).
			Bold(true).
			Padding(0, 1)

	// SelectedRowStyle highlights the cursor row in TUI tables.
	SelectedRowStyle = lipgloss.NewStyle().
				Foreground(ColorNearBlack).
				Background(ColorSecondaryGreen)

	// MutedStyle renders secondary/contextual text.
	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorSlateDark)

	// TitleStyle is used for view titles in the TUI.
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorPrimaryGreen).
			Bold(true)

	// StatusBarStyle is the top status bar background.
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorNearBlack).
			Background(ColorDarkestGreen).
			Padding(0, 1)

	// StatusBarHighlight is used for active view name in the status bar.
	StatusBarHighlight = lipgloss.NewStyle().
				Foreground(ColorNearBlack).
				Background(ColorPrimaryGreen).
				Bold(true).
				Padding(0, 1)

	// HelpBarStyle is the bottom help bar.
	HelpBarStyle = lipgloss.NewStyle().
			Foreground(ColorSlateDark).
			Background(ColorNearBlack).
			Padding(0, 1)

	// HelpKeyStyle highlights keybinding keys.
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(ColorPrimaryGreen)

	// ErrorStyle renders error messages.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorWarmHighlight)

	// BorderStyle is used for panel borders.
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDarkGreen)
)

// TaskStateBadge returns a coloured badge string for the given task state.
func TaskStateBadge(state string) string {
	switch state {
	case "running", "processing":
		return lipgloss.NewStyle().Foreground(ColorPrimaryGreen).Render(IconRunning + " " + state)
	case "failed", "error":
		return lipgloss.NewStyle().Foreground(ColorWarmHighlight).Render(IconFailed + " " + state)
	case "completed", "done":
		return lipgloss.NewStyle().Foreground(ColorLogoTeal).Render(IconDone + " " + state)
	default:
		return lipgloss.NewStyle().Foreground(ColorSlateLight).Render(IconPending + " " + state)
	}
}
