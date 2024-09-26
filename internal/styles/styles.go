package styles

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#44cbca")).MarginLeft(2)
	AnswerStyle       = lipgloss.NewStyle().MarginLeft(2)
	OutputStyle       = lipgloss.NewStyle().Bold(true)
	QuitStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	HelpHeaderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#44cbca")).Bold(true)
	HelpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)
