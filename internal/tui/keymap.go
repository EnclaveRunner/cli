package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

// keyMap holds all keybindings for the TUI.
type keyMap struct {
	Quit    key.Binding
	Refresh key.Binding
	Back    key.Binding
	Enter   key.Binding
	Up      key.Binding
	Down    key.Binding
	// View switches
	Tasks          key.Binding
	Users          key.Binding
	Roles          key.Binding
	ResourceGroups key.Binding
	Policies       key.Binding
	Artifacts      key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Tasks: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "tasks"),
	),
	Users: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "users"),
	),
	Roles: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "roles"),
	),
	ResourceGroups: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "rgs"),
	),
	Policies: key.NewBinding(
		key.WithKeys("5"),
		key.WithHelp("5", "policies"),
	),
	Artifacts: key.NewBinding(
		key.WithKeys("6"),
		key.WithHelp("6", "artifacts"),
	),
}
