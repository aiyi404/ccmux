package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/farion1231/ccmux/internal/cli/tui/styles"
	"github.com/farion1231/ccmux/internal/services"
)

type ProviderTable struct {
	table    table.Model
	selected string
}

type ProviderSelectedMsg struct{ Name string }

func NewProviderTable(providers []services.Provider, currentName string, width int) ProviderTable {
	columns := []table.Column{
		{Title: " ", Width: 2},
		{Title: "NAME", Width: 18},
		{Title: "BASE_URL", Width: 36},
		{Title: "MODEL", Width: 20},
	}
	var rows []table.Row
	for _, p := range providers {
		marker := " "
		if p.Name == currentName {
			marker = "→"
		}
		url := strings.TrimPrefix(strings.TrimPrefix(p.Env["ANTHROPIC_BASE_URL"], "http://"), "https://")
		rows = append(rows, table.Row{marker, p.Name, url, p.Env["ANTHROPIC_MODEL"]})
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Gray).
		BorderBottom(true).
		Bold(true).
		Foreground(styles.Lavender)
	s.Selected = s.Selected.Foreground(styles.Pink).Bold(true)
	t.SetStyles(s)
	return ProviderTable{table: t}
}

func (m ProviderTable) Init() tea.Cmd { return nil }

func (m ProviderTable) Update(msg tea.Msg) (ProviderTable, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			row := m.table.SelectedRow()
			if row != nil {
				m.selected = row[1]
				return m, func() tea.Msg { return ProviderSelectedMsg{Name: m.selected} }
			}
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ProviderTable) View() string { return m.table.View() }

func (m ProviderTable) SelectedName() string {
	row := m.table.SelectedRow()
	if row != nil {
		return row[1]
	}
	return ""
}
