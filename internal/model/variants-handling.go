package model

import (
	"github.com/charmbracelet/bubbles/list"
	v "github.com/liioan/faek/internal/variants"
)

func getVariantList(v []v.VariantData) []list.Item {
	res := []list.Item{}
	for _, i := range v {
		res = append(res, item(i.Value))
	}
	return res
}

func getVariantsValue(variants []v.VariantData, userInput string) v.Variant {
	for _, option := range variants {
		if option.Value == userInput {
			return option.Key
		}
	}
	return v.Variant("")
}
