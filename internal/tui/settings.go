package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

var settingsOptions = []string{
	"Manage Aliases",
	"Manage Production Connections",
}

func updateSettings(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Back):
			m.state = viewMenu
			m.mgmtCursor = 0
			m.err = nil
			return m, nil

		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Select):
			switch m.mgmtCursor {
			case 0:
				if m.activeProd == nil {
					m.err = fmt.Errorf("add a production connection first")
					return m, nil
				}
				m.state = viewSandboxMgmt
				m.mgmtCursor = 0
				m.err = nil
				return m, nil
			case 1:
				m.state = viewProdConnections
				m.mgmtCursor = 0
				m.err = nil
				return m, nil
			}

		case msg.String() == "up" || msg.String() == "k":
			if m.mgmtCursor > 0 {
				m.mgmtCursor--
			}

		case msg.String() == "down" || msg.String() == "j":
			if m.mgmtCursor < len(settingsOptions)-1 {
				m.mgmtCursor++
			}
		}
	}

	return m, nil
}

func renderSettings(m model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("SPIN Kit"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Settings"))
	b.WriteString("\n\n")

	for i, option := range settingsOptions {
		cursor := "  "
		if i == m.mgmtCursor {
			cursor = promptStyle.Render("> ")
			b.WriteString(cursor + activeOrgStyle.Render(option) + "\n")
		} else {
			b.WriteString(cursor + option + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("enter select • esc back"))

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.err.Error()))
	}

	return b.String()
}
