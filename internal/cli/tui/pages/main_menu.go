package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/farion1231/ccmux/internal/cli/i18n"
	"github.com/farion1231/ccmux/internal/cli/tui/styles"
)

type MainMenuChoice int

const (
	MainMenuProviders MainMenuChoice = iota
	MainMenuCurrent
	MainMenuSettings
	MainMenuExit
)

type MainMenuModel struct {
	choices  []string
	cursor   int
	selected MainMenuChoice
	decided  bool
}

type MainMenuSelectedMsg struct{ Choice MainMenuChoice }

func NewMainMenu() MainMenuModel {
	return MainMenuModel{
		choices: []string{
			i18n.T("manage_providers"),
			i18n.T("view_current"),
			i18n.T("settings_menu"),
			i18n.T("exit"),
		},
	}
}

func (m MainMenuModel) Init() tea.Cmd { return nil }

func (m MainMenuModel) Update(msg tea.Msg) (MainMenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.decided = true
			m.selected = MainMenuChoice(m.cursor)
			return m, func() tea.Msg { return MainMenuSelectedMsg{Choice: m.selected} }
		case "q", "esc":
			m.decided = true
			m.selected = MainMenuExit
			return m, func() tea.Msg { return MainMenuSelectedMsg{Choice: MainMenuExit} }
		}
	}
	return m, nil
}

func (m MainMenuModel) View() string {
	s := styles.TitleStyle.Render("  "+i18n.T("main_header")) + "\n\n"
	for i, choice := range m.choices {
		if i == m.cursor {
			s += styles.MenuSelectedStyle.Render("▸ "+choice) + "\n"
		} else {
			s += styles.MenuItemStyle.Render("  "+choice) + "\n"
		}
	}
	return s
}
