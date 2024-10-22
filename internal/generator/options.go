package generator

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/liioan/faek/internal/constance"
)

type Option string

type OptionData struct {
	key   string
	value Option
}

const (
	HorizontalImg     Option = "300x500"
	VerticalImg       Option = "500x300"
	ProfilePictureImg Option = "100x100"
	ArticleImg        Option = "600x400"
	Banner            Option = "600x240"
	Custom            Option = "custom"
)

var imgOptions = []OptionData{
	{"Horizontal (default) 300x500", HorizontalImg},
	{"Vertical 500x300", VerticalImg},
	{"Profile Picture 100x100", ProfilePictureImg},
	{"Article 600x400", ArticleImg},
	{"Banner 600x240", Banner},
	{"Custom", Custom},
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
	{"dateTime: e.g. 27.02.2024", DateTime},
	{"timestamp: e.g. 1718051654", Timestamp},
	{"day: 0-31", Day},
	{"month: 0-12", Month},
	{"year: current year", Year},
	{"object: new Date()", DateObject},
}

func newOptionsInput(fieldType string, instruction string) *listInputField {
	options := []OptionData{}
	switch fieldType {
	case "date":
		options = dateOptions
	case "img":
		options = imgOptions
	}
	l := []list.Item{}
	for _, option := range options {
		l = append(l, item(option.key))
	}
	return newListInputField(l, itemDelegate{}, constance.DefaultWidth, constance.ListHeight, instruction)
}

func getOptionsValue(fieldType string, userInput string) Option {
	switch fieldType {
	case "date":
		for _, dateOption := range dateOptions {
			if dateOption.key == userInput {
				return dateOption.value
			}
		}
	case "img":
		for _, imgOption := range imgOptions {
			if imgOption.key == userInput {
				return imgOption.value
			}
		}
	}
	return Option("")
}
