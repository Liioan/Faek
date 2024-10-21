package generator

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/liioan/faek/internal/constance"
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
	"strset":    "strSet",
}

const (
	Text inputMode = "text"
	List inputMode = "list"
)

var typesWithOptions = map[string]string{
	"date": "Choose a date format: ",
	"img":  "Choose a size of the image: ",
}

type Option string

const (
	HorizontalImg     Option = "300x500"
	VerticalImg       Option = "500x300"
	ProfilePictureImg Option = "100x100"
	ArticleImg        Option = "600x400"
	Banner            Option = "600x240"
	Custom            Option = "custom"
)

var imgOptions = map[string]Option{
	"Horizontal (default) 300x500": HorizontalImg,
	"Vertical 500x300":             VerticalImg,
	"Profile Picture 100x100":      ProfilePictureImg,
	"Article 600x400":              ArticleImg,
	"Banner 600x240":               Banner,
	"Custom":                       Custom,
}

const (
	DateTime   Option = "dateTime"
	Timestamp  Option = "timestamp"
	Day        Option = "day"
	Month      Option = "month"
	Year       Option = "year"
	DateObject Option = "object"
)

var dateOptions = map[string]Option{
	"dateTime: e.g. 27.02.2024":  DateTime,
	"timestamp: e.g. 1718051654": Timestamp,
	"day: 0-31":                  Day,
	"month: 0-12":                Month,
	"year: current year":         Year,
	"object: new Date()":         DateObject,
}

func newOptionsInput(fieldType string, instruction string) *listInputField {
	options := map[string]Option{}
	switch fieldType {
	case "date":
		options = dateOptions
	case "img":
		options = imgOptions
	}
	l := []list.Item{}
	for k := range options {
		l = append(l, item(k))
	}
	return newListInputField(l, itemDelegate{}, constance.DefaultWidth, constance.ListHeight, instruction)
}

type activeInput struct {
	input InputComponent
	mode  inputMode
}

type Field struct {
	name      string
	fieldType string
	option    Option
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

func NewListStep(instruction string, repeats bool, options []list.Item) *Step {
	i := activeInput{
		input: newListInputField(options, itemDelegate{}, constance.DefaultWidth, constance.ListHeight, instruction),
		mode:  List,
	}

	s := Step{Instruction: instruction, Repeats: repeats, StepInput: i}
	return &s
}

func NewTextStep(instruction, placeholder string, repeats bool) *Step {
	i := activeInput{
		input: newTextInputField(placeholder),
		mode:  Text,
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
		test := m.Steps
		for _, step := range test {
			if len(step.Answer.fields) > 0 {
				output += "\nFields: \n"
				for _, field := range step.Answer.fields {
					output += field.name + " " + field.fieldType
					if len(field.option) > 0 {
						output += " " + string(field.option)
					}

					output += "\n"
				}
				output += "\n"
				continue
			}
			output += step.Answer.text + "\n"
			utils.LogToFile(output)
		}

		return styles.QuitStyle.Render()
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
	if m.ActiveInput.mode == List {
		fieldsLen := len(m.Steps[m.Index].Answer.fields)
		currentField := &m.Steps[m.Index].Answer.fields[fieldsLen-1]

		switch currentField.fieldType {
		case "date":
			currentField.option = dateOptions[userInput]
		case "img":
			currentField.option = imgOptions[userInput]
		}
		m.ActiveInput = m.Steps[m.Index].StepInput
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
						optionsInput := newOptionsInput(k, v)
						m.ActiveInput.input = optionsInput
						m.ActiveInput.mode = List
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
