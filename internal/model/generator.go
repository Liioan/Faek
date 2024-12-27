package model

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	c "github.com/liioan/faek/internal/configuration"
	"github.com/liioan/faek/internal/data"
	e "github.com/liioan/faek/internal/errors"
	"github.com/liioan/faek/internal/utils"
	v "github.com/liioan/faek/internal/variants"
)

var underlyingTypes = map[string]string{
	"strSet": "string",
	"img":    "string",
}

var underlyingDateTypes = map[v.Variant]string{
	"dateTime":  "string",
	"timestamp": "number",
	"day":       "number",
	"month":     "number",
	"year":      "number",
	"date":      "number",
	"object":    "Date",
}

var predefinedValues = map[string][]string{
	"name":    data.Names,
	"surname": data.Surnames,
	"email":   data.Emails,
	"title":   data.Titles,
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
	switch {
	case l == 1:
		res += fmt.Sprintf("%s%s,\n", getIndent(&o.Settings, 1), insertValue(o.Fields[0]))
	case l > 1 && l <= 3:
		res += fmt.Sprintf("%s{", getIndent(&o.Settings, 1))
		res += " "
		for i, field := range o.Fields {
			coma := ","
			if i == l-1 {
				coma = ""
			}
			res += fmt.Sprintf("%s: %s%s ", field.name, insertValue(field), coma)
		}
		res += "},\n"
	case l >= 4:
		res += fmt.Sprintf("%s{", getIndent(&o.Settings, 1))
		res += "\n"
		for _, field := range o.Fields {
			res += fmt.Sprintf("%s%s: %s,\n", getIndent(&o.Settings, 2), field.name, insertValue(field))
		}

		res += fmt.Sprintf("%s},\n", getIndent(&o.Settings, 1))
	}
	return res
}

func insertValue(f Field) string {
	res := ""

	switch f.fieldType {
	case "string":
		if len(predefinedValues[f.name]) > 0 {
			values := predefinedValues[f.name]
			res = values[utils.Random(0, len(values)-1)]
			break
		}
		length := 39 // lorem(39) -> Lorem ipsum, dolor sit amet consectetur
		if len(f.options) > 0 {
			length = utils.ParseInt(f.options[0], length)
		}

		res = fmt.Sprintf("\"%s\"", data.Content[0:length])

	case "number":
		min := 0
		max := 100
		if len(f.options) == 1 {
			max = utils.ParseInt(f.options[0], max)
		} else if len(f.options) >= 2 {
			min = utils.ParseInt(f.options[0], min)
			max = utils.ParseInt(f.options[1], max)
		}
		res = fmt.Sprint(utils.Random(min, max))
	case "boolean":
		if utils.Random(0, 100) >= 50 {
			res = "true"
		} else {
			res = "false"
		}
	case "img":
		dimensions := strings.Split(string(f.variant), "x")
		width := dimensions[0]
		height := dimensions[1]
		res = fmt.Sprintf("\"https://unsplash.it/%s/%s\"", width, height)
	case "date":
		YEAR_IN_DAYS := 365
		YEAR_IN_MONTHS := 12
		MONTH_IN_DAYS := 31
		TEN_YEARS := 10

		multiplier := 1
		if utils.Random(0, 100) >= 50 {
			multiplier = -1
		}

		switch f.variant {
		case v.DateTime:
			res = fmt.Sprintf("\"%s\"", time.Now().AddDate(0, 0, multiplier*rand.Intn(YEAR_IN_DAYS+1)).Format("2.1.2006"))
		case v.Timestamp:
			res = fmt.Sprintf("%d", time.Now().AddDate(0, 0, multiplier*rand.Intn(YEAR_IN_DAYS+1)).Unix())
		case v.Day:
			res = fmt.Sprintf("%d", time.Now().AddDate(0, 0, multiplier*rand.Intn(MONTH_IN_DAYS+1)).Day())
		case v.Month:
			res = fmt.Sprintf("%d", time.Now().AddDate(0, multiplier*rand.Intn(YEAR_IN_MONTHS+1), 0).Month())
		case v.Year:
			res = fmt.Sprintf("%d", time.Now().AddDate(multiplier*rand.Intn(TEN_YEARS+1), 0, 0).Year())
		case v.DateObject:
			res = "new Date()"
		default:
			res = fmt.Sprintf("\"%s\"", time.Now().AddDate(0, 0, multiplier*rand.Intn(YEAR_IN_DAYS+1)).Format("2.1.2006"))
		}
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
				res += fmt.Sprintf("const %s: %s[]", o.AryName, t)
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
	if fieldType == "date" {
		for k, v := range underlyingDateTypes {
			if k == variant {
				return v
			}
		}
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
