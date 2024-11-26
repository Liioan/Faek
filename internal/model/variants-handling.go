package model

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/liioan/faek/internal/constance"
	v "github.com/liioan/faek/internal/variants"
)

func newVariantsInput(optionSet v.VariantSet, instruction string) *listInputField {
	variants := []v.VariantData{}

	switch optionSet {
	case v.DateSet:
		variants = v.DateVariants
	case v.ImgSet:
		variants = v.ImgVariants
	case v.OutputSet:
		variants = v.OutputVariants
	case v.LanguageSet:
		variants = v.LanguageVariants
	}
	l := []list.Item{}
	for _, option := range variants {
		l = append(l, item(option.Value))
	}
	return newListInputField(l, itemDelegate{}, constance.DefaultWidth, constance.ListHeight, instruction)
}

func getVariantsValue(optionSet v.VariantSet, userInput string) v.Variant {
	variants := []v.VariantData{}

	switch optionSet {
	case v.DateSet:
		variants = v.DateVariants
	case v.ImgSet:
		variants = v.ImgVariants
	case v.OutputSet:
		variants = v.OutputVariants
	case v.LanguageSet:
		variants = v.LanguageVariants
	}

	for _, option := range variants {
		if option.Value == userInput {
			return option.Key
		}
	}
	return v.Variant("")
}
