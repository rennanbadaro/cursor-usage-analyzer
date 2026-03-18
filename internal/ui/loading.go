package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadingDoneMsg struct{}

type loadingModel struct {
	spinner  spinner.Model
	message  string
	duration time.Duration
}

// ShowLoading displays a spinner with the given message for the specified
// duration before automatically dismissing.
func ShowLoading(message string, duration time.Duration) error {
	m := newLoadingModel(message, duration)
	_, err := tea.NewProgram(
		m,
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stdout),
	).Run()
	if err != nil {
		return fmt.Errorf("run loading screen: %w", err)
	}
	return nil
}

func newLoadingModel(message string, duration time.Duration) loadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(AccentColor)

	return loadingModel{
		spinner:  s,
		message:  message,
		duration: duration,
	}
}

func (m loadingModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.Tick(m.duration, func(time.Time) tea.Msg {
			return loadingDoneMsg{}
		}),
	)
}

func (m loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case loadingDoneMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m loadingModel) View() string {
	var content strings.Builder
	content.WriteString(TitleStyle.Render("Analyzing Cursor usage"))
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("%s %s", m.spinner.View(), m.message))
	content.WriteString("\n")
	content.WriteString(MutedStyle.Render("Just a moment..."))

	return "\n" + PanelStyle.Render(content.String()) + "\n"
}
