package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/aiyi404/ccmux/internal/cli/tui/styles"
)

type ConfirmModel struct {
	Message  string
	selected int
	Decided  bool
	Result   bool
}

type ConfirmResultMsg struct{ Confirmed bool }

func NewConfirm(message string) ConfirmModel {
	return ConfirmModel{Message: message, selected: 1}
}

func (m ConfirmModel) Init() tea.Cmd { return nil }

func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.selected = 0
		case "right", "l":
			m.selected = 1
		case "y", "Y":
			m.Decided = true
			m.Result = true
			return m, func() tea.Msg { return ConfirmResultMsg{Confirmed: true} }
		case "n", "N", "esc":
			m.Decided = true
			m.Result = false
			return m, func() tea.Msg { return ConfirmResultMsg{Confirmed: false} }
		case "enter":
			m.Decided = true
			m.Result = m.selected == 0
			return m, func() tea.Msg { return ConfirmResultMsg{Confirmed: m.Result} }
		}
	}
	return m, nil
}

func (m ConfirmModel) View() string {
	yes := " Yes "
	no := " No "
	if m.selected == 0 {
		yes = styles.MenuSelectedStyle.Render("▸ Yes ")
		no = styles.Dim.Render("  No ")
	} else {
		yes = styles.Dim.Render("  Yes ")
		no = styles.MenuSelectedStyle.Render("▸ No ")
	}
	return m.Message + "\n\n" + yes + "  " + no
}
