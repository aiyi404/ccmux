package animations

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var LogoLines = []string{
	"  ___  ___  ___  ",
	" / __)/ __)/ __) ",
	"( (__( (__( (__  ",
	" \\___)\\___)\\___) ",
}

type LogoModel struct {
	revealed int
	total    int
	done     bool
}

type logoTickMsg struct{}

func NewLogoModel() LogoModel {
	total := 0
	for _, l := range LogoLines {
		total += len([]rune(l))
	}
	return LogoModel{total: total}
}

func (m LogoModel) Init() tea.Cmd {
	return tea.Tick(30*time.Millisecond, func(t time.Time) tea.Msg { return logoTickMsg{} })
}

func (m LogoModel) Update(msg tea.Msg) (LogoModel, tea.Cmd) {
	switch msg.(type) {
	case logoTickMsg:
		if m.revealed < m.total {
			m.revealed += 2
			if m.revealed > m.total {
				m.revealed = m.total
			}
			return m, tea.Tick(30*time.Millisecond, func(t time.Time) tea.Msg { return logoTickMsg{} })
		}
		m.done = true
	}
	return m, nil
}

func (m LogoModel) View() string {
	if m.done {
		result := ""
		for _, l := range LogoLines {
			result += l + "\n"
		}
		return result
	}
	remaining := m.revealed
	result := ""
	for _, line := range LogoLines {
		runes := []rune(line)
		if remaining >= len(runes) {
			result += string(runes) + "\n"
			remaining -= len(runes)
		} else if remaining > 0 {
			result += string(runes[:remaining]) + "\n"
			remaining = 0
		} else {
			result += "\n"
		}
	}
	return result
}

func (m LogoModel) Done() bool { return m.done }
