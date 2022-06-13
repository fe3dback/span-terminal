package v2

import "github.com/charmbracelet/lipgloss"

const colorGreen = "2"
const colorYellow = "3"
const colorGray = "#777777"

var styleStatusDone = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(colorGreen))

var styleStatusActive = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(colorYellow))

var styleStatusWait = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(colorGray))

var styleProgressActive = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(colorYellow))

var styleProgressWait = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(colorGray))

var styleLogs = lipgloss.NewStyle().
	Foreground(lipgloss.Color(colorGray))
