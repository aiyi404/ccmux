package pages

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/farion1231/ccmux/internal/cli/i18n"
	"github.com/farion1231/ccmux/internal/cli/tui/styles"
	"github.com/farion1231/ccmux/internal/config"
	"github.com/farion1231/ccmux/internal/store"
	"github.com/mattn/go-runewidth"
)

type SettingsModel struct {
	state       *store.AppState
	choices     []string
	cursor      int
	langMenu    bool
	langChoices []string
	langCursor  int
}

type EditConfigMsg struct{}

func langDisplay(lang string) string {
	if lang == "zh" {
		return "中文"
	}
	return "English"
}

func NewSettings(state *store.AppState) SettingsModel {
	return SettingsModel{
		state: state,
		choices: []string{
			fmt.Sprintf("%s (%s)", i18n.T("switch_lang"), langDisplay(state.Lang)),
			i18n.T("open_config"),
			i18n.T("back"),
		},
		langChoices: []string{"English", "中文"},
	}
}

func (m SettingsModel) Init() tea.Cmd { return nil }

func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.langMenu {
			switch msg.String() {
			case "up", "k":
				if m.langCursor > 0 {
					m.langCursor--
				}
			case "down", "j":
				if m.langCursor < len(m.langChoices)-1 {
					m.langCursor++
				}
			case "enter":
				newLang := "en"
				if m.langCursor == 1 {
					newLang = "zh"
				}
				m.state.Config.Lang = newLang
				m.state.Lang = newLang
				i18n.SetLang(newLang)
				config.SaveConfig(m.state.Config)
				m.langMenu = false
				m.choices[0] = fmt.Sprintf("%s (%s)", i18n.T("switch_lang"), langDisplay(newLang))
			case "esc":
				m.langMenu = false
			}
			return m, nil
		}
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
// PLACEHOLDER_SETTINGS_UPDATE_CONT
		case "enter":
			switch m.cursor {
			case 0:
				m.langMenu = true
				m.langCursor = 0
			case 1:
				return m, func() tea.Msg { return EditConfigMsg{} }
			case 2:
				return m, func() tea.Msg { return GoBackMsg{} }
			}
		case "esc", "q":
			return m, func() tea.Msg { return GoBackMsg{} }
		}
	}
	return m, nil
}

// padLabel pads a label to a fixed display width, accounting for CJK characters
func padLabel(label string, width int) string {
	w := runewidth.StringWidth(label)
	if w >= width {
		return label
	}
	return label + strings.Repeat(" ", width-w)
}

func (m SettingsModel) View() string {
	s := styles.TitleStyle.Render("⚙ "+i18n.T("settings")) + "\n\n"
	const labelW = 16
	s += styles.Dim.Render(padLabel(i18n.T("config_path"), labelW)) + styles.ValueStyle.Render(config.CCCConfig) + "\n"
	s += styles.Dim.Render(padLabel(i18n.T("profiles_dir"), labelW)) + styles.ValueStyle.Render(config.CCCProfiles) + "\n"
	s += styles.Dim.Render(padLabel(i18n.T("lang_label"), labelW)) + styles.ValueStyle.Render(langDisplay(m.state.Lang)) + "\n\n"
	if m.langMenu {
		s += styles.Bold.Render(i18n.T("select_lang")) + "\n"
		for i, c := range m.langChoices {
			if i == m.langCursor {
				s += styles.MenuSelectedStyle.Render("▸ "+c) + "\n"
			} else {
				s += styles.MenuItemStyle.Render("  "+c) + "\n"
			}
		}
		return s
	}
	for i, choice := range m.choices {
		if i == m.cursor {
			s += styles.MenuSelectedStyle.Render("▸ "+choice) + "\n"
		} else {
			s += styles.MenuItemStyle.Render("  "+choice) + "\n"
		}
	}
	return s
}
