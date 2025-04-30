package model

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/liioan/faek/internal/constance"
	v "github.com/liioan/faek/internal/variants"
)

func newVariantsInput(variants []v.VariantData, instruction string) *listInputField {

	l := []list.Item{}
	for _, option := range variants {
		l = append(l, item(option.Value))
	}
	return newListInputField(l, itemDelegate{}, constance.DefaultWidth, constance.ListHeight, instruction)
}

func getVariantsValue(variants []v.VariantData, userInput string) v.Variant {
	for _, option := range variants {
		if option.Value == userInput {
			return option.Key
		}
	}
	return v.Variant("")
}
