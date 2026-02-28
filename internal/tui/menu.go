package tui

import "charm.land/bubbles/v2/list"

type menuAction int

const (
	actionListOrgs menuAction = iota
	actionRefreshSandbox
	actionReconnectSandbox
	actionSandboxStatus
	actionSettings
)

type menuItem struct {
	title       string
	description string
	action      menuAction
	sandboxName string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

func (m model) buildMenuItems() []list.Item {
	return []list.Item{
		menuItem{title: "List Org Connections", description: "Show all connected orgs", action: actionListOrgs},
		menuItem{title: "Refresh Sandbox", description: "Refresh a saved or custom sandbox", action: actionRefreshSandbox},
		menuItem{title: "Reconnect to Sandbox", description: "Re-authenticate to a sandbox via browser", action: actionReconnectSandbox},
		menuItem{title: "Sandbox Refresh Status", description: "Check progress of in-flight refreshes", action: actionSandboxStatus},
		menuItem{title: "Settings", description: "Manage sandboxes and production connections", action: actionSettings},
	}
}
