package model

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	c "github.com/liioan/faek/internal/constance"
	"github.com/liioan/faek/internal/styles"
	v "github.com/liioan/faek/internal/variants"
)

type stepType string

const (
	NormalStep   = "normal"
	EditStep     = "edit"
	OptionalStep = "optional"
	PropStep     = "prop"
)

type activeInput struct {
	input       InputComponent
	instruction string
}

type Field struct {
	name      string
	fieldType string
	variant   v.Variant
}

type Step struct {
	StepInput activeInput
	Answer    struct {
		text   string
		fields []Field
	}
	Variants []v.VariantData

	StepType stepType

	AvailableInputs []activeInput
	InputIdx        int
}

func CreateListStep(title, instruction string, variants []v.VariantData) *Step {
	l := []list.Item{}
	for _, option := range variants {
		l = append(l, item(option.Value))
	}

	fn := func(s ...string) *lipgloss.Style {
		if slices.Contains(s, "no") {
			return &styles.DestructiveItemStyle
		}
		return &styles.SelectedItemStyle
	}

	i := newListInputField(l, itemDelegate{getStyle: fn}, c.DefaultWidth, 12, instruction)

	s := Step{StepInput: activeInput{instruction: title, input: i}, StepType: NormalStep, Variants: variants}
	return &s
}

func CreateTextStep(instruction, placeholder string) *Step {
	i := newTextInputField(placeholder)
	s := Step{StepInput: activeInput{instruction: instruction, input: i}, StepType: NormalStep}
	return &s
}

func CreatePropsStep() *Step {
	types := []list.Item{}
	for _, option := range v.AllTypes {
		types = append(types, item(option))
	}

	nextType := []list.Item{item("yes"), item("no")}

	stringTypes := []list.Item{}
	for _, s := range v.StringTypes {
		stringTypes = append(stringTypes, item(s))
	}

	fn := func(s ...string) *lipgloss.Style {
		if strings.Contains(s[0], "no") {
			return &styles.DestructiveItemStyle
		}
		return &styles.SelectedItemStyle
	}

	dateTypes := getVariantList(v.DateVariants)
	imgTypes := getVariantList(v.ImgVariants)
	idTypes := getVariantList(v.IDVariants)

	i := []activeInput{
		{instruction: "Create your object", input: newListInputField(nextType, itemDelegate{fn}, c.DefaultWidth, 6, "Create another property?")},
		{instruction: "Write property name", input: newTextInputField("e.g. email")},
		{instruction: "Choose type", input: newListInputField(types, itemDelegate{listDefaultStyle}, c.DefaultWidth, c.ListHeight, "available types")},
		{instruction: "Choose type of string data", input: newListInputField(stringTypes, itemDelegate{listDefaultStyle}, c.DefaultWidth, c.ListHeight, "available data")},
		{instruction: "Choose date variant", input: newListInputField(dateTypes, itemDelegate{listDefaultStyle}, c.DefaultWidth, c.ListHeight, "available date variants")},
		{instruction: "Choose img variant", input: newListInputField(imgTypes, itemDelegate{listDefaultStyle}, c.DefaultWidth, c.ListHeight, "available img variants")},
		{instruction: "Write your dimensions: ", input: newTextInputField("e.g. 200x300")},
		{instruction: "Write your range", input: newTextInputField("e.g. 18 60")},
		{instruction: "Write your string set", input: newTextInputField("e.g. user mod admin")},
		{instruction: "Choose id type", input: newListInputField(idTypes, itemDelegate{listDefaultStyle}, c.DefaultWidth, c.ListHeight, "available id types")},
	}
	s := Step{StepType: PropStep, InputIdx: 1, AvailableInputs: i, StepInput: i[1]}
	return &s
}

func CreateEditStep(propStep *Step) *Step {
	fn := func(s ...string) *lipgloss.Style {
		if strings.Contains(s[0], "delete prop") {
			return &styles.DestructiveItemStyle
		}

		return &styles.SelectedItemStyle
	}

	options := []list.Item{item("confirm"), item("add prop"), item("delete prop")}

	i := []activeInput{
		{instruction: "Create your object", input: newListInputField(options, itemDelegate{fn}, c.DefaultWidth, 8, "confirm your object structure")},
		{instruction: "Delete prop"},
	}

	s := Step{StepType: EditStep, AvailableInputs: i, InputIdx: 0, StepInput: i[0]}
	return &s
}

func CreateOptionalStep(listInstruction, listTitle string, options []string, textInstruction, textPlaceholder string) *Step {
	l := []list.Item{}
	for _, o := range options {
		l = append(l, item(o))
	}

	fn := func(s ...string) *lipgloss.Style {
		if strings.Contains(s[0], "no") {
			return &styles.DestructiveItemStyle
		}
		return &styles.SelectedItemStyle
	}

	i := []activeInput{
		{instruction: listInstruction, input: newListInputField(l, itemDelegate{fn}, c.DefaultWidth, c.ListHeight, listTitle)},
		{instruction: textInstruction, input: newTextInputField(textPlaceholder)},
	}

	return &Step{StepType: OptionalStep, InputIdx: 0, AvailableInputs: i, StepInput: i[0]}
}
