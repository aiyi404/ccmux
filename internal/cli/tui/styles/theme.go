package styles

import "github.com/charmbracelet/lipgloss"

// Dracula palette
var (
	Background  = lipgloss.Color("#282a36")
	CurrentLine = lipgloss.Color("#44475a")
	Foreground  = lipgloss.Color("#f8f8f2")
	Comment     = lipgloss.Color("#6272a4")
	Cyan        = lipgloss.Color("#8be9fd")
	Green       = lipgloss.Color("#50fa7b")
	Orange      = lipgloss.Color("#ffb86c")
	Pink        = lipgloss.Color("#ff79c8")
	Purple      = lipgloss.Color("#bd93f9")
	Red         = lipgloss.Color("#ff5555")
	Yellow      = lipgloss.Color("#f1fa8c")

	// Backward compat aliases
	Lavender = Purple
	Gray     = Comment
)

var (
	Bold = lipgloss.NewStyle().Bold(true)
	Dim  = lipgloss.NewStyle().Foreground(Comment)

	// Logo
	LogoStyle    = lipgloss.NewStyle().Foreground(Cyan).Bold(true).Align(lipgloss.Center)
	VersionStyle = lipgloss.NewStyle().Foreground(Comment).Align(lipgloss.Center)

	// Header: full-width bar
	HeaderStyle = lipgloss.NewStyle().
			Background(CurrentLine).
			Foreground(Foreground).
			Bold(true).
			Padding(0, 1)

	HeaderAccent = lipgloss.NewStyle().
			Background(CurrentLine).
			Foreground(Cyan).
			Bold(true)

	// Sidebar
	SidebarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Comment).
			Padding(0, 1)

	SidebarFocusedStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Cyan).
				Padding(0, 1)

	// Nav items
	NavItemStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			PaddingLeft(1)

	NavSelectedStyle = lipgloss.NewStyle().
				Background(Cyan).
				Foreground(lipgloss.Color("#282a36")).
				Bold(true).
				PaddingLeft(1)

	// Content area
	ContentStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Comment).
			Padding(0, 1)

	ContentFocusedStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Cyan).
				Padding(0, 1)

	// Footer
	FooterStyle = lipgloss.NewStyle().
			Foreground(Comment)

	FooterKeyStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	FooterDescStyle = lipgloss.NewStyle().
			Foreground(Comment)

	// Title inside content
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan)

	// Menu items (for sub-menus inside content)
	MenuItemStyle     = lipgloss.NewStyle().PaddingLeft(1).Foreground(Foreground)
	MenuSelectedStyle = lipgloss.NewStyle().PaddingLeft(1).Background(Cyan).Foreground(lipgloss.Color("#282a36")).Bold(true)

	// Status
	SuccessStyle = lipgloss.NewStyle().Foreground(Green)
	ErrorStyle   = lipgloss.NewStyle().Foreground(Red)
	WarnStyle    = lipgloss.NewStyle().Foreground(Yellow)

	// Labels and values
	LabelStyle = lipgloss.NewStyle().Foreground(Comment).Width(16)
	ValueStyle = lipgloss.NewStyle().Foreground(Cyan)

	// Table
	TableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(Purple)
	TableCellStyle   = lipgloss.NewStyle().Foreground(Foreground)
	TableActiveRow   = lipgloss.NewStyle().Foreground(Green).Bold(true)

	// StatusBox kept for backward compat but restyled
	StatusBox = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Comment).
			Foreground(Foreground).
			Padding(0, 2)
)
