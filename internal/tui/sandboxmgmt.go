package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/TylerTwoForks/spin-kit/internal/sf"
)

func updateSandboxMgmt(m model, msg tea.Msg) (model, tea.Cmd) {
	if m.confirmDelete {
		return updateSandboxDeleteConfirm(m, msg)
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Back):
			m.state = viewSettings
			m.mgmtCursor = 0
			return m, nil

		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Add):
			m.inputPurpose = inputAddSandbox
			m.inputLabel = "Enter alias to save:"
			m.state = viewInput
			m.err = nil
			return m, m.textInput.Focus()

		case key.Matches(msg, keys.Delete):
			if len(m.sandboxes) > 0 && m.mgmtCursor < len(m.sandboxes) {
				sb := m.sandboxes[m.mgmtCursor]
				m.confirmDelete = true
				m.confirmDeleteName = sb.Name
				m.confirmDeleteID = sb.ID
				m.err = nil
				return m, nil
			}

		case msg.String() == "up" || msg.String() == "k":
			if m.mgmtCursor > 0 {
				m.mgmtCursor--
			}

		case msg.String() == "down" || msg.String() == "j":
			if m.mgmtCursor < len(m.sandboxes)-1 {
				m.mgmtCursor++
			}
		}
	}

	return m, nil
}

func updateSandboxDeleteConfirm(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			m.confirmDelete = false
			if err := m.store.RemoveSandbox(m.confirmDeleteID); err != nil {
				m.err = fmt.Errorf("remove alias: %w", err)
				return m, nil
			}
			if m.mgmtCursor >= len(m.sandboxes)-1 && m.mgmtCursor > 0 {
				m.mgmtCursor--
			}
			m.loading = true
			m.outputTitle = "Logging out " + m.confirmDeleteName
			m.state = viewOutput
			m.err = nil
			return m, tea.Batch(
				m.loadDataCmd(),
				sf.LogoutOrg(m.confirmDeleteName),
				func() tea.Msg { return m.spinner.Tick() },
			)

		case "n", "N", "esc":
			m.confirmDelete = false
			m.confirmDeleteName = ""
			m.confirmDeleteID = 0
			return m, nil
		}
	}
	return m, nil
}

func renderSandboxMgmt(m model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("SPIN Kit"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Manage Aliases"))
	b.WriteString("\n")

	if m.activeProd != nil {
		b.WriteString(dimStyle.Render("Prod Org: "))
		b.WriteString(activeOrgStyle.Render(m.activeProd.Alias))
		b.WriteString("\n\n")
	}

	if len(m.sandboxes) == 0 {
		b.WriteString(dimStyle.Render("  No aliases saved. Press 'a' to add one."))
		b.WriteString("\n")
	} else {
		for i, sb := range m.sandboxes {
			cursor := "  "
			if i == m.mgmtCursor {
				cursor = promptStyle.Render("> ")
			}
			b.WriteString(cursor)
			if i == m.mgmtCursor {
				b.WriteString(activeOrgStyle.Render(sb.Name))
			} else {
				b.WriteString(sb.Name)
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	if m.confirmDelete {
		b.WriteString(errorStyle.Render(
			fmt.Sprintf("Deleting %q will also log you out. Continue? (Y/n)", m.confirmDeleteName),
		))
	} else {
		b.WriteString(helpStyle.Render("a add • d delete • esc back"))
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.err.Error()))
	}

	return b.String()
}
