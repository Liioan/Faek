package model

import (
	"fmt"
	"log"

	c "github.com/liioan/faek/internal/configuration"
	v "github.com/liioan/faek/internal/variants"
)

func NewDebugModel(steps []Step, template string, length int, override Override) *Model {
	settings, err := c.GetUserSettings()
	if err != nil {
		log.Fatal(err)
	}

	overrideSettings(&settings, override)

	m := Model{Steps: steps, ConfigurationMode: false, ActiveInput: steps[len(steps)-1].StepInput, Settings: settings}
	m.Index = len(m.Steps) - 1
	m.Finished = true
	m.DebugMode = true

	switch template {
	default:
		m.Steps[0].Answer.text = ""
		m.Steps[1].Answer.fields = []Field{
			{name: "str", fieldType: "string"},
			{name: "int", fieldType: "number", variant: v.Variant("10000")},
			{name: "bool", fieldType: "boolean"},
			{name: "date", fieldType: "date", variant: v.Timestamp},
			{name: "img", fieldType: "img", variant: v.HorizontalImg},
			{name: "strSet", fieldType: "string enum", variant: v.Variant("a b a_b")},
		}
		m.Steps[3].Answer.text = ""
		m.Steps[4].Answer.text = fmt.Sprint(length)
	case "user":
		m.Steps[0].Answer.text = "users"
		m.Steps[1].Answer.fields = []Field{
			{name: "name", fieldType: "string"},
			{name: "surname", fieldType: "string"},
			{name: "age", fieldType: "number", variant: v.Variant("18 100")},
			{name: "email", fieldType: "string"},
			{name: "premiumAccount", fieldType: "boolean"},
			{name: "role", fieldType: "string enum", variant: v.Variant("user admin mod")},
		}
		m.Steps[3].Answer.text = "User"
		m.Steps[4].Answer.text = fmt.Sprint(length)
	case "dates":
		m.Steps[0].Answer.text = "dates"
		m.Steps[1].Answer.fields = []Field{
			{name: "dateTime", fieldType: "date", variant: v.DateTime},
			{name: "timestamp", fieldType: "date", variant: v.Timestamp},
			{name: "day", fieldType: "date", variant: v.Day},
			{name: "month", fieldType: "date", variant: v.Month},
			{name: "year", fieldType: "date", variant: v.Year},
			{name: "obj", fieldType: "date", variant: v.DateObject},
		}
		m.Steps[3].Answer.text = "Dates"
		m.Steps[4].Answer.text = fmt.Sprint(length)
	case "imgs":
		m.Steps[0].Answer.text = "images"
		m.Steps[1].Answer.fields = []Field{
			{name: "horizontal", fieldType: "img", variant: v.HorizontalImg},
			{name: "profile", fieldType: "img", variant: v.ProfilePictureImg},
			{name: "banner", fieldType: "img", variant: v.Banner},
			{name: "custom", fieldType: "img", variant: v.Variant("5x5")},
		}
		m.Steps[3].Answer.text = "Images"
		m.Steps[4].Answer.text = fmt.Sprint(length)
	}
	return &m
}
