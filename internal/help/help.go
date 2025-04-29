package help

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
	s "github.com/liioan/faek/internal/styles"
)

func ShowHelpScreen() {
	text := fmt.Sprintf("%s\n\n", s.TitleStyle.Render("----- Help -----"))
	text += fmt.Sprintf("usage: %s\n\n", s.OutputStyle.Margin(0).Render("faek [mode: -c | -h | -d [-template=<val>] [-len=<val>]] [language: -ts | -js | -json] [output: -file | -terminal]"))

	text += "Available debug templates:\n"
	rows := [][]string{
		{"types", "contains all types, no custom type, no array name"},
		{"user", "simple user template with custom type and array name"},
		{"imgs", "contains all variants for img field"},
		{"dates", "contains all variants for date field"},
	}
	text += createTable(rows, []string{}).Render()

	text += "\n\n"
	text += "Available types:\n"
	rows = [][]string{
		{"string", "[length]", "lorem ipsum text with given length"},
		{"number", "[max] | [min max]", "random number within given range"},
		{"boolean", "", "true/false"},
		{"date", "", "date in given format"},
		{"img", "", "img with given size"},
		{"strSet", "[str...]", "random word from given set"},
	}
	text += createTable(rows, []string{"Type", "Options", "Value"}).Render()

	text += "\n\n"

	text += "Predefined string fields:\n\n"
	l := list.New(
		"name",
		"surname",
		"email",
		"title",
		"content",
		"author",
	).
		EnumeratorStyle(s.HighlightStyle.Margin(0).MarginLeft(1)).
		ItemStyle(s.OutputStyle.Margin(0).MarginLeft(1))

	text += l.String()
	text += "\n\n"

	fmt.Print(text)

}

func createTable(rows [][]string, headers []string) *table.Table {
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return lipgloss.NewStyle().Margin(0).MarginLeft(1)
			case col == 0:
				return s.HighlightStyle.Margin(0).MarginLeft(1)
			default:
				return s.OutputStyle.Margin(0).MarginLeft(1)
			}
		}).
		Headers(headers...).
		Rows(rows...)
	return t
}
