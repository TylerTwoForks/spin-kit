package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/TylerTwoForks/spin-kit/internal/sf"
)

func updateInput(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Back):
			m.textInput.SetValue("")
			switch m.inputPurpose {
			case inputRefreshCustom:
				m.state = viewRefreshList
			case inputReconnect:
				m.state = viewReconnectList
			case inputAddSandbox:
				m.state = viewSandboxMgmt
			case inputAddProdOrg:
				m.state = viewProdConnections
			default:
				m.state = viewMenu
			}
			return m, nil

		case key.Matches(msg, keys.Select):
			val := strings.TrimSpace(m.textInput.Value())
			if val == "" {
				return m, nil
			}
			m.textInput.SetValue("")
			return m.handleInputSubmit(val)
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) handleInputSubmit(val string) (model, tea.Cmd) {
	switch m.inputPurpose {
	case inputProdConnect:
		m.loading = true
		m.outputTitle = "Connecting to " + val
		m.state = viewOutput
		return m, tea.Batch(
			sf.WebLogin(val),
			func() tea.Msg { return m.spinner.Tick() },
		)

	case inputRefreshCustom:
		if m.activeProd == nil {
			m.err = fmt.Errorf("no active production org configured")
			return m, nil
		}
		return m, sf.RefreshSandbox(val, m.activeProd.Alias)

	case inputReconnect:
		m.loading = true
		m.outputTitle = "Reconnecting to " + val
		m.state = viewOutput
		return m, tea.Batch(
			sf.ReconnectSandbox(val),
			func() tea.Msg { return m.spinner.Tick() },
		)

	case inputAddSandbox:
		if m.activeProd == nil {
			m.err = fmt.Errorf("no active production org configured")
			return m, nil
		}
		if err := m.store.AddSandbox(val, m.activeProd.ID); err != nil {
			m.err = fmt.Errorf("add sandbox: %w", err)
			return m, nil
		}
		m.state = viewSandboxMgmt
		return m, m.loadDataCmd()

	case inputAddProdOrg:
		if err := m.store.AddProdOrg(val); err != nil {
			m.err = fmt.Errorf("add prod org: %w", err)
			return m, nil
		}
		m.loading = true
		m.outputTitle = "Connecting to " + val
		m.state = viewOutput
		return m, tea.Batch(
			m.loadDataCmd(),
			sf.WebLogin(val),
			func() tea.Msg { return m.spinner.Tick() },
		)
	}

	return m, nil
}

func renderInput(m model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("SPIN Kit"))
	b.WriteString("\n\n")
	b.WriteString(promptStyle.Render(m.inputLabel))
	b.WriteString("\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("enter submit • esc cancel"))

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.err.Error()))
	}

	return b.String()
}
