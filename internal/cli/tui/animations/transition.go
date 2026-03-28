package animations

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
)

type TransitionModel struct {
	spring harmonica.Spring
	pos    float64
	vel    float64
	active bool
	done   bool
}

type transitionTickMsg struct{}

type TransitionDoneMsg struct{}

func NewTransition() TransitionModel {
	return TransitionModel{
		spring: harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5),
		done:   true,
	}
}

func (m *TransitionModel) Start() tea.Cmd {
	m.pos = 0
	m.vel = 0
	m.active = true
	m.done = false
	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg { return transitionTickMsg{} })
}

func (m TransitionModel) Update(msg tea.Msg) (TransitionModel, tea.Cmd) {
	switch msg.(type) {
	case transitionTickMsg:
		if !m.active {
			return m, nil
		}
		m.pos, m.vel = m.spring.Update(m.pos, m.vel, 1.0)
		if m.pos > 0.99 {
			m.pos = 1.0
			m.active = false
			m.done = true
			return m, func() tea.Msg { return TransitionDoneMsg{} }
		}
		return m, tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg { return transitionTickMsg{} })
	}
	return m, nil
}

func (m TransitionModel) Offset(width int) int {
	return int(float64(width) * (1.0 - m.pos))
}

func (m TransitionModel) Done() bool { return m.done }
