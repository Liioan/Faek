package model

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/liioan/faek/internal/configuration"
	"github.com/liioan/faek/internal/styles"

	v "github.com/liioan/faek/internal/variants"
	"github.com/muesli/reflow/wordwrap"
)

type inputMode string

type TypesArray []string

var ValidTypesArray = TypesArray{"string", "number", "boolean", "img", "strSet", "date"}

func (vt TypesArray) contains(t string) bool {
	for _, k := range ValidTypesArray {
		if t == k {
			return true
		}
	}
	return false
}

var typeConversion = map[string]string{
	"int":       "number",
	"float":     "number",
	"short":     "number",
	"str":       "string",
	"char":      "string",
	"bool":      "boolean",
	"stringSet": "strSet",
	"ss":        "strSet",
	"strs":      "strSet",
	"strset":    "strSet",
}

var typesWithOptions = map[string]string{
	"date": "Choose a date format: ",
	"img":  "Choose a size of the image: ",
}

const (
	TextInput   inputMode = "text"
	ListInput   inputMode = "list"
	CustomInput inputMode = "custom"
)

type activeInput struct {
	input InputComponent
	mode  inputMode
}

type Field struct {
	name      string
	fieldType string
	variant   v.Variant
	options   []string
}

type Step struct {
	Instruction string
	StepInput   activeInput
	Answer      struct {
		text   string
		fields []Field
	}
	OptionSet v.VariantSet
	Repeats   bool
}

func NewListStep(title, instruction string, repeats bool, optionSet v.VariantSet) *Step {
	i := activeInput{
		input: newVariantsInput(optionSet, instruction),
		mode:  ListInput,
	}

	s := Step{Instruction: title, Repeats: repeats, StepInput: i, OptionSet: optionSet}
	return &s
}

func NewTextStep(instruction, placeholder string, repeats bool) *Step {
	i := activeInput{
		input: newTextInputField(placeholder),
		mode:  TextInput,
	}

	s := Step{Instruction: instruction, Repeats: repeats, StepInput: i}
	return &s
}

type Model struct {
	Index         int
	Width         int
	Height        int
	Finished      bool
	Quitting      bool
	Configuration bool
	ActiveInput   activeInput
	Steps         []Step
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Next() {
	if m.Index < len(m.Steps)-1 {
		m.Index++
	} else if m.Index == len(m.Steps)-1 {
		m.Finished = true
	}
	m.ActiveInput = m.Steps[m.Index].StepInput
}

func (m Model) View() string {
	current := m.Steps[m.Index]

	if m.Width == 0 {
		return "Loading..."
	}

	if m.Finished {
		output := ""
		if m.Configuration {
			settings := configuration.Settings{
				OutputStyle: m.Steps[0].Answer.text,
				Language:    m.Steps[1].Answer.text,
				FileName:    strings.Trim(m.Steps[2].Answer.text, " "),
				Indent:      m.Steps[3].Answer.text,
			}

			configuration.SaveUserSettings(&settings)

			output += styles.OutputTitleStyle.Render("Your preferences have been saved!")
			output += "\n"

			rows := [][]string{
				{"Output style", settings.OutputStyle},
				{"Language", settings.Language},
				{"File name", settings.FileName},
				{"Indent", settings.Indent},
			}

			table := table.New().
				Border(lipgloss.NormalBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))).
				StyleFunc(func(row, col int) lipgloss.Style {
					switch {
					case row == 0:
						return styles.TableHeaderStyle
					case row%2 == 0:
						return styles.TableEvenRowStyle
					default:
						return styles.TableOddRowStyle
					}
				}).
				Headers("Setting", "Value").
				Rows(rows...)

			output += table.Render()
			output += styles.OutputStyle.Render("You can always change your settings by running ") + styles.HighlightStyle.Render("`>faek -c``")
		} else {
			output += generateOutput(&m)
		}

		output += styles.QuitStyle.Render("\n\npress q or ctrl+c to exit")

		return output
	}

	if m.Quitting {
		return styles.QuitStyle.Render("Quitting")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wordwrap.String(styles.TitleStyle.Render(current.Instruction), m.Width),
		styles.AnswerStyle.Render(m.ActiveInput.input.View()),
	)
}

func parseInput(m *Model, current *Step, userInput string) {
	userInput = strings.Trim(userInput, " ")

	if userInput == "Custom" {
		m.ActiveInput = activeInput{input: newTextInputField("custom dimension e.g. 200x300"), mode: CustomInput}
		return
	}

	if m.ActiveInput.mode == ListInput || m.ActiveInput.mode == CustomInput {
		if m.Configuration {
			current.Answer.text = string(getVariantsValue(current.OptionSet, userInput))
			m.Next()
			return
		} else {
			fieldsLen := len(m.Steps[m.Index].Answer.fields)
			currentField := &m.Steps[m.Index].Answer.fields[fieldsLen-1]

			if m.ActiveInput.mode == CustomInput {
				regex := regexp.MustCompile(`^\d+x\d+$`)
				if !regex.MatchString(userInput) {
					return
				}

				currentField.variant = v.Variant(userInput)
			} else {
				currentField.variant = getVariantsValue(v.VariantSet(currentField.fieldType), userInput)
			}
			m.ActiveInput = m.Steps[m.Index].StepInput
		}
		return
	}

	if current.Repeats {
		if userInput == "" {
			if len(current.Answer.fields) > 0 {
				m.Next()
				return
			} else {
				return
			}
		} else {
			stringFields := strings.Fields(strings.ToLower(userInput))
			l := len(stringFields)
			if l == 1 {
				return
			}
			if l >= 2 {
				for key, value := range typeConversion {
					if stringFields[1] == key {
						stringFields[1] = value
					}
				}
				if ValidTypesArray.contains(stringFields[1]) {
					current.Answer.fields = append(current.Answer.fields, Field{name: stringFields[0], fieldType: stringFields[1]})
				}

				for key, value := range typesWithOptions {
					if key == stringFields[1] {
						variantInput := newVariantsInput(v.VariantSet(key), value)
						m.ActiveInput.input = variantInput
						m.ActiveInput.mode = ListInput
						return
					}
				}

				if l > 2 {
					fieldsLen := len(m.Steps[m.Index].Answer.fields)
					currentField := &m.Steps[m.Index].Answer.fields[fieldsLen-1]
					currentField.options = append(currentField.options, stringFields[2:]...)
				}
			}
			return
		}
	}
	current.Answer.text = userInput
	m.Next()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.Steps[m.Index]
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "q":
			if m.Index == len(m.Steps)-1 {
				return m, tea.Quit
			}
		case "enter":
			parseInput(&m, current, m.ActiveInput.input.Value())
			m.ActiveInput.input.setValue("")
			return m, nil
		}
	}
	m.ActiveInput.input, cmd = m.ActiveInput.input.Update(msg)
	return m, cmd
}

func NewModel(steps []Step, configuration bool) *Model {
	m := Model{Steps: steps, Configuration: configuration, ActiveInput: steps[0].StepInput}
	return &m
}

// - debug
func NewDebugModel(steps []Step, template string, length int) *Model {
	m := Model{Steps: steps, Configuration: false, ActiveInput: steps[len(steps)-1].StepInput}
	m.Index = len(m.Steps) - 1
	m.Finished = true

	switch template {
	default:
		m.Steps[0].Answer.text = ""
		m.Steps[1].Answer.fields = []Field{
			{name: "str", fieldType: "string"},
			{name: "int", fieldType: "number", options: []string{"10000"}},
			{name: "bool", fieldType: "boolean"},
			{name: "date", fieldType: "date", variant: v.Timestamp},
			{name: "img", fieldType: "img", variant: v.HorizontalImg},
			{name: "strSet", fieldType: "strSet", options: []string{"a", "b", "a_b"}},
		}
		m.Steps[2].Answer.text = ""
		m.Steps[3].Answer.text = fmt.Sprint(length)
	case "user":
		m.Steps[0].Answer.text = "users"
		m.Steps[1].Answer.fields = []Field{
			{name: "name", fieldType: "string"},
			{name: "surname", fieldType: "string"},
			{name: "age", fieldType: "number", options: []string{"18", "100"}},
			{name: "email", fieldType: "string"},
			{name: "isAdmin", fieldType: "boolean"},
		}
		m.Steps[2].Answer.text = "User"
		m.Steps[3].Answer.text = fmt.Sprint(length)
	case "date":
		m.Steps[0].Answer.text = "dates"
		m.Steps[1].Answer.fields = []Field{
			{name: "dateTime", fieldType: "date", variant: v.DateTime},
			{name: "timestamp", fieldType: "date", variant: v.Timestamp},
			{name: "day", fieldType: "date", variant: v.Day},
			{name: "month", fieldType: "date", variant: v.Month},
			{name: "year", fieldType: "date", variant: v.Year},
			{name: "obj", fieldType: "date", variant: v.DateObject},
		}
		m.Steps[2].Answer.text = "Dates"
		m.Steps[3].Answer.text = fmt.Sprint(length)
	case "img":
		m.Steps[0].Answer.text = "images"
		m.Steps[1].Answer.fields = []Field{
			{name: "horizontal", fieldType: "img", variant: v.HorizontalImg},
			{name: "profile", fieldType: "img", variant: v.ProfilePictureImg},
			{name: "banner", fieldType: "img", variant: v.Banner},
			{name: "custom", fieldType: "img", variant: v.Variant("5x5")},
		}
		m.Steps[2].Answer.text = "Images"
		m.Steps[3].Answer.text = fmt.Sprint(length)
	}
	return &m
}
