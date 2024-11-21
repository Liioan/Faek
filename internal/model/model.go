package model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/liioan/faek/internal/configuration"
	o "github.com/liioan/faek/internal/options"
	"github.com/liioan/faek/internal/styles"
	"github.com/liioan/faek/internal/utils"
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

	"strset": "strSet",
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
	option    o.Option
}

type Step struct {
	Instruction string
	StepInput   activeInput
	Answer      struct {
		text   string
		fields []Field
	}
	Repeats bool
}

func NewListStep(title, instruction string, repeats bool, optionSet o.OptionSet) *Step {
	utils.LogToDebug(string(optionSet))

	i := activeInput{
		input: newOptionsInput(optionSet, instruction),
		mode:  ListInput,
	}

	s := Step{Instruction: title, Repeats: repeats, StepInput: i}
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
				FileName:    strings.Trim(m.Steps[1].Answer.text, " "),
				Language:    m.Steps[2].Answer.text,
			}

			configuration.SaveUserSettings(&settings)

			output += styles.OutputTitleStyle.Render("Your preferences have been saved!")
			output += "\n"
			output += styles.OutputStyle.Render(fmt.Sprintf(" output style: %s\n filename: %s\n language: %s", settings.OutputStyle, settings.FileName, settings.Language))
		} else {
			generateOutput(m)
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

func checkAnswer(m *Model, current *Step, userInput string) {
	if userInput == "Custom" {
		m.ActiveInput = activeInput{input: newTextInputField("custom dimension e.g. 200x300"), mode: CustomInput}
		return
	}

	if m.ActiveInput.mode == ListInput || m.ActiveInput.mode == CustomInput {
		if m.Configuration {
			current.Answer.text = string(getOptionsValue("output", userInput))
			m.Next()
			return
		} else {
			fieldsLen := len(m.Steps[m.Index].Answer.fields)
			currentField := &m.Steps[m.Index].Answer.fields[fieldsLen-1]

			if m.ActiveInput.mode == CustomInput {
				currentField.option = o.Option(userInput)
			} else {
				currentField.option = getOptionsValue(o.OptionSet(currentField.fieldType), userInput)
			}
			m.ActiveInput = m.Steps[m.Index].StepInput
		}
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
			values := strings.Fields(strings.ToLower(userInput))
			if len(values) == 1 {
				return
			}
			if len(values) == 2 {
				for k, v := range typeConversion {
					if values[1] == k {
						values[1] = v
					}
				}
				if ValidTypesArray.contains(values[1]) {
					current.Answer.fields = append(current.Answer.fields, Field{name: values[0], fieldType: values[1]})
				}

				for k, v := range typesWithOptions {
					if k == values[1] {
						optionsInput := newOptionsInput(o.OptionSet(k), v)
						m.ActiveInput.input = optionsInput
						m.ActiveInput.mode = ListInput
						return
					}
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
			checkAnswer(&m, current, m.ActiveInput.input.Value())
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
