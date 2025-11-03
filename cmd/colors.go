package cmd

import "github.com/charmbracelet/lipgloss"

const (
	ColorPrimary       lipgloss.Color = "2"
	ColorTextHighlight lipgloss.Color = "3"
)

var (
	TextPrimary   = lipgloss.NewStyle().Foreground(ColorPrimary)
	TextHighlight = lipgloss.NewStyle().
			Foreground(ColorTextHighlight)
)
