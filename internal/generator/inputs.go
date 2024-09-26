package generator

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liioan/faek/internal/styles"
)

type InputComponent interface {
	Value() string
	setValue(string)
	Update(tea.Msg) (InputComponent, tea.Cmd)
	View() string
}

//------- TextInput component -------

type textInputField struct {
	textInput textinput.Model
}

func newTextInputField(placeholder string) *textInputField {
	a := textInputField{}

	model := textinput.New()
	model.Placeholder = placeholder

	model.Focus()

	a.textInput = model
	return &a
}

func (a *textInputField) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *textInputField) Init() tea.Cmd {
	return nil
}

func (a *textInputField) Update(msg tea.Msg) (InputComponent, tea.Cmd) {
	var cmd tea.Cmd
	a.textInput, cmd = a.textInput.Update(msg)
	return a, cmd
}

func (a *textInputField) View() string {
	return a.textInput.View()
}

func (a *textInputField) Focus() tea.Cmd {
	return a.textInput.Focus()
}

func (a *textInputField) SetValue(s string) {
	a.textInput.SetValue(s)
}

func (a *textInputField) Blur() tea.Msg {
	return a.textInput.Blur
}

func (a *textInputField) Value() string {
	return a.textInput.Value()
}

func (a *textInputField) SelectedItem() list.Item {
	return nil
}

func (a *textInputField) setValue(v string) {
	a.textInput.SetValue(v)
}

//------- ListInput component -------

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := styles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type listInputField struct {
	listInput list.Model
}

func newListInputField(items []list.Item, delegate list.ItemDelegate, width, height int, title string) *listInputField {
	a := listInputField{}

	model := list.New(items, delegate, width, height)
	model.Title = title
	model.SetShowStatusBar(false)
	model.SetFilteringEnabled(false)
	model.SetShowHelp(false)
	model.Styles.Title = styles.TitleStyle

	a.listInput = model
	return &a
}

func (a *listInputField) Init() tea.Cmd {
	return nil
}

func (a *listInputField) Update(msg tea.Msg) (InputComponent, tea.Cmd) {
	var cmd tea.Cmd
	a.listInput, cmd = a.listInput.Update(msg)
	return a, cmd
}

func (a *listInputField) View() string {
	return a.listInput.View()
}

func (a *listInputField) Value() string {
	i, ok := a.SelectedItem().(item)
	if ok {
		return string(i)
	}
	return ""
}

func (a *listInputField) SelectedItem() list.Item {
	return a.listInput.SelectedItem()
}

func (a *listInputField) Focus() tea.Cmd {
	return nil
}

func (a *listInputField) Blur() tea.Msg {
	return nil
}

func (a *listInputField) setValue(v string) {
	searchedItem := item(v)
	searchedItemIndex := 0
	for i, item := range a.listInput.Items() {
		if item == searchedItem {
			searchedItemIndex = i
		}
	}
	a.listInput.Select(searchedItemIndex)
}
