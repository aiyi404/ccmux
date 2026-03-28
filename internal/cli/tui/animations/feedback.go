package animations

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/farion1231/ccmux/internal/cli/tui/styles"
)

type FeedbackType int

const (
	FeedbackSuccess FeedbackType = iota
	FeedbackError
)

type FeedbackModel struct {
	Message  string
	Type     FeedbackType
	visible  bool
	ticks    int
	maxTicks int
}

type feedbackTickMsg struct{}

type FeedbackDoneMsg struct{}

func NewFeedback(msg string, ft FeedbackType) FeedbackModel {
	return FeedbackModel{Message: msg, Type: ft, visible: true, maxTicks: 6}
}

func (m FeedbackModel) Init() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg { return feedbackTickMsg{} })
}

func (m FeedbackModel) Update(msg tea.Msg) (FeedbackModel, tea.Cmd) {
	switch msg.(type) {
	case feedbackTickMsg:
		m.ticks++
		m.visible = !m.visible
		if m.ticks >= m.maxTicks {
			m.visible = false
			return m, func() tea.Msg { return FeedbackDoneMsg{} }
		}
		return m, tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg { return feedbackTickMsg{} })
	}
	return m, nil
}

func (m FeedbackModel) View() string {
	if !m.visible {
		return ""
	}
	var style lipgloss.Style
	prefix := "✓ "
	if m.Type == FeedbackError {
		style = styles.ErrorStyle
		prefix = "✗ "
	} else {
		style = styles.SuccessStyle
	}
	return style.Render(prefix + m.Message)
}
