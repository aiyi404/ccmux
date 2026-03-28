package styles

import "github.com/charmbracelet/lipgloss"

var (
	Purple   = lipgloss.Color("99")
	Pink     = lipgloss.Color("212")
	Lavender = lipgloss.Color("141")
	Gray     = lipgloss.Color("240")
	Green    = lipgloss.Color("82")
	Red      = lipgloss.Color("196")
	Yellow   = lipgloss.Color("220")
	Cyan     = lipgloss.Color("86")

	Bold = lipgloss.NewStyle().Bold(true)
	Dim  = lipgloss.NewStyle().Faint(true)

	LogoStyle    = lipgloss.NewStyle().Foreground(Purple).Bold(true).Align(lipgloss.Center)
	VersionStyle = lipgloss.NewStyle().Foreground(Gray).Align(lipgloss.Center)

	StatusBox = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Purple).
			Foreground(Pink).
			Padding(1, 3)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).Foreground(Lavender).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Gray).
			Padding(0, 4).Align(lipgloss.Center)

	MenuItemStyle     = lipgloss.NewStyle().PaddingLeft(2)
	MenuSelectedStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(Pink).Bold(true)

	SuccessStyle = lipgloss.NewStyle().Foreground(Green)
	ErrorStyle   = lipgloss.NewStyle().Foreground(Red)
	WarnStyle    = lipgloss.NewStyle().Foreground(Yellow)

	TableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(Lavender)
	TableCellStyle   = lipgloss.NewStyle()
	TableActiveRow   = lipgloss.NewStyle().Foreground(Green).Bold(true)
)
