package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/aiyi404/ccmux/internal/cli/tui/styles"
)

type FormField struct {
	Label       string
	Placeholder string
	Required    bool
	Password    bool
}

type FormModel struct {
	fields  []FormField
	inputs  []textinput.Model
	focused int
	done    bool
}

type FormSubmitMsg struct{ Values map[string]string }

func NewForm(fields []FormField) FormModel {
	var inputs []textinput.Model
	for _, f := range fields {
		ti := textinput.New()
		ti.Placeholder = f.Placeholder
		ti.CharLimit = 256
		ti.Width = 50
		if f.Password {
			ti.EchoMode = textinput.EchoPassword
		}
		inputs = append(inputs, ti)
	}
	if len(inputs) > 0 {
		inputs[0].Focus()
	}
	return FormModel{fields: fields, inputs: inputs}
}

func (m FormModel) Init() tea.Cmd { return textinput.Blink }

func (m FormModel) Update(msg tea.Msg) (FormModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.focused = (m.focused + 1) % len(m.inputs)
			return m, m.focusCurrent()
		case "shift+tab", "up":
			m.focused = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
			return m, m.focusCurrent()
		case "enter":
			if m.focused == len(m.inputs)-1 {
				values := make(map[string]string)
				for i, f := range m.fields {
					values[f.Label] = m.inputs[i].Value()
				}
				m.done = true
				return m, func() tea.Msg { return FormSubmitMsg{Values: values} }
			}
			m.focused++
			return m, m.focusCurrent()
		}
	}
	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return m, cmd
}

func (m FormModel) focusCurrent() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focused {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m FormModel) View() string {
	s := ""
	for i, f := range m.fields {
		label := f.Label
		if f.Required {
			label += " *"
		}
		style := styles.Dim
		if i == m.focused {
			style = styles.Bold
		}
		s += style.Render(label) + "\n" + m.inputs[i].View() + "\n\n"
	}
	return s
}

func (m FormModel) Done() bool { return m.done }
