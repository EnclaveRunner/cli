package tui

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
