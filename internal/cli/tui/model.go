package tui

import (
	"os"
	"os/exec"

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
	PageMainMenu
	PageProviders
	PageSettings
)

// ExecSignal holds the result for post-TUI syscall.Exec.
type ExecSignal struct {
	Result *commands.ExecResult
}

type Model struct {
	state      *store.AppState
	page       Page
	width      int
	height     int

	logo     animations.LogoModel
	feedback *animations.FeedbackModel

	mainMenu   pages.MainMenuModel
	providers  pages.ProvidersModel
	settings   pages.SettingsModel

	execSignal *ExecSignal
	quitting   bool
}

func NewModel(state *store.AppState) Model {
	return Model{
		state:      state,
		page:       PageLogo,
		logo:     animations.NewLogoModel(),
		mainMenu:   pages.NewMainMenu(),
		providers:  pages.NewProviders(state),
		settings:   pages.NewSettings(state),
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

	// Update feedback animation
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

	switch m.page {
	case PageLogo:
		logo, cmd := m.logo.Update(msg)
		m.logo = logo
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		if m.logo.Done() {
			m.page = PageMainMenu
		}
// PLACEHOLDER_MODEL_UPDATE_CONT
	case PageMainMenu:
		mm, cmd := m.mainMenu.Update(msg)
		m.mainMenu = mm
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
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

	// Handle page navigation messages
	switch msg := msg.(type) {
	case pages.MainMenuSelectedMsg:
		switch msg.Choice {
		case pages.MainMenuProviders:
			m.providers.Refresh()
			m.page = PageProviders
		case pages.MainMenuCurrent:
			p, err := m.state.Service.GetCurrent()
			info := ""
			if err != nil {
				info = styles.WarnStyle.Render("no active provider")
			} else {
				info = styles.SuccessStyle.Render("→ "+p.Name) + "\n"
				info += "  base_url: " + p.Env["ANTHROPIC_BASE_URL"] + "\n"
				info += "  model:    " + p.Env["ANTHROPIC_MODEL"]
			}
			m.providers.SetInfo(info)
			m.providers.Refresh()
			m.page = PageProviders
		case pages.MainMenuSettings:
			m.settings = pages.NewSettings(m.state)
			m.page = PageSettings
		case pages.MainMenuExit:
			m.quitting = true
			return m, tea.Quit
		}
// PLACEHOLDER_MODEL_ACTIONS
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
// PLACEHOLDER_MODEL_ACTIONS_END
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
		m.page = PageMainMenu
		m.mainMenu = pages.NewMainMenu()

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

	var content string

	switch m.page {
	case PageLogo:
		logo := styles.LogoStyle.Render(m.logo.View())
		version := styles.VersionStyle.Render("v0.2.0")
		content = "\n\n\n" + logo + "\n" + version
	case PageMainMenu:
		providerInfo := i18n.T("none")
		if p, err := m.state.Service.GetCurrent(); err == nil {
			providerInfo = p.Name
		}
		status := styles.StatusBox.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				i18n.T("provider")+"  :  "+providerInfo,
				i18n.T("mode")+"      :  "+m.state.Mode,
				i18n.T("app")+"       :  claude",
			),
		)
		logo := styles.LogoStyle.Render(
			animations.LogoLines[0] + "\n" +
				animations.LogoLines[1] + "\n" +
				animations.LogoLines[2] + "\n" +
				animations.LogoLines[3],
		)
		content = "\n" + logo + "\n\n" + status + "\n\n" + m.mainMenu.View()
	case PageProviders:
		content = "\n" + m.providers.View()
	case PageSettings:
		content = "\n" + m.settings.View()
	}

	// Overlay feedback animation
	if m.feedback != nil {
		content += "\n\n" + m.feedback.View()
	}

	return content
}

// GetExecSignal returns the signal for post-TUI syscall.Exec.
func (m Model) GetExecSignal() *ExecSignal {
	return m.execSignal
}
