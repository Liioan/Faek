package model

import (
	"errors"
	"strconv"

	c "github.com/liioan/faek/internal/configuration"
	e "github.com/liioan/faek/internal/errors"
)

type OutputModel struct {
	AryName    string
	Fields     []Field
	CustomType string
	Len        int

	Settings c.Settings
}

func generateOutput(m *Model) string {
	output := ""

	// outputData, err := NewOutputModel(m)
	// if err != nil {
	// 	return err.Error()
	// }

	test := m.Steps
	for _, step := range test {
		output += "\n"
		if len(step.Answer.fields) > 0 {
			output += "fields: \n"
			for _, field := range step.Answer.fields {
				output += field.name + " " + field.fieldType
				if len(field.variant) > 0 {
					output += " " + string(field.variant)
				}

				if len(field.options) > 0 {
					for _, option := range field.options {
						output += " " + option
					}
				}

				output += "\n"
			}
			output += "\n"
			continue
		}
		output += step.Instruction + ": \n"
		output += step.Answer.text + "\n"
	}

	return output
}

func NewOutputModel(m *Model) (*OutputModel, error) {
	o := OutputModel{}

	//. get data from user interview
	o.AryName = m.Steps[0].Answer.text
	o.Fields = m.Steps[1].Answer.fields
	o.CustomType = m.Steps[2].Answer.text
	l, err := strconv.Atoi(m.Steps[3].Answer.text)
	if err != nil {
		l = 5
	}
	o.Len = l

	if o.AryName == "" {
		o.AryName = "arr"
	}

	//. settings
	settings, err := c.GetUserSettings()
	if err != nil {
		return nil, errors.New(e.SettingsUnavailable)
	}

	o.Settings = settings

	return &o, nil
}
