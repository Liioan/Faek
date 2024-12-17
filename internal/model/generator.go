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

var underlyingTypes = map[string]string{
	"strSet": "string",
	"img":    "string",
	"date":   "string",
}

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
	res += " = [\n"
	for range outputModel.Len {
		res += handleObject(outputModel)
	}
	res += "];"

	return res
}

func handleObject(o *OutputModel) string {
	res := ""

	l := len(o.Fields)
	res += fmt.Sprintf("%s{", getIndent(&o.Settings, 1))
	switch {
	case l == 1:
	case l <= 3:
		res += " "
		for i, field := range o.Fields {
			coma := ","
			if i == l-1 {
				coma = ""
			}
			res += fmt.Sprintf("%s: test%s ", field.name, coma)
		}
		res += "},\n"
	case l >= 4:
		res += "\n"
		for _, field := range o.Fields {
			res += fmt.Sprintf("%s%s: test,\n", getIndent(&o.Settings, 2), field.name)
		}

		res += fmt.Sprintf("%s},\n", getIndent(&o.Settings, 1))
	}
	return res
}

func handleType(o *OutputModel) string {
	lang := o.Settings.Language
	res := ""

	switch v.Variant(lang) {
	case v.JavaScript:
		res += fmt.Sprintf("const %s", o.AryName)
	case v.TypeScript:
		l := len(o.Fields)
		switch {
		case l == 1:
			t := getUnderlyingType(o.Fields[0].fieldType, o.Fields[0].variant)
			if o.CustomType != "" {
				res += fmt.Sprintf("type %s = %s;\n\nconst %s: %s[]", o.CustomType, t, o.AryName, o.CustomType)
			} else {
				res += fmt.Sprintf("const %s: %s[]\n", o.AryName, t)
			}
		case l > 1 && l <= 3:
			if o.CustomType != "" {
				res += fmt.Sprintf("type %s = {\n", o.CustomType)
				for _, field := range o.Fields {
					t := getUnderlyingType(field.fieldType, field.variant)
					res += fmt.Sprintf("%s%s: %s\n", getIndent(&o.Settings, 1), field.name, t)
				}
				res += fmt.Sprintf("}\n\nconst %s: %s[]", o.AryName, o.CustomType)
			} else {
				res += fmt.Sprintf("const %s: { ", o.AryName)
				for i, field := range o.Fields {
					coma := ","
					if i == l-1 {
						coma = ""
					}
					t := getUnderlyingType(field.fieldType, field.variant)
					res += fmt.Sprintf("%s: %s%s ", field.name, t, coma)
				}
				res += "}[]"
			}
		case l >= 4:
			if o.CustomType != "" {
				res += fmt.Sprintf("type %s = {\n", o.CustomType)
				for _, field := range o.Fields {
					t := getUnderlyingType(field.fieldType, field.variant)
					res += fmt.Sprintf("%s%s: %s\n", getIndent(&o.Settings, 1), field.name, t)
				}
				res += fmt.Sprintf("}\n\nconst %s: %s[]", o.AryName, o.CustomType)
			} else {
				res += fmt.Sprintf("const %s: {\n", o.AryName)
				for _, field := range o.Fields {
					t := getUnderlyingType(field.fieldType, field.variant)
					res += fmt.Sprintf("%s%s: %s;\n", getIndent(&o.Settings, 1), field.name, t)
				}
				res += "}"
			}
		}
	}

	return res
}

func getUnderlyingType(fieldType string, variant v.Variant) string {
	if fieldType == "date" && variant == v.DateObject {
		return "Date"
	}

	for k, v := range underlyingTypes {
		if k == fieldType {
			return v
		}
	}

	return fieldType
}

func getIndent(s *c.Settings, level int) string {
	indent, err := strconv.Atoi(s.Indent)
	if err != nil {
		indent = 2
	}
	length := indent * level
	str := ""
	for {
		str = " " + str
		if len(str) > length {
			return str[0:length]
		}
	}
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
