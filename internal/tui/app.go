package tui

import (
	"cli/internal/styles"
	"cli/internal/tui/views"
	"context"
	"fmt"
	"strings"

	iv "cli/internal/version"

	"github.com/EnclaveRunner/sdk-go/enclave"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	minWidth  = 80
	minHeight = 20
)

// AppModel is the root Bubbletea model for the TUI.
type AppModel struct {
	client     *enclave.Client
	activeView View
	prevView   View

	tasks          views.TasksModel
	users          views.UsersModel
	roles          views.RolesModel
	resourceGroups views.ResourceGroupsModel
	policies       views.PoliciesModel
	artifacts      views.ArtifactsModel
	taskDetail     views.TaskDetailModel

	header headerPanel
	tabs   tabRibbon

	width  int
	height int
}

// versionCheckedMsg is sent when an asynchronous remote version check
// completes.
type versionCheckedMsg struct {
	Remote string
	Newer  bool
}

// checkVersionCmd fetches the remote version asynchronously.
func checkVersionCmd(local string) tea.Cmd {
	return func() tea.Msg {
		remote, newer, err := iv.CheckRemote(local)
		if err != nil {
			return nil
		}

		return versionCheckedMsg{Remote: remote, Newer: newer}
	}
}

// New creates a new TUI app model.
func New(c *enclave.Client, apiURL, username, version string) AppModel {
	m := AppModel{
		client:     c,
		activeView: ViewTasks,
		header:     newHeaderPanel(apiURL, username, version),
		tabs:       newTabRibbon(),
	}
	m.tabs.setView(ViewTasks)
	m.tasks.Loading = true

	return m
}

// Init loads initial data (tasks view).
func (m AppModel) Init() tea.Cmd {
	// Kick off tasks load and an async version check.
	return tea.Batch(m.tasks.Load(m.client), checkVersionCmd(m.header.version))
}

// Update is the main event loop.
func (m AppModel) Update(
	msg tea.Msg,
) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case versionCheckedMsg:
		if msg.Newer {
			m.header.updateNotice = fmt.Sprintf("⚡️ %s (latest)", msg.Remote)
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.setWidth(m.width)
		m.tabs.setWidth(m.width)
		contentH := m.height - m.header.Height() - 1 // 1 for tab ribbon
		if contentH < 1 {
			contentH = 1
		}
		m.tasks.SetSize(m.width, contentH)
		m.users.SetSize(m.width, contentH)
		m.roles.SetSize(m.width, contentH)
		m.resourceGroups.SetSize(m.width, contentH)
		m.policies.SetSize(m.width, contentH)
		m.artifacts.SetSize(m.width, contentH)
		m.taskDetail.SetSize(m.width, contentH)

		return m, nil

	// --- data loaded ---
	case views.TasksLoadedMsg:
		m.tasks, _ = m.tasks.Update(msg)
	case views.UsersLoadedMsg:
		m.users, _ = m.users.Update(msg)
	case views.RolesLoadedMsg:
		m.roles, _ = m.roles.Update(msg)
	case views.ResourceGroupsLoadedMsg:
		m.resourceGroups, _ = m.resourceGroups.Update(msg)
	case views.PoliciesLoadedMsg:
		m.policies, _ = m.policies.Update(msg)
	case views.ArtifactsLoadedMsg:
		m.artifacts, _ = m.artifacts.Update(msg, m.client)
	case views.TaskLogsLoadedMsg:
		m.taskDetail, _ = m.taskDetail.Update(msg)

	// --- form/modal lifecycle routed to active view ---
	case views.FormSubmittedMsg, views.FormCancelledMsg,
		views.ModalConfirmedMsg, views.ModalCancelledMsg:
		return m.delegateMsg(msg)

	// --- user operations ---
	case views.FormDeleteUserMsg:
		return m, m.deleteUserCmd(msg.Name)
	case views.FormCreateUserMsg:
		return m, m.createUserCmd(msg.Name, msg.Display, msg.Pass)
	case views.UserDeletedMsg:
		m.users, _ = m.users.Update(msg)
		if msg.Err == nil {
			m.users.Loading = true
			return m, m.users.Load(m.client)
		}
	case views.UserCreatedMsg:
		m.users, _ = m.users.Update(msg)
		if msg.Err == nil {
			m.users.Loading = true
			return m, m.users.Load(m.client)
		}

	// --- role operations ---
	case views.FormDeleteRoleMsg:
		return m, m.deleteRoleCmd(msg.Name)
	case views.FormCreateRoleMsg:
		return m, m.createRoleCmd(msg.Name, msg.UsersRaw)
	case views.RoleDeletedMsg:
		m.roles, _ = m.roles.Update(msg)
		if msg.Err == nil {
			m.roles.Loading = true
			return m, m.roles.Load(m.client)
		}
	case views.RoleCreatedMsg:
		m.roles, _ = m.roles.Update(msg)
		if msg.Err == nil {
			m.roles.Loading = true
			return m, m.roles.Load(m.client)
		}

	// --- resource group operations ---
	case views.FormDeleteRGMsg:
		return m, m.deleteRGCmd(msg.Name)
	case views.FormCreateRGMsg:
		return m, m.createRGCmd(msg.Name, msg.EndpointsRaw)
	case views.ResourceGroupDeletedMsg:
		m.resourceGroups, _ = m.resourceGroups.Update(msg)
		if msg.Err == nil {
			m.resourceGroups.Loading = true
			return m, m.resourceGroups.Load(m.client)
		}
	case views.ResourceGroupCreatedMsg:
		m.resourceGroups, _ = m.resourceGroups.Update(msg)
		if msg.Err == nil {
			m.resourceGroups.Loading = true
			return m, m.resourceGroups.Load(m.client)
		}

	// --- policy operations ---
	case views.FormDeletePolicyMsg:
		return m, m.deletePolicyCmd(msg.Policy)
	case views.FormCreatePolicyMsg:
		return m, m.createPolicyCmd(msg.Role, msg.ResourceGroup, msg.Method)
	case views.PolicyDeletedMsg:
		m.policies, _ = m.policies.Update(msg)
		if msg.Err == nil {
			m.policies.Loading = true
			return m, m.policies.Load(m.client)
		}
	case views.PolicyCreatedMsg:
		m.policies, _ = m.policies.Update(msg)
		if msg.Err == nil {
			m.policies.Loading = true
			return m, m.policies.Load(m.client)
		}

	case tea.KeyMsg:
		// Always allow quit.
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		// Ignore all other input when terminal is too small.
		if m.width < minWidth || m.height < minHeight {
			return m, nil
		}

		// When a sub-view is capturing input (form/modal/describe), delegate
		// everything directly so global hotkeys don't interfere.
		if m.isCapturing() {
			return m.delegateKey(msg)
		}

		// "q" only quits from top-level list views, not inside sub-views.
		if msg.String() == "q" {
			return m, tea.Quit
		}

		switch msg.String() {
		case "1":
			return m.switchToView(ViewTasks)
		case "2":
			return m.switchToView(ViewUsers)
		case "3":
			return m.switchToView(ViewRoles)
		case "4":
			return m.switchToView(ViewResourceGroups)
		case "5":
			return m.switchToView(ViewPolicies)
		case "6":
			return m.switchToView(ViewArtifacts)

		case "r":
			return m.doRefresh()

		case "esc":
			if m.activeView == ViewTaskDetail {
				m.activeView = m.prevView
				m.tabs.setView(m.activeView)

				return m, nil
			}
			if m.activeView == ViewArtifacts {
				var cmd tea.Cmd
				m.artifacts, cmd = m.artifacts.Update(msg, m.client)

				return m, cmd
			}

		case "enter":
			if m.activeView == ViewTasks {
				if t, ok := m.tasks.SelectedTask(); ok {
					m.prevView = m.activeView
					var cmd tea.Cmd
					m.taskDetail, cmd = m.taskDetail.SetTask(t, m.client)
					m.activeView = ViewTaskDetail
					m.tabs.setView(m.activeView)

					return m, cmd
				}
			}
			if m.activeView == ViewArtifacts {
				var cmd tea.Cmd
				m.artifacts, cmd = m.artifacts.Update(msg, m.client)

				return m, cmd
			}
		}

		return m.delegateKey(msg)
	}

	return m, nil
}

// View renders the full TUI.
func (m AppModel) View() string {
	if m.width < minWidth || m.height < minHeight {
		return m.tooSmallView()
	}

	return m.header.View() + "\n" + m.tabs.View() + "\n" + m.activeContent()
}

func (m AppModel) tooSmallView() string {
	msg := fmt.Sprintf(
		"Terminal too small (%dx%d). Minimum: %dx%d. Press q to quit.",
		m.width, m.height, minWidth, minHeight,
	)
	lines := []string{"", "  " + styles.MutedStyle.Render(msg)}

	return strings.Join(lines, "\n")
}

func (m AppModel) activeContent() string {
	switch m.activeView {
	case ViewTasks:
		return m.tasks.View()
	case ViewUsers:
		return m.users.View()
	case ViewRoles:
		return m.roles.View()
	case ViewResourceGroups:
		return m.resourceGroups.View()
	case ViewPolicies:
		return m.policies.View()
	case ViewArtifacts:
		return m.artifacts.View()
	case ViewTaskDetail:
		return m.taskDetail.View()
	}

	return ""
}

func (m AppModel) switchToView(
	v View,
) (AppModel, tea.Cmd) {
	m.prevView = m.activeView
	m.activeView = v
	m.tabs.setView(v)

	switch v {
	case ViewTasks:
		if len(m.tasks.Tasks) == 0 && !m.tasks.Loading {
			m.tasks.Loading = true

			return m, m.tasks.Load(m.client)
		}
	case ViewUsers:
		if len(m.users.Users) == 0 && !m.users.Loading {
			m.users.Loading = true

			return m, m.users.Load(m.client)
		}
	case ViewRoles:
		if len(m.roles.Roles) == 0 && !m.roles.Loading {
			m.roles.Loading = true

			return m, m.roles.Load(m.client)
		}
	case ViewResourceGroups:
		if len(m.resourceGroups.RGs) == 0 && !m.resourceGroups.Loading {
			m.resourceGroups.Loading = true

			return m, m.resourceGroups.Load(m.client)
		}
	case ViewPolicies:
		if len(m.policies.Policies) == 0 && !m.policies.Loading {
			m.policies.Loading = true

			return m, m.policies.Load(m.client)
		}
	case ViewArtifacts:
		if len(m.artifacts.Items) == 0 && !m.artifacts.Loading {
			m.artifacts.Loading = true

			return m, m.artifacts.Load(m.client)
		}
	case ViewTaskDetail:
		// TaskDetail is entered via enter key, not direct navigation.
	}

	return m, nil
}

func (m AppModel) doRefresh() (AppModel, tea.Cmd) {
	switch m.activeView {
	case ViewTasks:
		m.tasks.Loading = true

		return m, m.tasks.Load(m.client)
	case ViewUsers:
		m.users.Loading = true

		return m, m.users.Load(m.client)
	case ViewRoles:
		m.roles.Loading = true

		return m, m.roles.Load(m.client)
	case ViewResourceGroups:
		m.resourceGroups.Loading = true

		return m, m.resourceGroups.Load(m.client)
	case ViewPolicies:
		m.policies.Loading = true

		return m, m.policies.Load(m.client)
	case ViewArtifacts:
		m.artifacts.Loading = true

		return m, m.artifacts.Load(m.client)
	case ViewTaskDetail:
		// TaskDetail refreshes by reloading its task logs.
	}

	return m, nil
}

func (m AppModel) isCapturing() bool {
	switch m.activeView {
	case ViewUsers:
		return m.users.IsCapturing()
	case ViewRoles:
		return m.roles.IsCapturing()
	case ViewResourceGroups:
		return m.resourceGroups.IsCapturing()
	case ViewPolicies:
		return m.policies.IsCapturing()
	}

	return false
}

func (m AppModel) delegateKey(
	msg tea.KeyMsg,
) (AppModel, tea.Cmd) {
	return m.delegateMsg(msg)
}

func (m AppModel) delegateMsg(
	msg tea.Msg,
) (AppModel, tea.Cmd) {
	switch m.activeView {
	case ViewTasks:
		m.tasks, _ = m.tasks.Update(msg)
	case ViewUsers:
		var cmd tea.Cmd
		m.users, cmd = m.users.Update(msg)

		return m, cmd
	case ViewRoles:
		var cmd tea.Cmd
		m.roles, cmd = m.roles.Update(msg)

		return m, cmd
	case ViewResourceGroups:
		var cmd tea.Cmd
		m.resourceGroups, cmd = m.resourceGroups.Update(msg)

		return m, cmd
	case ViewPolicies:
		var cmd tea.Cmd
		m.policies, cmd = m.policies.Update(msg)

		return m, cmd
	case ViewArtifacts:
		var cmd tea.Cmd
		m.artifacts, cmd = m.artifacts.Update(msg, m.client)

		return m, cmd
	case ViewTaskDetail:
		m.taskDetail, _ = m.taskDetail.Update(msg)
	}

	return m, nil
}

// --- async API helpers ---

func (m AppModel) deleteUserCmd(name string) tea.Cmd {
	c := m.client

	return func() tea.Msg {
		_, err := c.DeleteUser(context.Background(), name)

		return views.UserDeletedMsg{Err: err}
	}
}

func (m AppModel) createUserCmd(name, display, pass string) tea.Cmd {
	c := m.client

	return func() tea.Msg {
		_, err := c.CreateUser(context.Background(), name, pass, display)

		return views.UserCreatedMsg{Err: err}
	}
}

func (m AppModel) deleteRoleCmd(name string) tea.Cmd {
	c := m.client

	return func() tea.Msg {
		_, err := c.DeleteRole(context.Background(), name)

		return views.RoleDeletedMsg{Err: err}
	}
}

func (m AppModel) createRoleCmd(name, usersRaw string) tea.Cmd {
	c := m.client
	users := splitTrim(usersRaw)

	return func() tea.Msg {
		_, err := c.CreateRole(context.Background(), name, users)

		return views.RoleCreatedMsg{Err: err}
	}
}

func (m AppModel) deleteRGCmd(name string) tea.Cmd {
	c := m.client

	return func() tea.Msg {
		_, err := c.DeleteResourceGroup(context.Background(), name)

		return views.ResourceGroupDeletedMsg{Err: err}
	}
}

func (m AppModel) createRGCmd(name, endpointsRaw string) tea.Cmd {
	c := m.client
	endpoints := splitTrim(endpointsRaw)

	return func() tea.Msg {
		_, err := c.CreateResourceGroup(context.Background(), name, endpoints)

		return views.ResourceGroupCreatedMsg{Err: err}
	}
}

func (m AppModel) deletePolicyCmd(p enclave.Policy) tea.Cmd {
	c := m.client

	return func() tea.Msg {
		err := c.DeletePolicy(context.Background(), p)

		return views.PolicyDeletedMsg{Err: err}
	}
}

func (m AppModel) createPolicyCmd(role, rg, method string) tea.Cmd {
	c := m.client
	p := enclave.Policy{
		Role:          role,
		ResourceGroup: rg,
		Method:        enclave.PolicyMethod(method),
	}

	return func() tea.Msg {
		err := c.CreatePolicy(context.Background(), p)

		return views.PolicyCreatedMsg{Err: err}
	}
}

// splitTrim splits a comma-separated string and trims whitespace from each part.
func splitTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}
