package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ErrInputCancelled is returned when the user cancels the CSV path prompt.
var ErrInputCancelled = errors.New("input cancelled")

type inputModel struct {
	input     textinput.Model
	errMsg    string
	submitted bool
	cancelled bool
}

// PromptCSVPath displays an interactive text input for the user to
// enter a CSV file path, returning the sanitized path.
func PromptCSVPath() (string, error) {
	m := newInputModel()
	finalModel, err := tea.NewProgram(
		m,
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stdout),
	).Run()
	if err != nil {
		return "", fmt.Errorf("run input screen: %w", err)
	}

	fm, ok := finalModel.(inputModel)
	if !ok {
		return "", fmt.Errorf("unexpected input model type %T", finalModel)
	}

	if fm.cancelled {
		return "", ErrInputCancelled
	}

	path := sanitizeCSVPath(fm.input.Value())
	if path == "" {
		return "", fmt.Errorf("csv path cannot be empty")
	}

	return path, nil
}

func newInputModel() inputModel {
	ti := textinput.New()
	ti.Placeholder = "./usage.example.csv"
	ti.Prompt = "CSV Path: "
	ti.Focus()

	return inputModel{
		input: ti,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			if sanitizeCSVPath(m.input.Value()) == "" {
				m.errMsg = "Please type a CSV path."
				return m, nil
			}
			m.submitted = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	var content strings.Builder
	content.WriteString(TitleStyle.Render("Cursor Usage Analyzer"))
	content.WriteString("\n")
	content.WriteString(MutedStyle.Render("Enter the CSV path and press Enter"))
	content.WriteString("\n\n")
	content.WriteString(m.input.View())
	content.WriteString("\n")
	content.WriteString(MutedStyle.Render("Esc / Ctrl+C to cancel"))

	if m.errMsg != "" {
		content.WriteString("\n\n")
		content.WriteString(ErrorStyle.Render(m.errMsg))
	}

	return "\n" + PanelStyle.Render(content.String()) + "\n"
}

func sanitizeCSVPath(raw string) string {
	path := strings.TrimSpace(raw)
	path = strings.ReplaceAll(path, "\n", "")
	path = strings.ReplaceAll(path, "\r", "")

	for len(path) >= 2 {
		if (path[0] == '"' && path[len(path)-1] == '"') || (path[0] == '\'' && path[len(path)-1] == '\'') {
			path = strings.TrimSpace(path[1 : len(path)-1])
			continue
		}
		break
	}

	path = expandHome(path)
	return path
}

func expandHome(path string) string {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return home
	}

	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}

	return path
}
