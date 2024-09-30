package generator

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liioan/faek/internal/constance"
)

type inputMode string

const (
	Text inputMode = "text"
	List inputMode = "list"
)

type activeInput struct {
	input InputComponent
	mode  inputMode
}

type Step struct {
	Instruction string
	StepInput   activeInput
	Answer      struct {
		text   string
		fields []struct {
			name      string
			fieldType string
			option    string
		}
	}
	Repeat bool
}

func newListStep(instruction string, repeats bool, options []list.Item) *Step {
	i := activeInput{
		input: newListInputField(options, itemDelegate{}, constance.DefaultWidth, constance.ListHeight, instruction),
		mode:  List,
	}

	s := Step{Instruction: instruction, Repeat: repeats, StepInput: i}
	return &s
}

func newTextStep(instruction, placeholder string, repeats bool) *Step {
	i := activeInput{
		input: newTextInputField(placeholder),
		mode:  Text,
	}

	s := Step{Instruction: instruction, Repeat: repeats, StepInput: i}
	return &s
}

type Model struct {
	Index       int
	Width       int
	Height      int
	Finished    bool
	ActiveInput activeInput
	Steps       []Step
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Next() {
	if m.Index < len(m.Steps)-1 {
		m.Index++
	} else if m.Index == len(m.Steps)-1 {
		m.Finished = true
	}
	m.ActiveInput.input = m.Steps[m.Index].StepInput.input
}

func (m Model) View() {

}
