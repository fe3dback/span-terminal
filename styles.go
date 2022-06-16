package terminal

import "github.com/charmbracelet/lipgloss"

const colorGreen = "2"
const colorYellow = "3"
const colorCyan = "4"
const colorPurple = "5"

var styleStatusDone = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorGreen))

var styleStatusActive = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(colorYellow))

var styleHeader = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(colorCyan))

var styleLogs = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorPurple))
