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

var validTypes = ValidTypes{"string", "number", "boolean", "img"}

var typeConversions = map[string]string{
	"int":   "number",
	"float": "number",
	"short": "number",
	"str":   "string",
	"char":  "string",
	"bool":  "boolean",
}

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
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#44cbca")).MarginLeft(2)
	answerStyle     = lipgloss.NewStyle().MarginLeft(2)
	outputStyle     = lipgloss.NewStyle().Bold(true)
	quitStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	helpHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#44cbca")).Bold(true)
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))
)

type Field struct {
	fieldName string
	fieldType string
	options   []string
}

type Output struct {
	arrName        string
	customType     bool
	customTypeName string
	fields         []Field
	length         int
}

type Step struct {
	instruction string
	answer      string
	isRepeating bool
	fields      []Field
	placeholder string
}

func (s Step) containsField(name string) bool {
	for _, field := range s.fields {
		if field.fieldName == name {
			return true
		}
	}
	return false
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
		{instruction: "press enter to continue, ctrl+c to quit or ctrl+h for help", placeholder: ""},
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
	if len(m.steps[1].answer) > 0 {
		output.arrName = strings.Fields(m.steps[1].answer)[0]
	}

	//. fields
	output.fields = m.steps[2].fields

	//. custom type
	if len(m.steps[3].answer) > 0 {
		output.customType = true
		customType := m.steps[3].answer
		customType = strings.ToUpper(string(customType[0])) + customType[1:]
		output.customTypeName = strings.Fields(customType)[0]
	}

	//. array length
	length, _ := strconv.Atoi(m.steps[4].answer)
	if length > 0 {
		output.length = length
	}

	//- generating output
	outputStr := ""

	if len(output.fields) == 1 {
		fieldType := output.fields[0].fieldType

		field := output.fields[0].fieldName
		fieldOptions := output.fields[0].options

		if fieldType != "img" {
			if output.customType {
				outputStr += fmt.Sprintf("type %s = %s;\n\n", output.customTypeName, fieldType)
			} else {
				output.customTypeName = fieldType
			}
		} else {
			if output.customType {
				outputStr += fmt.Sprintf("type %s = %s;\n\n", output.customTypeName, "string")
			} else {
				output.customTypeName = "string"
			}
		}

		outputStr += fmt.Sprintf("const %s: %s[] = [\n", output.arrName, output.customTypeName)
		for i := 0; i < output.length; i++ {
			outputStr += insertData(field, fieldType, len(output.fields), fieldOptions)
		}
		outputStr += "];\n"

		return outputStr
	}

	if output.customType {
		//. type declaration
		outputStr += fmt.Sprintf("type %s = {\n", output.customTypeName)
		for _, field := range output.fields {
			fieldType := field.fieldType
			if fieldType == "img" {
				fieldType = "string"
			}
			fieldName := field.fieldName
			outputStr += fmt.Sprintf("  %s: %s;\n", fieldName, fieldType)
		}
		outputStr += "};\n\n"

		//. arr declaration
		outputStr += fmt.Sprintf("const %s: %s[] = [\n", output.arrName, output.customTypeName)
	} else {
		outputStr += fmt.Sprintf("const %s: { ", output.arrName)
		for _, field := range output.fields {
			fieldType := field.fieldType
			if fieldType == "img" {
				fieldType = "string"
			}
			fieldName := field.fieldName
			outputStr += fmt.Sprintf("%s: %s; ", fieldName, fieldType)
		}
		outputStr += "}[] = [\n"
	}

	fieldAmount := len(output.fields)

	for i := 0; i < output.length; i++ {
		outputStr += "  { "
		for _, field := range output.fields {
			fieldType := field.fieldType
			fieldName := field.fieldName
			fieldOptions := field.options
			data := insertData(fieldName, fieldType, fieldAmount, fieldOptions)
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

func insertData(field string, fieldType string, fieldAmount int, fieldOptions []string) string {
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

	imageTypes := map[string]string{
		"profile": profileImg,
		"article": articleImg,
		"banner":  bannerImg,
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
			var number int
			switch len(fieldOptions) {
			case 0:
				number = rand.Intn(101)
			case 1:
				MaxNum, err := strconv.Atoi(fieldOptions[0])
				if err != nil {
					MaxNum = 100
				}
				number = rand.Intn(MaxNum + 1)
			case 2:
				LowNum, err := strconv.Atoi(fieldOptions[0])
				if err != nil {
					LowNum = 0
				}
				MaxNum, err := strconv.Atoi(fieldOptions[1])
				if err != nil {
					MaxNum = 100
				}
				number = rand.Intn((MaxNum-LowNum)+1) + LowNum
			}
			data += fmt.Sprintf("  %d,\n", number)
		case "boolean":
			boolean := false
			if rand.Intn(101) >= 50 {
				boolean = true
			}
			data += fmt.Sprintf("  %t,\n", boolean)
		case "img":
			switch len(fieldOptions) {
			case 0:
				data += fmt.Sprintf("  '%s',\n", img)
			case 1:
				for typeName, imgType := range imageTypes {
					if fieldOptions[0] == typeName {
						data += fmt.Sprintf(" '%s',\n", imgType)
					}
				}
			case 3:
				data += fmt.Sprintf(" 'unsplash.it/%s/%s'\n", fieldOptions[1], fieldOptions[2])
			}
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
		var number int
		switch len(fieldOptions) {
		case 0:
			number = rand.Intn(101)
		case 1:
			MaxNum, err := strconv.Atoi(fieldOptions[0])
			if err != nil {
				MaxNum = 100
			}
			number = rand.Intn(MaxNum + 1)
		case 2:
			LowNum, err := strconv.Atoi(fieldOptions[0])
			if err != nil {
				LowNum = 0
			}
			MaxNum, err := strconv.Atoi(fieldOptions[1])
			if err != nil {
				MaxNum = 100
			}
			number = rand.Intn((MaxNum-LowNum)+1) + LowNum
		}
		data += fmt.Sprintf("%s: %d, ", field, number)
	case "boolean":
		boolean := false
		if rand.Intn(101) >= 50 {
			boolean = true
		}
		data += fmt.Sprintf("%s: %t, ", field, boolean)
	case "img":
		switch len(fieldOptions) {
		case 0:
			data += fmt.Sprintf("%s: '%s',", field, img)
		case 1:
			for typeName, imgType := range imageTypes {
				if fieldOptions[0] == typeName {
					data += fmt.Sprintf("%s: '%s',", field, imgType)
				}
			}
		case 3:
			data += fmt.Sprintf("%s: 'https://unsplash.it/%s/%s',", field, fieldOptions[1], fieldOptions[2])
		}

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
		case "ctrl+h":
			if m.steps[0].answer == "help" {
				m.steps[0].answer = ""
			} else {
				m.steps[0].answer = "help"
			}
			return m, nil
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
			values := strings.Fields(input)
			if len(values) == 1 {
				return
			}
			values[1] = strings.ToLower(values[1])
			if len(values) >= 2 {
				for k, v := range typeConversions {
					if k == values[1] {
						values[1] = v
					}
				}
				if validTypes.contains(values[1]) && !current.containsField(values[0]) {
					newField := Field{fieldName: values[0], fieldType: values[1]}
					if len(values) >= 3 {
						newField.options = values[2:]
					}
					current.fields = append(current.fields, newField)
				}
			}
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

	if m.steps[0].answer == "help" {
		output := fmt.Sprintf("%s\n", helpHeaderStyle.Render("Help"))
		for _, info := range helpInfo {
			output += fmt.Sprintf("%s\n", info.style.Render(info.text))
		}
		return wordwrap.String(output, m.width)
	}

	if m.done {
		output := m.generateOutput()
		if len(strings.Split(output, "\n")) > m.height {
			filePath := "output.ts"
			file, err := os.Create(filePath)
			if err != nil {
				log.Fatal("cannot create file")
			}
			defer file.Close()
			file.Write([]byte(output))
			return wordwrap.String(
				fmt.Sprintf(
					"%s\n%s",
					outputStyle.Render("Output length exceeded terminal height, generating output file"),
					quitStyle.Render("press q or ctrl+c to exit"),
				),
				m.width,
			)
		} else {
			return wordwrap.String(
				fmt.Sprintf(
					"%s\n%s",
					outputStyle.Render(output),
					quitStyle.Render("press q or ctrl+c to exit"),
				),
				m.width,
			)
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wordwrap.String(titleStyle.Render(current.instruction), m.width),
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
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
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
