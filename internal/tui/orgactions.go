package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/TylerTwoForks/spin-kit/internal/sf"
)

var orgActionOptions = []string{
	"Open",
	"Reconnect",
	"Refresh",
	"Logout",
}

func updateOrgActions(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Back):
			m.state = viewOutput
			m.mgmtCursor = 0
			m.statusMsg = ""
			m.err = nil
			return m, nil

		case key.Matches(msg, keys.Quit):
			m.state = viewMenu
			m.showTable = false
			m.mgmtCursor = 0
			m.statusMsg = ""
			m.err = nil
			return m, nil

		case key.Matches(msg, keys.Select):
			return m.handleSelectedOrgAction()

		case msg.String() == "up" || msg.String() == "k":
			if m.mgmtCursor > 0 {
				m.mgmtCursor--
			}

		case msg.String() == "down" || msg.String() == "j":
			if m.mgmtCursor < len(orgActionOptions)-1 {
				m.mgmtCursor++
			}
		}
	}

	return m, nil
}

func (m model) handleSelectedOrgAction() (model, tea.Cmd) {
	alias := strings.TrimSpace(m.selectedOrgAlias)
	if alias == "" {
		m.err = fmt.Errorf("selected org has no alias")
		return m, nil
	}

	switch m.mgmtCursor {
	case 0:
		m.loading = true
		m.outputTitle = "Opening " + alias
		m.state = viewOutput
		m.statusMsg = ""
		m.err = nil
		return m, tea.Batch(
			sf.OpenOrg(alias),
			func() tea.Msg { return m.spinner.Tick() },
		)

	case 1:
		m.loading = true
		m.outputTitle = "Reconnecting to " + alias
		m.state = viewOutput
		m.statusMsg = ""
		m.err = nil
		return m, tea.Batch(
			sf.ReconnectSandbox(alias),
			func() tea.Msg { return m.spinner.Tick() },
		)

	case 2:
		if m.activeProd == nil {
			m.err = fmt.Errorf("no active production org configured")
			return m, nil
		}
		m.statusMsg = ""
		m.err = nil
		return m, sf.RefreshSandbox(alias, m.activeProd.Alias)

	case 3:
		m.loading = true
		m.outputTitle = "Logging out " + alias
		m.state = viewOutput
		m.statusMsg = ""
		m.err = nil
		return m, tea.Batch(
			sf.LogoutOrg(alias),
			func() tea.Msg { return m.spinner.Tick() },
		)
	}

	return m, nil
}

func renderOrgActions(m model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("SPIN Kit"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Org Connection Actions"))
	b.WriteString("\n")

	b.WriteString(dimStyle.Render("Selected alias: "))
	b.WriteString(activeOrgStyle.Render(m.selectedOrgAlias))
	b.WriteString("\n")

	if m.selectedOrgUsername != "" {
		b.WriteString(dimStyle.Render("Username: " + m.selectedOrgUsername))
		b.WriteString("\n")
	}

	if m.activeProd != nil {
		b.WriteString(dimStyle.Render("Active prod: "))
		b.WriteString(activeOrgStyle.Render(m.activeProd.Alias))
	} else {
		b.WriteString(errorStyle.Render("No active production org configured."))
	}
	b.WriteString("\n\n")

	for i, option := range orgActionOptions {
		cursor := "  "
		if i == m.mgmtCursor {
			cursor = promptStyle.Render("> ")
			b.WriteString(cursor + activeOrgStyle.Render(option) + "\n")
			continue
		}
		b.WriteString(cursor + option + "\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("enter select • esc back to list • q menu"))

	if m.statusMsg != "" {
		b.WriteString("\n")
		b.WriteString(successStyle.Render(m.statusMsg))
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.err.Error()))
	}

	return b.String()
}
