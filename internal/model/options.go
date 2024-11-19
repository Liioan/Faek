package model

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/liioan/faek/internal/constance"
)

type Option string

type OptionData struct {
	key   Option
	value string
}

// - generator options
const (
	HorizontalImg     Option = "300x500"
	VerticalImg       Option = "500x300"
	ProfilePictureImg Option = "100x100"
	ArticleImg        Option = "600x400"
	Banner            Option = "600x240"
	Custom            Option = "custom"
)

var imgOptions = []OptionData{
	{HorizontalImg, "Horizontal (default) 300x500"},
	{VerticalImg, "Vertical 500x300"},
	{ProfilePictureImg, "Profile Picture 100x100"},
	{ArticleImg, "Article 600x400"},
	{Banner, "Banner 600x240"},
	{Custom, "Custom"},
}

const (
	DateTime   Option = "dateTime"
	Timestamp  Option = "timestamp"
	Day        Option = "day"
	Month      Option = "month"
	Year       Option = "year"
	DateObject Option = "object"
)

var dateOptions = []OptionData{
	{DateTime, "dateTime: e.g. 27.02.2024"},
	{Timestamp, "timestamp: e.g. 1718051654"},
	{Day, "day: 0-31"},
	{Month, "month: 0-12"},
	{Year, "year: current year"},
	{DateObject, "object: new Date()"},
}

//- configuration options

const (
	Terminal Option = "terminal"
	File     Option = "file"
)

var outputOptions = []OptionData{
	{Terminal, "In terminal"},
	{File, "Output file"},
}

func newOptionsInput(fieldType string, instruction string) *listInputField {
	options := []OptionData{}
	switch fieldType {
	case "date":
		options = dateOptions
	case "img":
		options = imgOptions
	case "output":
		options = outputOptions
	}
	l := []list.Item{}
	for _, option := range options {
		l = append(l, item(option.value))
	}
	return newListInputField(l, itemDelegate{}, constance.DefaultWidth, constance.ListHeight, instruction)
}

func getOptionsValue(fieldType string, userInput string) Option {
	optionsArr := []OptionData{}

	switch fieldType {
	case "date":
		optionsArr = dateOptions
	case "img":
		optionsArr = imgOptions
	case "output":
		optionsArr = outputOptions
	}

	for _, option := range optionsArr {
		if option.value == userInput {
			return option.key
		}
	}
	return Option("")
}
