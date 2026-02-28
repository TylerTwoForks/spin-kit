package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/TylerTwoForks/spin-kit/internal/db"
	"github.com/TylerTwoForks/spin-kit/internal/sf"
)

type viewState int

const (
	viewMenu viewState = iota
	viewInput
	viewOutput
	viewOrgActions
	viewRefreshList
	viewReconnectList
	viewSandboxMgmt
	viewSettings
	viewProdConnections
)

type inputPurpose int

const (
	inputProdConnect inputPurpose = iota
	inputRefreshCustom
	inputReconnect
	inputAddSandbox
	inputAddProdOrg
)

type dataLoadedMsg struct {
	activeProd *db.ProdOrg
	sandboxes  []db.Sandbox
	prodOrgs   []db.ProdOrg
	err        error
}

type model struct {
	state viewState

	store      db.Store
	activeProd *db.ProdOrg
	sandboxes  []db.Sandbox
	prodOrgs   []db.ProdOrg

	menu      list.Model
	textInput textinput.Model

	viewport    viewport.Model
	resultTable table.Model
	showTable   bool
	spinner     spinner.Model
	loading     bool
	outputTitle string

	inputPurpose inputPurpose
	inputLabel   string

	mgmtCursor int

	err       error
	statusMsg string

	selectedOrgAlias    string
	selectedOrgUsername string

	confirmDelete     bool
	confirmDeleteName string
	confirmDeleteID   int64

	width    int
	height   int
	ready    bool
	firstRun bool
}

func NewModel(store db.Store) model {
	ti := textinput.New()
	ti.Placeholder = "type here..."
	ti.CharLimit = 100

	sp := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(spinnerStyle),
	)

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	menuList := list.New([]list.Item{}, delegate, 0, 0)
	menuList.Title = "SPIN Kit"
	menuList.Styles.Title = titleStyle

	vp := viewport.New()
	rt := table.New(
		table.WithColumns([]table.Column{}),
		table.WithRows([]table.Row{}),
		table.WithStyles(table.Styles{
			Header:   tableHeaderStyle,
			Cell:     lipgloss.NewStyle().Padding(0, 1),
			Selected: tableSelectedStyle,
		}),
	)

	return model{
		state:       viewMenu,
		store:       store,
		menu:        menuList,
		textInput:   ti,
		viewport:    vp,
		resultTable: rt,
		spinner:     sp,
	}
}

func (m model) Init() tea.Cmd {
	return m.loadDataCmd()
}

func (m model) loadDataCmd() tea.Cmd {
	store := m.store
	return func() tea.Msg {
		activeProd, err := store.GetActiveProdOrg()
		if err != nil {
			return dataLoadedMsg{err: err}
		}

		var sandboxes []db.Sandbox
		if activeProd != nil {
			sandboxes, err = store.ListSandboxes(activeProd.ID)
			if err != nil {
				return dataLoadedMsg{err: err}
			}
		}

		prodOrgs, err := store.ListProdOrgs()
		if err != nil {
			return dataLoadedMsg{err: err}
		}

		return dataLoadedMsg{
			activeProd: activeProd,
			sandboxes:  sandboxes,
			prodOrgs:   prodOrgs,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global quit
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		if keyMsg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu.SetSize(m.width, m.height-2)
		m.viewport.SetWidth(m.width)
		m.viewport.SetHeight(m.height - 8)
		m.ready = true
		return m, nil

	case dataLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.activeProd = msg.activeProd
		m.sandboxes = msg.sandboxes
		m.prodOrgs = msg.prodOrgs
		m.err = nil

		if m.activeProd == nil {
			m.firstRun = true
			m.state = viewProdConnections
		} else {
			m.firstRun = false
		}

		items := m.buildMenuItems()
		cmd := m.menu.SetItems(items)
		if m.activeProd != nil {
			m.menu.Title = fmt.Sprintf("SPIN Kit  %s",
				dimStyle.Render("org: ")+activeOrgStyle.Render(m.activeProd.Alias))
		}
		return m, cmd

	case sf.OrgListMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			m.state = viewOutput
			return m, nil
		}

		if m.activeProd != nil {
			for _, o := range msg.Orgs {
				alias := strings.TrimSpace(o.Alias)
				if alias != "" && alias != m.activeProd.Alias {
					_ = m.store.EnsureSandbox(alias, m.activeProd.ID)
				}
			}
		}

		m.showTable = true
		m.resultTable = buildOrgTable(msg.Orgs, m.width)
		m.outputTitle = "Org Connections"
		m.selectedOrgAlias = ""
		m.selectedOrgUsername = ""
		m.state = viewOutput
		return m, m.loadDataCmd()

	case sf.SandboxStatusMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			m.state = viewOutput
			return m, nil
		}
		m.showTable = true
		m.resultTable = buildStatusTable(msg.Records, m.width)
		m.outputTitle = "Sandbox Refresh Status"
		m.state = viewOutput
		return m, nil

	case sf.SandboxRefreshMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			m.state = viewOutput
			return m, nil
		}
		m.showTable = false
		m.viewport.SetContent(successStyle.Render(msg.Result))
		m.outputTitle = "Refresh: " + msg.Name
		m.state = viewOutput
		return m, nil

	case sf.OrgLogoutMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			m.state = viewOutput
			return m, nil
		}
		if m.activeProd != nil {
			_ = m.store.RemoveSandboxByName(msg.Alias, m.activeProd.ID)
		}
		m.showTable = false
		m.viewport.SetContent(successStyle.Render(fmt.Sprintf("Logged out of %s and removed from saved aliases.", msg.Alias)))
		m.outputTitle = "Logout: " + msg.Alias
		m.state = viewOutput
		return m, m.loadDataCmd()

	case sf.WebLoginMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			m.state = viewOutput
			return m, nil
		}
		m.showTable = false
		result := successStyle.Render(fmt.Sprintf("Successfully connected to %s", msg.Alias))
		if msg.URL != "" {
			result += "\n" + dimStyle.Render(msg.URL)
		}
		m.viewport.SetContent(result)
		m.outputTitle = "Login: " + msg.Alias
		m.state = viewOutput
		return m, nil
	}

	// Spinner ticks while loading
	if m.loading {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Screen-specific handling
	switch m.state {
	case viewMenu:
		return updateMenuView(m, msg)
	case viewInput:
		return updateInput(m, msg)
	case viewOutput:
		return updateOutput(m, msg)
	case viewOrgActions:
		return updateOrgActions(m, msg)
	case viewRefreshList:
		return updateRefreshList(m, msg)
	case viewReconnectList:
		return updateReconnectList(m, msg)
	case viewSandboxMgmt:
		return updateSandboxMgmt(m, msg)
	case viewSettings:
		return updateSettings(m, msg)
	case viewProdConnections:
		return updateProdConnections(m, msg)
	}

	return m, nil
}

func updateMenuView(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			selected := m.menu.SelectedItem()
			if selected == nil {
				return m, nil
			}
			item := selected.(menuItem)
			return m.handleMenuAction(item)

		case "s":
			m.state = viewSettings
			m.mgmtCursor = 0
			m.err = nil
			return m, nil

		case "m":
			m.state = viewSettings
			m.mgmtCursor = 0
			m.err = nil
			return m, nil

		case "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

func (m model) handleMenuAction(item menuItem) (model, tea.Cmd) {
	switch item.action {
	case actionRefreshSandbox:
		if m.activeProd == nil {
			m.err = fmt.Errorf("no active production org configured")
			return m, nil
		}
		m.state = viewRefreshList
		m.mgmtCursor = 0
		m.err = nil
		return m, nil

	case actionListOrgs:
		m.loading = true
		m.outputTitle = "Org Connections"
		m.state = viewOutput
		m.err = nil
		return m, tea.Batch(
			sf.ListOrgs(),
			func() tea.Msg { return m.spinner.Tick() },
		)

	case actionSandboxStatus:
		if m.activeProd == nil {
			m.err = fmt.Errorf("no active production org configured")
			return m, nil
		}
		m.loading = true
		m.outputTitle = "Sandbox Refresh Status"
		m.state = viewOutput
		m.err = nil
		return m, tea.Batch(
			sf.SandboxStatus(m.activeProd.Alias),
			func() tea.Msg { return m.spinner.Tick() },
		)

	case actionReconnectSandbox:
		if m.activeProd == nil {
			m.err = fmt.Errorf("no active production org configured")
			return m, nil
		}
		m.state = viewReconnectList
		m.mgmtCursor = 0
		m.err = nil
		return m, nil

	case actionSettings:
		m.state = viewSettings
		m.mgmtCursor = 0
		m.err = nil
		return m, nil
	}

	return m, nil
}

func (m model) View() tea.View {
	if !m.ready {
		v := tea.NewView("Loading...")
		v.AltScreen = true
		return v
	}

	var s string
	switch m.state {
	case viewMenu:
		s = m.menu.View()
		if m.err != nil {
			s += "\n" + errorStyle.Render(m.err.Error())
		}
	case viewInput:
		s = renderInput(m)
	case viewOutput:
		s = renderOutput(m)
	case viewOrgActions:
		s = renderOrgActions(m)
	case viewRefreshList:
		s = renderRefreshList(m)
	case viewReconnectList:
		s = renderReconnectList(m)
	case viewSandboxMgmt:
		s = renderSandboxMgmt(m)
	case viewSettings:
		s = renderSettings(m)
	case viewProdConnections:
		s = renderProdConnections(m)
	}

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}
