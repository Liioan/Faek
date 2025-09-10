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
		m.Steps[1].Answer.fields = []Property{
			StringProperty{name: "str", propertyType: "string"},
			NumberProperty{name: "int", propertyType: "number", variant: v.Variant("10000")},
			BooleanProperty{name: "bool", propertyType: "boolean"},
			DateProperty{name: "date", propertyType: "date", variant: v.Timestamp},
			ImageProperty{name: "img", propertyType: "img", variant: v.HorizontalImg},
			EnumProperty{name: "strSet", propertyType: "string enum", variant: v.Variant("a b a_b")},
		}
		m.Steps[3].Answer.text = ""
		m.Steps[4].Answer.text = fmt.Sprint(length)
	case "user":
		m.Steps[0].Answer.text = "users"
		m.Steps[1].Answer.fields = []Property{
			StringProperty{name: "name", propertyType: "string", variant: v.Variant("name")},
			StringProperty{name: "surname", propertyType: "string", variant: v.Variant("surname")},
			NumberProperty{name: "age", propertyType: "number", variant: v.Variant("18 100")},
			StringProperty{name: "email", propertyType: "string", variant: v.Variant("email")},
			BooleanProperty{name: "premiumAccount", propertyType: "boolean"},
			EnumProperty{name: "role", propertyType: "string enum", variant: v.Variant("user admin mod")},
		}
		m.Steps[3].Answer.text = "User"
		m.Steps[4].Answer.text = fmt.Sprint(length)
	case "dates":
		m.Steps[0].Answer.text = "dates"
		m.Steps[1].Answer.fields = []Property{
			DateProperty{name: "dateTime", propertyType: "date", variant: v.DateTime},
			DateProperty{name: "timestamp", propertyType: "date", variant: v.Timestamp},
			DateProperty{name: "day", propertyType: "date", variant: v.Day},
			DateProperty{name: "month", propertyType: "date", variant: v.Month},
			DateProperty{name: "year", propertyType: "date", variant: v.Year},
			DateProperty{name: "obj", propertyType: "date", variant: v.DateObject},
		}
		m.Steps[3].Answer.text = "Dates"
		m.Steps[4].Answer.text = fmt.Sprint(length)
	case "imgs":
		m.Steps[0].Answer.text = "images"
		m.Steps[1].Answer.fields = []Property{
			ImageProperty{name: "horizontal", propertyType: "img", variant: v.HorizontalImg},
			ImageProperty{name: "profile", propertyType: "img", variant: v.ProfilePictureImg},
			ImageProperty{name: "banner", propertyType: "img", variant: v.Banner},
			ImageProperty{name: "custom", propertyType: "img", variant: v.Variant("5x5")},
		}
		m.Steps[3].Answer.text = "Images"
		m.Steps[4].Answer.text = fmt.Sprint(length)
	}
	return &m
}
