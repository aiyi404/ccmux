package animations

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/aiyi404/ccmux/internal/cli/tui/styles"
)

type SpinnerModel struct {
	spinner spinner.Model
	Label   string
	Active  bool
}

func NewSpinner(label string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Bold.Foreground(styles.Purple)
	return SpinnerModel{spinner: s, Label: label, Active: true}
}

func (m SpinnerModel) Init() tea.Cmd { return m.spinner.Tick }

func (m SpinnerModel) Update(msg tea.Msg) (SpinnerModel, tea.Cmd) {
	if !m.Active {
		return m, nil
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m SpinnerModel) View() string {
	if !m.Active {
		return ""
	}
	return m.spinner.View() + " " + m.Label
}
