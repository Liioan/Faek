package model

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/liioan/faek/internal/constance"
	o "github.com/liioan/faek/internal/options"
	"github.com/liioan/faek/internal/utils"
)

func newOptionsInput(optionSet o.OptionSet, instruction string) *listInputField {
	optionsArr := []o.OptionData{}

	switch optionSet {
	case o.DateSet:
		optionsArr = o.DateOptions
	case o.ImgSet:
		optionsArr = o.ImgOptions
	case o.OutputSet:
		optionsArr = o.OutputOptions
	case o.LanguageSet:
		utils.LogToDebug(string(optionSet))
		optionsArr = o.LanguageOptions
	}
	l := []list.Item{}
	for _, option := range optionsArr {
		l = append(l, item(option.Value))
	}
	return newListInputField(l, itemDelegate{}, constance.DefaultWidth, constance.ListHeight, instruction)
}

func getOptionsValue(optionSet o.OptionSet, userInput string) o.Option {
	optionsArr := []o.OptionData{}

	switch optionSet {
	case o.DateSet:
		optionsArr = o.DateOptions
	case o.ImgSet:
		optionsArr = o.ImgOptions
	case o.OutputSet:
		optionsArr = o.OutputOptions
	case o.LanguageSet:
		optionsArr = o.LanguageOptions
	}

	for _, option := range optionsArr {
		if option.Value == userInput {
			return option.Key
		}
	}
	return o.Option("")
}
