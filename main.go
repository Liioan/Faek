package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type ValidTypes []string

var validTypes = ValidTypes{"string", "number", "boolean"}

func (vt ValidTypes) contains(item string) bool {
	for _, v := range vt {
		if v == item {
			return true
		}
	}
	return false
}

const LONG_OBJ = 4

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#44cbca")).MarginLeft(2)
	answerStyle = lipgloss.NewStyle().MarginLeft(2)
	outputStyle = lipgloss.NewStyle().Bold(true)
	quitStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
)

type Output struct {
	arrName        string
	customType     bool
	customTypeName string
	fields         []string
	types          []string
	length         int
}

type Step struct {
	instruction string
	answer      string
	isRepeating bool
	fields      []string
	placeholder string
}

type Model struct {
	index       int
	steps       []Step
	width       int
	height      int
	done        bool
	answerField textinput.Model
}

func main() {
	clearConsole()
	steps := []Step{
		{instruction: "what will the array be called? (default: arr)", placeholder: "E.g. userPosts"},
		{instruction: "write your field (to continue press enter without input)", isRepeating: true, placeholder: "E.g. email string"},
		{instruction: "Create type for your object? (default: no type, input: type name)", placeholder: "E.g. Post"},
		{instruction: "how many items will be in this array (default 5)", placeholder: "E.g. 5"},
	}

	model := New(steps)

	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer file.Close()

	program := tea.NewProgram(*model)
	_, err = program.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Model) generateOutput() string {
	output := Output{arrName: "arr", customType: false, length: 5}
	//. array name
	if len(m.steps[0].answer) > 0 {
		output.arrName = strings.Fields(m.steps[0].answer)[0]
	}
	//. fields
	for _, field := range m.steps[1].fields {
		values := strings.Fields(field)
		if len(values) == 2 {
			if !validTypes.contains(values[1]) {
				continue
			}
			output.fields = append(output.fields, values[0])
			output.types = append(output.types, values[1])
		}
	}
	//. custom type
	if len(m.steps[2].answer) > 0 {
		output.customType = true
		customType := m.steps[2].answer
		customType = strings.ToUpper(string(customType[0])) + customType[1:]
		output.customTypeName = strings.Fields(customType)[0]
	}
	//. array length
	length, _ := strconv.Atoi(m.steps[3].answer)
	if length > 0 {
		output.length = length
	}

	//- generating output
	outputStr := ""

	if len(output.fields) == 1 {
		fieldType := output.types[0]
		field := output.fields[0]

		if output.customType {
			outputStr += fmt.Sprintf("type %s = %s;\n\n", output.customTypeName, fieldType)
		} else {
			output.customTypeName = fieldType
		}

		outputStr += fmt.Sprintf("const %s: %s[] = [\n", output.arrName, output.customTypeName)
		for i := 0; i < output.length; i++ {
			outputStr += insertData(field, fieldType, len(output.fields))
		}
		outputStr += "];\n"

		return outputStr
	}

	if output.customType {
		//. type declaration
		outputStr += fmt.Sprintf("type %s = {\n", output.customTypeName)
		for i, field := range output.fields {
			fieldType := output.types[i]
			outputStr += fmt.Sprintf("  %s: %s;\n", field, fieldType)
		}
		outputStr += "};\n\n"

		//. arr declaration
		outputStr += fmt.Sprintf("const %s: %s[] = [\n", output.arrName, output.customTypeName)
	} else {
		outputStr += fmt.Sprintf("const %s: { ", output.arrName)
		for i, field := range output.fields {
			fieldType := output.types[i]
			outputStr += fmt.Sprintf("%s: %s; ", field, fieldType)
		}
		outputStr += "}[] = [\n"
	}

	fieldAmount := len(output.fields)

	for i := 0; i < output.length; i++ {
		outputStr += "  { "
		for i, field := range output.fields {
			fieldType := output.types[i]
			data := insertData(field, fieldType, fieldAmount)
			outputStr += data
		}
		if fieldAmount >= LONG_OBJ {
			outputStr += "\n  "
		}
		outputStr += "},\n"
	}

	outputStr += "];\n"

	return outputStr
}

func insertData(field string, fieldType string, fieldAmount int) string {
	//! this is very dirty, but I`m a pepega, and this works
	const itemAmount = 20
	recognizedFields := map[string][]string{
		"name":      names,
		"author":    names,
		"surname":   surnames,
		"lastName":  surnames,
		"last_name": surnames,
		"email":     emails,
		"title":     titles,
		"content":   content,
	}

	data := ""

	if fieldAmount == 1 {
		switch fieldType {
		case "string":
			if recognizedFields[field] != nil {
				randItem := recognizedFields[field][rand.Intn(itemAmount)]
				data += fmt.Sprintf("  '%s',\n", randItem)
			} else {
				data += fmt.Sprintf("  '%s',\n", "lorem ipsum dolor sit amet")
			}
		case "number":
			number := rand.Intn(101)
			data += fmt.Sprintf("  %d,\n", number)
		case "boolean":
			boolean := false
			if rand.Intn(101) >= 50 {
				boolean = true
			}
			data += fmt.Sprintf("  %t,\n", boolean)
		}
		return data
	}

	if fieldAmount >= LONG_OBJ {
		data += "\n    "
	}

	switch fieldType {
	case "string":
		if recognizedFields[field] != nil {
			randItem := recognizedFields[field][rand.Intn(itemAmount)]
			data += fmt.Sprintf("%s: '%s', ", field, randItem)
		} else {
			data += fmt.Sprintf("%s: '%s', ", field, "lorem ipsum dolor sit amet")
		}
	case "number":
		number := rand.Intn(101)
		data += fmt.Sprintf("%s: %d, ", field, number)
	case "boolean":
		boolean := false
		if rand.Intn(101) >= 50 {
			boolean = true
		}
		data += fmt.Sprintf("%s: %t, ", field, boolean)
	}

	return data
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Next() {
	if m.index < len(m.steps)-1 {
		m.index++
		m.answerField.Placeholder = m.steps[m.index].placeholder
	} else if m.index == len(m.steps)-1 {
		m.done = true
		m.answerField.Blur()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.steps[m.index]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.index == len(m.steps)-1 {
				return m, tea.Quit
			}
		case "enter":
			checkAnswer(&m, current, m.answerField.Value())
			m.answerField.SetValue("")
			return m, nil
		}
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func checkAnswer(m *Model, current *Step, input string) {
	if current.isRepeating {
		if input == "" {
			if len(current.fields) > 0 {
				m.Next()
				return
			} else {
				return
			}
		} else {
			current.fields = append(current.fields, input)
			return
		}
	}
	current.answer = input
	m.Next()
}

func (m Model) View() string {
	current := m.steps[m.index]
	if m.width == 0 {
		return "loading..."
	}

	if m.done {
		output := m.generateOutput()
		return wordwrap.String(
			fmt.Sprintf(
				"%s\n%s",
				outputStyle.Render(output),
				quitStyle.Render("press q or ctrl+c to exit"),
			),
			m.width,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(current.instruction),
		answerStyle.Render(m.answerField.View()),
	)
}

func New(steps []Step) *Model {
	answerField := textinput.New()
	answerField.Placeholder = steps[0].placeholder
	answerField.Focus()
	return &Model{steps: steps, answerField: answerField}
}

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clearConsole() {
	clear, available := clear[runtime.GOOS]
	if available {
		clear()
	}
}
