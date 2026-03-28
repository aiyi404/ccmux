package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/farion1231/ccmux/internal/cli/i18n"
	"github.com/farion1231/ccmux/internal/cli/tui/components"
	"github.com/farion1231/ccmux/internal/cli/tui/styles"
	"github.com/farion1231/ccmux/internal/services"
	"github.com/farion1231/ccmux/internal/store"
)

type ProviderAction int

const (
	ProviderActionNone ProviderAction = iota
	ProviderActionViewDetails
	ProviderActionSwitchGlobal
	ProviderActionUseSession
	ProviderActionAdd
	ProviderActionEdit
	ProviderActionRemove
	ProviderActionImport
	ProviderActionBack
)

type ProvidersModel struct {
	state      *store.AppState
	table      components.ProviderTable
	actionMenu bool
	actions    []string
	actionMap  []ProviderAction
	actionIdx  int
	info       string
}

type ProviderActionMsg struct {
	Action ProviderAction
	Name   string
}

type GoBackMsg struct{}

func NewProviders(state *store.AppState) ProvidersModel {
	providers, _ := state.Service.List()
	current, _ := state.Service.GetCurrent()
	currentName := ""
	if current != nil {
		currentName = current.Name
	}

	actions := []string{i18n.T("view_details"), i18n.T("switch_global"), i18n.T("use_session")}
	actionMap := []ProviderAction{ProviderActionViewDetails, ProviderActionSwitchGlobal, ProviderActionUseSession}

	if state.Mode == "standalone" {
		actions = append(actions, i18n.T("add_provider"), i18n.T("edit_provider"), i18n.T("remove_provider"), i18n.T("import_settings"))
		actionMap = append(actionMap, ProviderActionAdd, ProviderActionEdit, ProviderActionRemove, ProviderActionImport)
	}
	actions = append(actions, i18n.T("back"))
	actionMap = append(actionMap, ProviderActionBack)

	return ProvidersModel{
		state:     state,
		table:     components.NewProviderTable(providers, currentName, 80),
		actions:   actions,
		actionMap: actionMap,
	}
}

func (m ProvidersModel) Init() tea.Cmd { return nil }

func (m ProvidersModel) Update(msg tea.Msg) (ProvidersModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.info != "" {
			m.info = ""
			return m, nil
		}
		if m.actionMenu {
			switch msg.String() {
			case "up", "k":
				if m.actionIdx > 0 {
					m.actionIdx--
				}
			case "down", "j":
				if m.actionIdx < len(m.actions)-1 {
					m.actionIdx++
				}
			case "enter":
				action := m.actionMap[m.actionIdx]
				name := m.table.SelectedName()
				m.actionMenu = false
				if action == ProviderActionBack {
					return m, func() tea.Msg { return GoBackMsg{} }
				}
				return m, func() tea.Msg { return ProviderActionMsg{Action: action, Name: name} }
			case "esc":
				m.actionMenu = false
			}
			return m, nil
		}
		switch msg.String() {
		case "enter":
			m.actionMenu = true
			m.actionIdx = 0
			return m, nil
		case "esc", "q":
			return m, func() tea.Msg { return GoBackMsg{} }
		}
	case components.ProviderSelectedMsg:
		m.actionMenu = true
		m.actionIdx = 0
		return m, nil
	}
	if !m.actionMenu {
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m ProvidersModel) View() string {
	s := styles.TitleStyle.Render(i18n.T("providers")) + "\n\n"
	s += m.table.View() + "\n\n"
	if m.info != "" {
		s += m.info + "\n\n" + styles.Dim.Render(i18n.T("press_enter"))
		return s
	}
	if m.actionMenu {
		s += styles.Bold.Render(i18n.T("select_action")) + "\n"
		for i, a := range m.actions {
			if i == m.actionIdx {
				s += styles.MenuSelectedStyle.Render("▸ "+a) + "\n"
			} else {
				s += styles.MenuItemStyle.Render("  "+a) + "\n"
			}
		}
	} else {
		s += styles.Dim.Render("Enter: actions  q: back")
	}
	return s
}

func (m *ProvidersModel) SetInfo(msg string) { m.info = msg }

func (m *ProvidersModel) Refresh() {
	providers, _ := m.state.Service.List()
	current, _ := m.state.Service.GetCurrent()
	currentName := ""
	if current != nil {
		currentName = current.Name
	}
	m.table = components.NewProviderTable(providers, currentName, 80)
}

func ShowDetails(state *store.AppState, name string) string {
	p, err := state.Service.GetByName(name)
	if err != nil {
		return styles.ErrorStyle.Render(fmt.Sprintf("error: %v", err))
	}
	s := styles.Bold.Render(p.Name) + "\n"
	if p.Description != "" {
		s += fmt.Sprintf("  description: %s\n", p.Description)
	}
	for k, v := range p.Env {
		s += fmt.Sprintf("  %s: %s\n", k, v)
	}
	if p.ModelAlias != "" {
		s += fmt.Sprintf("  model_alias: %s\n", p.ModelAlias)
	}
	return s
}

// compile-time check that services.Provider is used
var _ = services.Provider{}
