package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	c "github.com/liioan/faek/internal/configuration"
	e "github.com/liioan/faek/internal/errors"
	"github.com/liioan/faek/internal/utils"
	v "github.com/liioan/faek/internal/variants"
)

type OutputModel struct {
	AryName    string
	Fields     []Field
	CustomType string
	Len        int

	Settings c.Settings
}

func generateOutput(m *Model) string {
	res := ""

	outputModel, err := NewOutputModel(m)
	if err != nil {
		return err.Error()
	}
	utils.LogToDebug(PrintInterview(outputModel))

	res += handleType(outputModel)

	res += "];"
	return res
}

func handleType(o *OutputModel) string {
	// https://excalidraw.com/#json=W7TYYmjH3zP67GIeOQnTl,QFsPbGhJNR9xBtnLQKgBPg
	lang := o.Settings.Language
	res := ""

	switch v.Variant(lang) {
	case v.JavaScript:
		res += fmt.Sprintf("const %s = [", o.AryName)
	case v.TypeScript:
		switch len(o.Fields) {
		case 1:
			t := o.Fields[0].fieldType
			if o.CustomType != "" {
				res += fmt.Sprintf("type %s = %s;\n\nconst %s: %s[] = [\n", o.CustomType, t, o.AryName, o.CustomType)
			} else {
				res += fmt.Sprintf("const %s: %s[] = [\n", o.AryName, t)
			}
		}
	}

	return res
}

func NewOutputModel(m *Model) (*OutputModel, error) {
	o := OutputModel{}

	//. get data from user interview
	o.AryName = m.Steps[0].Answer.text
	o.Fields = m.Steps[1].Answer.fields
	o.CustomType = m.Steps[2].Answer.text
	if o.CustomType != "" {
		o.CustomType = strings.ToUpper(string(o.CustomType[0])) + o.CustomType[1:]
	}
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

func PrintInterview(o *OutputModel) string {
	res := ""

	res += "Array name: "
	res += o.AryName + "\n\n"
	res += "Fields: \n"
	for _, f := range o.Fields {
		res += fmt.Sprintf("%s %s %v %v \n", f.name, f.fieldType, f.options, f.variant)
	}
	res += "\n"
	res += "Custom type: "
	res += o.CustomType + "\n\n"
	res += "Length: "
	res += fmt.Sprint(o.Len)

	return res
}
