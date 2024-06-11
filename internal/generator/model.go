package generator

import (
	"faek/internal/constance"
	"faek/internal/data"
	"faek/internal/styles"
	"faek/internal/types"
	"faek/internal/utils"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Model struct {
	Index       int
	Steps       []types.Step
	Width       int
	Height      int
	Done        bool
	AnswerField textinput.Model
	Help        bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Next() {
	if m.Index < len(m.Steps)-1 {
		m.Index++
		m.AnswerField.Placeholder = m.Steps[m.Index].Placeholder
	} else if m.Index == len(m.Steps)-1 {
		m.Done = true
		m.AnswerField.Blur()
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
		case "ctrl+h":
			m.Help = !m.Help
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.Index == len(m.Steps)-1 {
				return m, tea.Quit
			}
		case "enter":
			checkAnswer(&m, current, m.AnswerField.Value())
			m.AnswerField.SetValue("")
			return m, nil
		}
	}
	m.AnswerField, cmd = m.AnswerField.Update(msg)
	return m, cmd
}

func checkAnswer(m *Model, current *types.Step, input string) {
	if current.IsRepeating {
		if input == "" {
			if len(current.Fields) > 0 {
				m.Next()
				return
			} else {
				return
			}
		} else {
			values := strings.Fields(input)
			if len(values) == 1 {
				return
			}
			values[1] = strings.ToLower(values[1])
			if len(values) >= 2 {
				for k, v := range types.TypeConversions {
					if k == values[1] {
						values[1] = v
					}
				}
				if types.ValidTypes.Contains(types.ValidTypesArray, values[1]) && !current.ContainsField(values[0]) {
					newField := types.Field{FieldName: values[0], FieldType: values[1]}
					if len(values) >= 3 {
						newField.Options = values[2:]
					}
					current.Fields = append(current.Fields, newField)
				}
			}
			return
		}
	}
	current.Answer = input
	m.Next()
}

func (m Model) View() string {
	current := m.Steps[m.Index]
	if m.Width == 0 {
		return "loading..."
	}

	if m.Help {
		output := fmt.Sprintf("%s\n", styles.HelpHeaderStyle.Render("Help"))
		for _, info := range data.HelpInfo {
			output += fmt.Sprintf("%s\n", info.Style.Render(info.Text))
		}
		return wordwrap.String(output, m.Width)
	}

	if m.Done {
		outputGenerator := InitGenerator()
		output := outputGenerator.GenerateOutput(&m)
		directory, err := os.Getwd()
		if err != nil {
			log.Fatal("cannot get user path")
		}
		if len(strings.Split(output, "\n")) > m.Height {
			if !utils.FileExists(constance.OUTPUT_FILEPATH) {
				file, err := os.Create(constance.OUTPUT_FILEPATH)
				if err != nil {
					log.Fatal("cannot create file")
				}
				defer file.Close()
				file.Write([]byte(output))
			}
			return wordwrap.String(
				fmt.Sprintf(
					"%s\n%s",
					styles.OutputStyle.Render(fmt.Sprintf("Output length exceeded terminal Height, generating output file at %s/faekOutput.ts", directory)),
					styles.QuitStyle.Render("press q or ctrl+c to exit"),
				),
				m.Width,
			)
		} else {
			return wordwrap.String(
				fmt.Sprintf(
					"%s\n%s",
					styles.OutputStyle.Render(output),
					styles.QuitStyle.Render("press q or ctrl+c to exit"),
				),
				m.Width,
			)
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wordwrap.String(styles.TitleStyle.Render(current.Instruction), m.Width),
		styles.AnswerStyle.Render(m.AnswerField.View()),
	)
}

func New(steps []types.Step) *Model {
	AnswerField := textinput.New()
	AnswerField.Placeholder = steps[0].Placeholder
	AnswerField.Focus()
	return &Model{Steps: steps, AnswerField: AnswerField}
}
