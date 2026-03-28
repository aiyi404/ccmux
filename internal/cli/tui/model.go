package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/farion1231/ccmux/internal/cli/commands"
	"github.com/farion1231/ccmux/internal/cli/i18n"
	"github.com/farion1231/ccmux/internal/cli/tui/animations"
	"github.com/farion1231/ccmux/internal/cli/tui/pages"
	"github.com/farion1231/ccmux/internal/cli/tui/styles"
	"github.com/farion1231/ccmux/internal/config"
	"github.com/farion1231/ccmux/internal/store"
)

type Page int

const (
	PageLogo Page = iota
	PageHome
	PageProviders
	PageSettings
)

const sidebarWidth = 16

// ExecSignal holds the result for post-TUI syscall.Exec.
type ExecSignal struct {
	Result *commands.ExecResult
}

type Model struct {
	state  *store.AppState
	page   Page
	width  int
	height int

	logo     animations.LogoModel
	feedback *animations.FeedbackModel

	focusSidebar bool
	navCursor    int
	navItems     []string

	providers pages.ProvidersModel
	settings  pages.SettingsModel

	execSignal *ExecSignal
	quitting   bool
}
func NewModel(state *store.AppState) Model {
	return Model{
		state:        state,
		page:         PageLogo,
		logo:         animations.NewLogoModel(),
		focusSidebar: true,
		navCursor:    0,
		navItems:     []string{"Home", i18n.T("providers"), i18n.T("settings"), i18n.T("exit")},
		providers:    pages.NewProviders(state),
		settings:     pages.NewSettings(state),
	}
}

func (m Model) Init() tea.Cmd {
	return m.logo.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Feedback animation
	if m.feedback != nil {
		fb, cmd := m.feedback.Update(msg)
		m.feedback = &fb
		if _, ok := msg.(animations.FeedbackDoneMsg); ok {
			m.feedback = nil
		}
		if cmd != nil {
			return m, cmd
		}
	}

	var cmds []tea.Cmd

	// Logo page
	if m.page == PageLogo {
		logo, cmd := m.logo.Update(msg)
		m.logo = logo
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		if m.logo.Done() {
			m.page = PageHome
			m.focusSidebar = true
		}
		return m, tea.Batch(cmds...)
	}
	// Handle keyboard
	if msg, ok := msg.(tea.KeyMsg); ok {
		if m.focusSidebar {
			switch msg.String() {
			case "up", "k":
				if m.navCursor > 0 {
					m.navCursor--
				}
			case "down", "j":
				if m.navCursor < len(m.navItems)-1 {
					m.navCursor++
				}
			case "enter":
				switch m.navCursor {
				case 0:
					m.page = PageHome
				case 1:
					m.providers.Refresh()
					m.page = PageProviders
					m.focusSidebar = false
				case 2:
					m.settings = pages.NewSettings(m.state)
					m.page = PageSettings
					m.focusSidebar = false
				case 3:
					m.quitting = true
					return m, tea.Quit
				}
			case "right", "l", "tab":
				if m.page != PageHome {
					m.focusSidebar = false
				}
			case "q":
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}

		// Content has focus
		switch msg.String() {
		case "left", "h":
			m.focusSidebar = true
			return m, nil
		case "tab":
			m.focusSidebar = true
			return m, nil
		}
	}

	// Delegate to active page
	switch m.page {
	case PageProviders:
		pp, cmd := m.providers.Update(msg)
		m.providers = pp
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case PageSettings:
		ss, cmd := m.settings.Update(msg)
		m.settings = ss
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	// Handle page messages
	switch msg := msg.(type) {
	case pages.ProviderActionMsg:
		switch msg.Action {
		case pages.ProviderActionViewDetails:
			info := pages.ShowDetails(m.state, msg.Name)
			m.providers.SetInfo(info)
		case pages.ProviderActionSwitchGlobal:
			err := commands.RunSwitch(m.state, msg.Name)
			if err != nil {
				fb := animations.NewFeedback(err.Error(), animations.FeedbackError)
				m.feedback = &fb
				cmds = append(cmds, fb.Init())
			} else {
				fb := animations.NewFeedback("switched to '"+msg.Name+"'", animations.FeedbackSuccess)
				m.feedback = &fb
				cmds = append(cmds, fb.Init())
			}
			m.providers.Refresh()
		case pages.ProviderActionUseSession:
			result, err := commands.BuildExecResult(m.state, msg.Name, nil)
			if err != nil {
				fb := animations.NewFeedback(err.Error(), animations.FeedbackError)
				m.feedback = &fb
				cmds = append(cmds, fb.Init())
			} else {
				m.execSignal = &ExecSignal{Result: result}
				m.quitting = true
				return m, tea.Quit
			}
		case pages.ProviderActionEdit:
			p, err := m.state.Service.GetByName(msg.Name)
			if err == nil {
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vi"
				}
				if svc, ok := m.state.Service.(interface{ ProfilePath(string) string }); ok {
					path := svc.ProfilePath(p.Name)
					c := exec.Command(editor, path)
					return m, tea.ExecProcess(c, func(err error) tea.Msg {
						return pages.GoBackMsg{}
					})
				}
			}
		case pages.ProviderActionRemove:
			err := m.state.Service.Remove(msg.Name)
			if err != nil {
				fb := animations.NewFeedback(err.Error(), animations.FeedbackError)
				m.feedback = &fb
				cmds = append(cmds, fb.Init())
			} else {
				fb := animations.NewFeedback("removed '"+msg.Name+"'", animations.FeedbackSuccess)
				m.feedback = &fb
				cmds = append(cmds, fb.Init())
			}
			m.providers.Refresh()
		case pages.ProviderActionImport:
			err := m.state.Service.Import("")
			if err != nil {
				fb := animations.NewFeedback(err.Error(), animations.FeedbackError)
				m.feedback = &fb
				cmds = append(cmds, fb.Init())
			} else {
				fb := animations.NewFeedback("imported", animations.FeedbackSuccess)
				m.feedback = &fb
				cmds = append(cmds, fb.Init())
			}
			m.providers.Refresh()
		}

	case pages.GoBackMsg:
		m.focusSidebar = true

	case pages.EditConfigMsg:
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		c := exec.Command(editor, config.CCCConfig)
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			return pages.GoBackMsg{}
		})
	}

	return m, tea.Batch(cmds...)
}
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Logo page — full screen
	if m.page == PageLogo {
		logo := styles.LogoStyle.Render(m.logo.View())
		version := styles.VersionStyle.Render("v0.2.0")
		return "\n\n\n" + logo + "\n" + version
	}

	w := m.width
	h := m.height
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}

	// Header
	header := m.renderHeader(w)

	// Footer
	footer := m.renderFooter(w)

	// Body height = total - header(3) - footer(1) - borders
	bodyHeight := h - 4
	if bodyHeight < 5 {
		bodyHeight = 5
	}

	// Sidebar
	sidebar := m.renderSidebar(bodyHeight)

	// Content
	contentWidth := w - lipgloss.Width(sidebar) - 2
	if contentWidth < 20 {
		contentWidth = 20
	}
	content := m.renderContent(contentWidth, bodyHeight)

	// Compose
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
	full := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)

	// Overlay feedback
	if m.feedback != nil {
		full += "\n" + m.feedback.View()
	}

	return full
}
func (m Model) renderHeader(width int) string {
	title := styles.HeaderAccent.Render(" ccc ") + styles.HeaderStyle.Render(" v0.2.0")

	providerName := i18n.T("none")
	if p, err := m.state.Service.GetCurrent(); err == nil {
		providerName = p.Name
	}
	right := styles.HeaderStyle.Render(
		fmt.Sprintf("Provider: %s  Mode: %s",
			styles.HeaderAccent.Render(providerName),
			m.state.Mode),
	)

	gap := width - lipgloss.Width(title) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}

	return styles.HeaderStyle.Width(width).Render(
		title + strings.Repeat(" ", gap) + right,
	)
}

func (m Model) renderSidebar(height int) string {
	var items []string
	icons := []string{"🏠", "🔌", "⚙ ", "🚪"}

	for i, name := range m.navItems {
		icon := icons[i]
		label := icon + " " + name
		if i == m.navCursor {
			items = append(items, styles.NavSelectedStyle.Render(label))
		} else {
			items = append(items, styles.NavItemStyle.Render(label))
		}
	}

	nav := strings.Join(items, "\n")

	// Pad to fill height
	lines := len(m.navItems)
	for lines < height-2 {
		nav += "\n"
		lines++
	}

	sidebarStyle := styles.SidebarStyle
	if m.focusSidebar {
		sidebarStyle = styles.SidebarFocusedStyle
	}

	return sidebarStyle.Width(sidebarWidth).Height(height - 2).Render(nav)
}
func (m Model) renderContent(width, height int) string {
	var content string

	switch m.page {
	case PageHome:
		content = m.renderHome()
	case PageProviders:
		content = m.providers.View()
	case PageSettings:
		content = m.settings.View()
	}

	contentStyle := styles.ContentStyle
	if !m.focusSidebar {
		contentStyle = styles.ContentFocusedStyle
	}

	return contentStyle.Width(width).Height(height - 2).Render(content)
}

func (m Model) renderHome() string {
	s := styles.TitleStyle.Render("Home") + "\n\n"

	p, err := m.state.Service.GetCurrent()
	if err != nil {
		s += styles.Dim.Render("No active provider. Use Providers to switch.")
		return s
	}

	s += styles.LabelStyle.Render("Provider") + styles.ValueStyle.Render(p.Name) + "\n"
	s += styles.LabelStyle.Render("Base URL") + styles.ValueStyle.Render(p.Env["ANTHROPIC_BASE_URL"]) + "\n"
	s += styles.LabelStyle.Render("Model") + styles.ValueStyle.Render(p.Env["ANTHROPIC_MODEL"]) + "\n"
	if p.ModelAlias != "" {
		s += styles.LabelStyle.Render("Model Alias") + styles.ValueStyle.Render(p.ModelAlias) + "\n"
	}
	s += "\n"

	// Provider count
	providers, _ := m.state.Service.List()
	s += styles.LabelStyle.Render("Providers") + styles.ValueStyle.Render(fmt.Sprintf("%d", len(providers))) + "\n"
	s += styles.LabelStyle.Render("Mode") + styles.ValueStyle.Render(m.state.Mode) + "\n"
	s += styles.LabelStyle.Render("Language") + styles.ValueStyle.Render(m.state.Lang) + "\n"

	return s
}

func (m Model) renderFooter(width int) string {
	var keys []string

	if m.focusSidebar {
		keys = []string{
			styles.FooterKeyStyle.Render("↑↓") + styles.FooterDescStyle.Render(" Navigate"),
			styles.FooterKeyStyle.Render("Enter") + styles.FooterDescStyle.Render(" Select"),
			styles.FooterKeyStyle.Render("→/Tab") + styles.FooterDescStyle.Render(" Content"),
			styles.FooterKeyStyle.Render("q") + styles.FooterDescStyle.Render(" Quit"),
		}
	} else {
		keys = []string{
			styles.FooterKeyStyle.Render("←/Tab") + styles.FooterDescStyle.Render(" Menu"),
			styles.FooterKeyStyle.Render("↑↓") + styles.FooterDescStyle.Render(" Navigate"),
			styles.FooterKeyStyle.Render("Enter") + styles.FooterDescStyle.Render(" Select"),
			styles.FooterKeyStyle.Render("Esc") + styles.FooterDescStyle.Render(" Back"),
		}
	}

	return styles.FooterStyle.Width(width).Render("  " + strings.Join(keys, "  │  "))
}

// GetExecSignal returns the signal for post-TUI syscall.Exec.
func (m Model) GetExecSignal() *ExecSignal {
	return m.execSignal
}
