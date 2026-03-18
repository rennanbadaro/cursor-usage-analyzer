package ui

import "github.com/charmbracelet/lipgloss"

// AccentColor is the primary accent used for titles and spinners.
var AccentColor = lipgloss.Color("205")

// TitleStyle renders bold section headings.
var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(AccentColor)

// MutedStyle renders secondary / hint text.
var MutedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241"))

// ErrorStyle renders error messages.
var ErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("196")).
	Bold(true)

// LabelStyle renders row labels in the report.
var LabelStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("252"))

// ValueStyle renders numeric values in the report.
var ValueStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("230"))

// PanelStyle renders the outer bordered panel.
var PanelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63")).
	Padding(1, 2)
