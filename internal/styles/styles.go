package styles

import "github.com/charmbracelet/lipgloss"

const (
	white   = "#fff"
	primary = "#44cbca"
	danger  = "9"
)

var (
	//- model styles
	OutputStyle      = lipgloss.NewStyle().Bold(true).MarginLeft(1).MarginTop(1)
	OutputTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(primary))
	QuitStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(danger)).Bold(true)
	HelpHeaderStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(primary)).Bold(true)
	HelpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(white))
	HighlightStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(primary)).Bold(false)

	//- step styles
	TitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(primary)).MarginLeft(2)
	AnswerStyle = lipgloss.NewStyle().MarginLeft(2).MarginTop(1)

	//- list styles
	ListTitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(white))
	ItemStyle         = lipgloss.NewStyle().MarginLeft(2)
	SelectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(white))

	//- table styles
	TableHeaderStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(primary)).Bold(true).Width(16).PaddingLeft(1)
	TableEvenRowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(1)
	TableOddRowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("248")).PaddingLeft(1)
)
