package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/TylerTwoForks/spin-kit/internal/sf"
)

func updateRefreshList(m model, msg tea.Msg) (model, tea.Cmd) {
	maxIdx := len(m.sandboxes) // index 0 = "non-listed", 1..N = sandboxes

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Back):
			m.state = viewMenu
			m.mgmtCursor = 0
			return m, nil

		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Select):
			if m.mgmtCursor == 0 {
				m.inputPurpose = inputRefreshCustom
				m.inputLabel = "Enter sandbox name to refresh:"
				m.state = viewInput
				m.err = nil
				return m, m.textInput.Focus()
			}

			sb := m.sandboxes[m.mgmtCursor-1]
			m.err = nil
			return m, sf.RefreshSandbox(sb.Name, m.activeProd.Alias)

		case msg.String() == "up" || msg.String() == "k":
			if m.mgmtCursor > 0 {
				m.mgmtCursor--
			}

		case msg.String() == "down" || msg.String() == "j":
			if m.mgmtCursor < maxIdx {
				m.mgmtCursor++
			}
		}
	}

	return m, nil
}

func renderRefreshList(m model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("SPIN Kit"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Refresh Sandbox"))
	b.WriteString("\n")

	if m.activeProd != nil {
		b.WriteString(dimStyle.Render("Prod Org: "))
		b.WriteString(activeOrgStyle.Render(m.activeProd.Alias))
		b.WriteString("\n\n")
	}

	// First item: type a custom name
	cursor := "  "
	if m.mgmtCursor == 0 {
		cursor = promptStyle.Render("> ")
	}
	label := "Refresh non-listed sandbox"
	if m.mgmtCursor == 0 {
		label = activeOrgStyle.Render(label)
	}
	b.WriteString(fmt.Sprintf("%s%s\n", cursor, label))

	if len(m.sandboxes) > 0 {
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("  Saved sandboxes:"))
		b.WriteString("\n")
	}

	for i, sb := range m.sandboxes {
		idx := i + 1
		cursor := "  "
		if m.mgmtCursor == idx {
			cursor = promptStyle.Render("> ")
		}
		name := sb.Name
		if m.mgmtCursor == idx {
			name = activeOrgStyle.Render(name)
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, name))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("enter select • esc back"))

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.err.Error()))
	}

	return b.String()
}
