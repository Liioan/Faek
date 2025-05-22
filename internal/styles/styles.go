package styles

import "github.com/charmbracelet/lipgloss"

const (
	White     = lipgloss.Color("15")
	Primary   = lipgloss.Color("85")
	Secondary = lipgloss.Color("246")
	Danger    = lipgloss.Color("9")
	Warning   = lipgloss.Color("215")
	Disabled  = lipgloss.Color("240")
)

var (
	//- model styles
	OutputStyle      = lipgloss.NewStyle().Bold(true).MarginLeft(1).MarginTop(1).Foreground(Secondary)
	OutputTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(Primary)
	QuitStyle        = lipgloss.NewStyle().Foreground(Danger).Bold(true)
	HighlightStyle   = lipgloss.NewStyle().Foreground(Primary).Bold(true)

	//- step styles
	TitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(Primary).MarginLeft(2)
	AnswerStyle = lipgloss.NewStyle().MarginLeft(2).MarginTop(1)

	//- list styles
	ListTitleStyle       = lipgloss.NewStyle().Bold(true).Foreground(White)
	ItemStyle            = lipgloss.NewStyle().MarginLeft(2)
	SelectedItemStyle    = lipgloss.NewStyle().Foreground(Primary)
	DestructiveItemStyle = lipgloss.NewStyle().Foreground(Danger)
	WarningItemStyle     = lipgloss.NewStyle().Foreground(Warning)
	DisabledItemStyle    = lipgloss.NewStyle().Foreground(Disabled)

	//- table styles
	TableHeaderStyle  = lipgloss.NewStyle().Foreground(Primary).Bold(true).Width(16).PaddingLeft(1)
	TableEvenRowStyle = lipgloss.NewStyle().Foreground(Secondary).PaddingLeft(1)
	TableOddRowStyle  = lipgloss.NewStyle().Foreground(White).PaddingLeft(1)
)
