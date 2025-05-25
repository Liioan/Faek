package model

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"

	list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	table "github.com/charmbracelet/lipgloss/table"

	"github.com/muesli/reflow/wordwrap"

	// internal
	c "github.com/liioan/faek/internal/configuration"
	e "github.com/liioan/faek/internal/errors"
	"github.com/liioan/faek/internal/styles"
	v "github.com/liioan/faek/internal/variants"
)

type Override struct {
	Language v.Variant
	Output   v.Variant
	Export   v.Variant
}

// ------- model -------

type Model struct {
	Index       int
	Width       int
	Height      int
	Finished    bool
	Quitting    bool
	ActiveInput activeInput
	Steps       []Step
	Settings    c.Settings

	ConfigurationMode bool
	DebugMode         bool
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
		if m.ConfigurationMode {
			settings := c.Settings{
				OutputStyle: v.Variant(m.Steps[0].Answer.text),
				Language:    v.Variant(m.Steps[1].Answer.text),
				FileName:    strings.Trim(m.Steps[2].Answer.text, " "),
				Indent:      m.Steps[3].Answer.text,
				Export:      v.Variant(m.Steps[4].Answer.text),
			}

			c.SaveUserSettings(&settings)

			output += styles.OutputTitleStyle.Render("Your preferences have been saved!")
			output += "\n"

			rows := [][]string{
				{"Output style", string(settings.OutputStyle)},
				{"Language", string(settings.Language)},
				{"File name", settings.GetFullFileName()},
				{"Indent", settings.Indent},
				{"Export", string(settings.Export)},
			}

			table := table.New().
				Border(lipgloss.RoundedBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(styles.White)).
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
			text := m.generateOutput()

			settings := m.Settings

			output += styles.TitleStyle.Margin(0).Render("Output: \n")

			switch settings.OutputStyle {
			case v.File:
				file, err := os.Create(settings.GetFullFileName())
				if err != nil {
					log.Fatal(err)
				}
				file.Write([]byte(text))
				output += styles.OutputStyle.Render(fmt.Sprintf("Created output file: `%s`\n\n", settings.GetFullFileName()))
			case v.Terminal:
				if m.Settings.Language == v.JSON {
					file, err := os.Create(settings.GetFullFileName())
					if err != nil {
						log.Fatal(err)
					}
					file.Write([]byte(text))
					output += styles.OutputStyle.Render(fmt.Sprintf("JSON is not supported in terminal, created output file: `%s`\n\n", settings.GetFullFileName()))

				} else if len(strings.Split(text, "\n")) > m.Height {
					file, err := os.Create(settings.GetFullFileName())
					if err != nil {
						log.Fatal(err)
					}
					file.Write([]byte(text))
					output += styles.OutputStyle.Render(fmt.Sprintf("Output is too big for your terminal, created output file: `%s`\n\n", settings.GetFullFileName()))
				} else {
					output += styles.OutputStyle.Render(wordwrap.String(text, m.Width))
				}
			}
		}

		output += styles.QuitStyle.Render("\n\npress q or ctrl+c to exit")

		return output
	}

	if m.Quitting {
		return styles.QuitStyle.Render("Quitting")
	}

	instruction := ""
	if m.Index == 1 && !m.ConfigurationMode {
		instruction = current.AvailableInputs[current.InputIdx].instruction
	} else {
		instruction = current.StepInput.instruction
	}

	isNextStepInput := m.Index == 1 && current.InputIdx == NEXT_STEP_INPUT
	isReviewStepInput := m.Index == 2 && current.InputIdx == CONFIRM_OBJ_INPUT

	if (isNextStepInput || isReviewStepInput) && !m.ConfigurationMode {
		rows := [][]string{}

		for _, field := range m.Steps[1].Answer.fields {
			rows = append(rows, []string{field.name, field.fieldType, string(field.variant)})
		}

		return lipgloss.JoinVertical(
			lipgloss.Left,
			wordwrap.String(styles.TitleStyle.Render(instruction), m.Width),
			styles.AnswerStyle.Render(m.ActiveInput.input.View()),
			styles.AnswerStyle.Render("current properties: "),
			createTable(rows, []string{"Name", "Type", "Option"}).Render())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wordwrap.String(styles.TitleStyle.Render(current.StepInput.instruction), m.Width),
		styles.AnswerStyle.Render(m.ActiveInput.input.View()),
	)
}

const (
	NEXT_STEP_INPUT     = 0
	PROPERTY_NAME_INPUT = 1
	PROPERTY_TYPE_INPUT = 2
	STRING_DATA_INPUT   = 3
	DATE_VARIANT_INPUT  = 4
	IMG_VARIANT_INPUT   = 5
	CUSTOM_SIZE_INPUT   = 6
	RANGE_INPUT         = 7
	STRING_SET_INPUT    = 8
)

const (
	OPTIONAL_STEP_1 = 0
	OPTIONAL_STEP_2 = 1
)

const (
	CONFIRM_OBJ_INPUT = 0
	DELETE_PROP_INPUT = 1
)

func updateLastField(userInput v.Variant, current *Step, m *Model) {
	lastField := &current.Answer.fields[len(current.Answer.fields)-1]
	lastField.variant = userInput
	current.InputIdx = NEXT_STEP_INPUT
	m.ActiveInput = current.AvailableInputs[current.InputIdx]
}

func parseInput(m *Model, current *Step, userInput string) {

	userInput = strings.Trim(userInput, " ")

	switch current.StepType {
	case NormalStep:
		if len(userInput) == 0 && m.Index != 0 {
			return
		}
		if len(current.Variants) != 0 {
			current.Answer.text = string(getVariantsValue(current.Variants, userInput))
		} else {
			current.Answer.text = userInput
		}
		m.Next()
		return
	case PropStep:
		switch current.InputIdx {
		case PROPERTY_NAME_INPUT:
			for _, f := range current.Answer.fields {
				if f.name == userInput {
					return
				}
			}

			current.Answer.text = userInput
			current.InputIdx++
			m.ActiveInput = current.AvailableInputs[current.InputIdx]

		case PROPERTY_TYPE_INPUT:
			fieldName := current.Answer.text
			fieldType := userInput
			current.Answer.fields = append(current.Answer.fields, Field{name: fieldName, fieldType: fieldType})

			switch fieldType {
			case "string":
				current.InputIdx = STRING_DATA_INPUT
			case "number":
				current.InputIdx = RANGE_INPUT
			case "string enum":
				current.InputIdx = STRING_SET_INPUT
			case "date":
				current.InputIdx = DATE_VARIANT_INPUT
			case "img":
				current.InputIdx = IMG_VARIANT_INPUT
			default:
				current.InputIdx = NEXT_STEP_INPUT
			}

			m.ActiveInput = current.AvailableInputs[current.InputIdx]

		case NEXT_STEP_INPUT:
			if userInput == "no" {
				current.Answer.text = ""
				m.Next()
				return
			} else {
				current.InputIdx++
				m.ActiveInput = current.AvailableInputs[current.InputIdx]
			}

		case IMG_VARIANT_INPUT:
			if userInput == "Custom" {
				current.InputIdx = CUSTOM_SIZE_INPUT
				m.ActiveInput = current.AvailableInputs[current.InputIdx]
			} else {
				selectedVariant := getVariantsValue(v.ImgVariants, userInput)
				updateLastField(selectedVariant, current, m)
			}

		case CUSTOM_SIZE_INPUT:
			regex := regexp.MustCompile(`^\d+x\d+$`)
			if regex.MatchString(userInput) {
				updateLastField(v.Variant(userInput), current, m)
			}

		case RANGE_INPUT:
			regex := regexp.MustCompile(`^(\d+\s\d+|\d+)?$`)
			if regex.MatchString(userInput) {
				updateLastField(v.Variant(userInput), current, m)
			}

		case DATE_VARIANT_INPUT:
			selectedVariant := getVariantsValue(v.DateVariants, userInput)
			updateLastField(selectedVariant, current, m)
		default:
			updateLastField(v.Variant(userInput), current, m)
		}
	case OptionalStep:
		switch current.InputIdx {
		case OPTIONAL_STEP_1:
			if userInput == "no" {
				m.Next()
			} else {
				current.InputIdx++
				m.ActiveInput = current.AvailableInputs[current.InputIdx]
			}
		case OPTIONAL_STEP_2:
			if len(userInput) == 0 {
				return
			}
			current.Answer.text = userInput
			m.Next()
		}
	case EditStep:
		if current.InputIdx == DELETE_PROP_INPUT {

			propStep := &m.Steps[1]

			propIdx := 0
			for i, prop := range propStep.Answer.fields {
				if userInput == strings.Trim(fmt.Sprintf("%s %s %s", prop.name, prop.fieldType, prop.variant), " ") {
					propIdx = i
					break
				}
			}

			propStep.Answer.fields = slices.Delete(propStep.Answer.fields, propIdx, propIdx+1)
			current.InputIdx = 0
			m.ActiveInput = current.AvailableInputs[current.InputIdx]
			if len(propStep.Answer.fields) == 0 {
				m.Index--
				propStep.InputIdx = 1
				m.ActiveInput = propStep.AvailableInputs[propStep.InputIdx]
				return
			}

			return

		}

		switch userInput {
		case "confirm":
			m.Next()
		case "add prop":
			m.Index--
			m.Steps[1].InputIdx = 1
			m.ActiveInput = m.Steps[1].AvailableInputs[m.Steps[1].InputIdx]

		case "delete prop":
			current.InputIdx = DELETE_PROP_INPUT
			props := []list.Item{item("cancel")}
			for _, f := range m.Steps[1].Answer.fields {
				props = append(props, item(fmt.Sprintf("%s %s %s", f.name, f.fieldType, f.variant)))
			}

			i := newListInputField(props, itemDelegate{func(s ...string) *lipgloss.Style {
				if strings.Contains(s[0], "cancel") {
					return &styles.SelectedItemStyle
				}
				return &styles.DestructiveItemStyle
			}}, 50, 14, "Which prop should be deleted?")
			current.AvailableInputs[DELETE_PROP_INPUT].input = i
			m.ActiveInput = current.AvailableInputs[current.InputIdx]
		}

	}

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

			if m.Index < len(m.Steps)-1 {
				parseInput(&m, current, m.ActiveInput.input.Value())
				m.ActiveInput.input.setValue("")
				return m, nil
			} else {
				parseInput(&m, current, m.ActiveInput.input.Value())
				m.Finished = true
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}

	if m.Finished {
		return m, cmd
	}

	m.ActiveInput.input, cmd = m.ActiveInput.input.Update(msg)
	return m, cmd
}

func NewModel(steps []Step, configMode bool, override Override) (*Model, error) {
	settings, err := c.GetUserSettings()
	if err != nil && !configMode {
		return nil, errors.New(e.SettingsUnavailable)
	}

	overrideSettings(&settings, override)

	m := Model{Steps: steps, ConfigurationMode: configMode, ActiveInput: steps[0].StepInput, Settings: settings}
	return &m, nil
}

// - debug
func NewDebugModel(steps []Step, template string, length int, override Override) *Model {
	settings, err := c.GetUserSettings()
	if err != nil {
		log.Fatal(err)
	}

	overrideSettings(&settings, override)

	m := Model{Steps: steps, ConfigurationMode: false, ActiveInput: steps[len(steps)-1].StepInput, Settings: settings}
	m.Index = len(m.Steps) - 1
	m.Finished = true
	m.DebugMode = true

	switch template {
	default:
		m.Steps[0].Answer.text = ""
		m.Steps[1].Answer.fields = []Field{
			{name: "str", fieldType: "string"},
			{name: "int", fieldType: "number", variant: v.Variant("10000")},
			{name: "bool", fieldType: "boolean"},
			{name: "date", fieldType: "date", variant: v.Timestamp},
			{name: "img", fieldType: "img", variant: v.HorizontalImg},
			{name: "strSet", fieldType: "string enum", variant: v.Variant("a b a_b")},
		}
		m.Steps[3].Answer.text = ""
		m.Steps[4].Answer.text = fmt.Sprint(length)
	case "user":
		m.Steps[0].Answer.text = "users"
		m.Steps[1].Answer.fields = []Field{
			{name: "name", fieldType: "string"},
			{name: "surname", fieldType: "string"},
			{name: "age", fieldType: "number", variant: v.Variant("18 100")},
			{name: "email", fieldType: "string"},
			{name: "premiumAccount", fieldType: "boolean"},
			{name: "role", fieldType: "string enum", variant: v.Variant("user admin mod")},
		}
		m.Steps[3].Answer.text = "User"
		m.Steps[4].Answer.text = fmt.Sprint(length)
	case "dates":
		m.Steps[0].Answer.text = "dates"
		m.Steps[1].Answer.fields = []Field{
			{name: "dateTime", fieldType: "date", variant: v.DateTime},
			{name: "timestamp", fieldType: "date", variant: v.Timestamp},
			{name: "day", fieldType: "date", variant: v.Day},
			{name: "month", fieldType: "date", variant: v.Month},
			{name: "year", fieldType: "date", variant: v.Year},
			{name: "obj", fieldType: "date", variant: v.DateObject},
		}
		m.Steps[3].Answer.text = "Dates"
		m.Steps[4].Answer.text = fmt.Sprint(length)
	case "imgs":
		m.Steps[0].Answer.text = "images"
		m.Steps[1].Answer.fields = []Field{
			{name: "horizontal", fieldType: "img", variant: v.HorizontalImg},
			{name: "profile", fieldType: "img", variant: v.ProfilePictureImg},
			{name: "banner", fieldType: "img", variant: v.Banner},
			{name: "custom", fieldType: "img", variant: v.Variant("5x5")},
		}
		m.Steps[3].Answer.text = "Images"
		m.Steps[4].Answer.text = fmt.Sprint(length)
	}
	return &m
}

func overrideSettings(s *c.Settings, o Override) {
	if o.Language != v.Config {
		s.Language = o.Language
	}

	if o.Output != v.Config {
		s.OutputStyle = o.Output
	}

	if o.Export == v.NoExport || o.Export == v.ExportDefault || o.Export == v.Inline {
		s.Export = o.Export
	}
}

func createTable(rows [][]string, headers []string) *table.Table {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styles.White)).
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
		Headers(headers...).
		Rows(rows...)
	return t
}
