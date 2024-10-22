package styles

import "github.com/charmbracelet/lipgloss"

var (
	//- model styles
	OutputStyle     = lipgloss.NewStyle().Bold(true)
	QuitStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	HelpHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#44cbca")).Bold(true)
	HelpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))

	//- step styles
	TitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#44cbca")).MarginLeft(2)
	AnswerStyle = lipgloss.NewStyle().MarginLeft(2).MarginTop(1)

	//- list styles
	ListTitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fff"))
	ItemStyle         = lipgloss.NewStyle().MarginLeft(2)
	SelectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))
)
