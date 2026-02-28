package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func updateProdConnections(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Back):
			if m.firstRun && m.activeProd == nil {
				return m, nil
			}
			m.state = viewSettings
			m.mgmtCursor = 0
			return m, nil

		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Add):
			m.inputPurpose = inputAddProdOrg
			m.inputLabel = "Enter production org alias:"
			m.state = viewInput
			m.err = nil
			return m, m.textInput.Focus()

		case key.Matches(msg, keys.Delete):
			if len(m.prodOrgs) > 0 && m.mgmtCursor < len(m.prodOrgs) {
				org := m.prodOrgs[m.mgmtCursor]
				if org.IsActive && len(m.prodOrgs) > 1 {
					m.err = fmt.Errorf("cannot delete the active org; switch active org first")
					return m, nil
				}
				if err := m.store.RemoveProdOrg(org.ID); err != nil {
					m.err = fmt.Errorf("remove prod org: %w", err)
					return m, nil
				}
				if m.mgmtCursor >= len(m.prodOrgs)-1 && m.mgmtCursor > 0 {
					m.mgmtCursor--
				}
				return m, m.loadDataCmd()
			}

		case key.Matches(msg, keys.Select):
			if len(m.prodOrgs) > 0 && m.mgmtCursor < len(m.prodOrgs) {
				org := m.prodOrgs[m.mgmtCursor]
				if err := m.store.SetActiveProdOrg(org.ID); err != nil {
					m.err = fmt.Errorf("set active org: %w", err)
					return m, nil
				}
				return m, m.loadDataCmd()
			}

		case msg.String() == "up" || msg.String() == "k":
			if m.mgmtCursor > 0 {
				m.mgmtCursor--
			}

		case msg.String() == "down" || msg.String() == "j":
			if m.mgmtCursor < len(m.prodOrgs)-1 {
				m.mgmtCursor++
			}
		}
	}

	return m, nil
}

func renderProdConnections(m model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("SPIN Kit"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Manage Production Connections"))
	b.WriteString("\n\n")

	if m.firstRun && m.activeProd == nil {
		b.WriteString(promptStyle.Render("Welcome! Add your first production org to get started."))
		b.WriteString("\n\n")
	}

	if len(m.prodOrgs) == 0 {
		b.WriteString(dimStyle.Render("  No production orgs configured. Press 'a' to add one."))
		b.WriteString("\n")
	} else {
		for i, org := range m.prodOrgs {
			cursor := "  "
			if i == m.mgmtCursor {
				cursor = promptStyle.Render("> ")
			}
			b.WriteString(cursor)

			label := org.Alias
			if org.IsActive {
				label += " " + activeOrgStyle.Render("(active)")
			}

			if i == m.mgmtCursor {
				b.WriteString(activeOrgStyle.Render(org.Alias))
				if org.IsActive {
					b.WriteString(" " + activeOrgStyle.Render("(active)"))
				}
			} else {
				b.WriteString(label)
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	help := "a add • enter set active • d delete • esc back"
	if m.firstRun && m.activeProd == nil {
		help = "a add your first production org"
	}
	b.WriteString(helpStyle.Render(help))

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.err.Error()))
	}

	return b.String()
}
