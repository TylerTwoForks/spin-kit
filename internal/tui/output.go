package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/TylerTwoForks/spin-kit/internal/sf"
)

func updateOutput(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Back), key.Matches(msg, keys.Quit):
			m.state = viewMenu
			m.showTable = false
			m.err = nil
			return m, nil

		case key.Matches(msg, keys.Select):
			if m.showTable && m.outputTitle == "Org Connections" {
				row := m.resultTable.SelectedRow()
				if len(row) == 0 {
					m.err = fmt.Errorf("no org selected")
					return m, nil
				}

				alias := strings.TrimSpace(row[0])
				if alias == "" {
					m.err = fmt.Errorf("selected org has no alias; set an alias in sf and retry")
					return m, nil
				}

				m.selectedOrgAlias = alias
				if len(row) > 1 {
					m.selectedOrgUsername = strings.TrimSpace(row[1])
				} else {
					m.selectedOrgUsername = ""
				}
				m.mgmtCursor = 0
				m.statusMsg = ""
				m.err = nil
				m.state = viewOrgActions
				return m, nil
			}
		}
	}

	if m.showTable {
		var cmd tea.Cmd
		m.resultTable, cmd = m.resultTable.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func renderOutput(m model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("SPIN Kit"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render(m.outputTitle))
	b.WriteString("\n")

	if m.loading {
		b.WriteString(spinnerStyle.Render(m.spinner.View()))
		b.WriteString(" Running command...")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("esc/q back to menu"))
		return b.String()
	}

	if m.showTable {
		b.WriteString(m.resultTable.View())
	} else {
		b.WriteString(m.viewport.View())
	}

	b.WriteString("\n")
	help := "↑↓ scroll • esc/q back to menu"
	if m.showTable && m.outputTitle == "Org Connections" {
		help = "↑↓ scroll • enter actions • esc/q back to menu"
	}
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

func fitColumns(titles []string, rows []table.Row, padding, maxWidth int) []table.Column {
	widths := make([]int, len(titles))
	for i, t := range titles {
		widths[i] = lipgloss.Width(t)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				if w := lipgloss.Width(cell); w > widths[i] {
					widths[i] = w
				}
			}
		}
	}

	cols := make([]table.Column, len(titles))
	for i, t := range titles {
		w := widths[i] + padding
		if maxWidth > 0 && w > maxWidth {
			w = maxWidth
		}
		cols[i] = table.Column{Title: t, Width: w}
	}
	return cols
}

func buildOrgTable(orgs []sf.OrgInfo, width int) table.Model {
	rows := make([]table.Row, 0, len(orgs))
	for _, o := range orgs {
		status := o.ConnectedStatus
		if strings.EqualFold(status, "Connected") {
			status = successStyle.Render(status)
		} else {
			status = errorStyle.Render(status)
		}
		rows = append(rows, table.Row{o.Alias, o.Username, o.OrgID, status, o.InstanceURL})
	}

	titles := []string{"Alias", "Username", "Org ID", "Status", "Instance URL"}
	maxCol := max(width/len(titles), 12)
	columns := fitColumns(titles, rows, 2, maxCol)

	totalWidth := 0
	for _, c := range columns {
		totalWidth += c.Width
	}

	s := table.Styles{
		Header:   tableHeaderStyle,
		Cell:     lipgloss.NewStyle().Padding(0, 1),
		Selected: tableSelectedStyle,
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithWidth(totalWidth),
		table.WithHeight(min(len(rows)+1, 20)),
		table.WithStyles(s),
		table.WithFocused(true),
		table.WithRows(rows),
	)

	return t
}

func buildStatusTable(records []sf.SandboxProcessRecord, width int) table.Model {
	rows := make([]table.Row, 0, len(records))
	for _, r := range records {
		rows = append(rows, table.Row{r.SandboxName, r.Status, fmt.Sprintf("%d%%", r.CopyProgress)})
	}

	titles := []string{"Sandbox Name", "Status", "Copy Progress"}
	maxCol := max(width/len(titles), 15)
	columns := fitColumns(titles, rows, 2, maxCol)

	totalWidth := 0
	for _, c := range columns {
		totalWidth += c.Width
	}

	s := table.Styles{
		Header:   tableHeaderStyle,
		Cell:     lipgloss.NewStyle().Padding(0, 1),
		Selected: tableSelectedStyle,
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithWidth(totalWidth),
		table.WithHeight(min(len(rows)+1, 20)),
		table.WithStyles(s),
		table.WithFocused(true),
		table.WithRows(rows),
	)

	return t
}
