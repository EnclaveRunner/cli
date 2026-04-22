package tui

import (
	"github.com/EnclaveRunner/sdk-go/enclave"
	tea "github.com/charmbracelet/bubbletea"
)

// View identifies the active TUI pane.
type View int

const (
	ViewTasks View = iota
	ViewUsers
	ViewRoles
	ViewResourceGroups
	ViewPolicies
	ViewArtifacts
	ViewTaskDetail
)

// viewName returns the display label for a view.
func viewName(v View) string {
	switch v {
	case ViewTasks:
		return "Tasks"
	case ViewUsers:
		return "Users"
	case ViewRoles:
		return "Roles"
	case ViewResourceGroups:
		return "Resource Groups"
	case ViewPolicies:
		return "Policies"
	case ViewArtifacts:
		return "Artifacts"
	case ViewTaskDetail:
		return "Task Detail"
	}
	return ""
}

// Ensure the SDK and tea packages are used.
var _ *enclave.Client
var _ tea.Cmd
