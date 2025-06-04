package model

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	c "github.com/liioan/faek/internal/configuration"
	"github.com/liioan/faek/internal/data"
	"github.com/liioan/faek/internal/utils"
	v "github.com/liioan/faek/internal/variants"
)

var underlyingTypes = map[string]string{
	"strSet": "string",
	"img":    "string",
	"id":     "string",
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

type OutputMetadata struct {
	AryName    string
	Fields     []Field
	CustomType string
	Len        int

	Settings c.Settings
}

func (m *Model) generateOutput() string {
	res := ""

	outputMetadata := CreateOutputMetadata(m)

	res += handleDeclaration(outputMetadata)

	if outputMetadata.Settings.Language != v.JSON {
		res += " = [\n"
	} else {
		res += "["
	}

	for i := range outputMetadata.Len {
		res += handleObject(outputMetadata, i)
	}
	if outputMetadata.Settings.Language != v.JSON {
		res += "];"
	} else {
		res += "]"
	}

	res += handleExport(outputMetadata, v.ExportDefault)

	return res
}

func handleObject(o *OutputMetadata, iteration int) string {
	res := ""

	if o.Settings.Language == v.JSON {
		separator := ","
		if iteration == o.Len-1 {
			separator = ""
		}
		res += "{"
		for i, field := range o.Fields {
			separator := ","
			if i == len(o.Fields)-1 {
				separator = ""
			}
			res += fmt.Sprintf("\"%s\":%s%s", field.name, strings.ReplaceAll(insertValue(field), "`", "\""), separator)
		}
		res += "}"
		res += separator
		return res
	}

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
		if len(predefinedValues[string(f.variant)]) > 0 {
			values := predefinedValues[string(f.variant)]
			res = fmt.Sprintf("`%s`", values[utils.Random(0, len(values)-1)])
			break
		}
		length := 39 // lorem(39) -> Lorem ipsum, dolor sit amet consectetur

		if f.variant == v.Variant("content") {
			length = len(data.Content) - 1
		}

		text := data.Content
		res = fmt.Sprintf("`%s`", text[0:length])

	case "number":
		min := 0
		max := 100
		numRange := strings.Split(string(f.variant), " ")

		if len(numRange) == 1 {
			max = utils.ParseInt(numRange[0], max)
		} else if len(numRange) >= 2 {
			min = utils.ParseInt(numRange[0], min)
			max = utils.ParseInt(numRange[1], max)
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
		res = fmt.Sprintf("`https://unsplash.it/%s/%s`", width, height)
	case "date":
		YEAR_IN_DAYS := 365
		YEAR_IN_MONTHS := 12
		MONTH_IN_DAYS := 31
		TEN_YEARS := 10
		switch f.variant {
		case v.DateTime:
			res = fmt.Sprintf("`%s`", time.Now().AddDate(0, 0, -1*rand.Intn(YEAR_IN_DAYS+1)).Format("2.1.2006"))
		case v.Timestamp:
			res = fmt.Sprintf("%d", time.Now().AddDate(0, 0, -1*rand.Intn(YEAR_IN_DAYS+1)).Unix()*1000) // unix time to js timestamp
		case v.Day:
			res = fmt.Sprintf("%d", time.Now().AddDate(0, 0, -1*rand.Intn(MONTH_IN_DAYS+1)).Day())
		case v.Month:
			res = fmt.Sprintf("%d", time.Now().AddDate(0, -1*rand.Intn(YEAR_IN_MONTHS+1), 0).Month())
		case v.Year:
			res = fmt.Sprintf("%d", time.Now().AddDate(-1*rand.Intn(TEN_YEARS+1), 0, 0).Year())
		case v.DateObject:
			res = "new Date()"
		default:
			res = fmt.Sprintf("`%s`", time.Now().AddDate(0, 0, -1*rand.Intn(YEAR_IN_DAYS+1)).Format("2.1.2006"))
		}
	case "string enum":
		wordSet := strings.Split(data.Content[0:39], " ")
		v := strings.Split(string(f.variant), " ")
		if len(v) != 0 {
			wordSet = v
		}
		wordSet = parseStringEnum(wordSet)
		randStr := wordSet[utils.Random(0, len(wordSet)-1)]
		res = fmt.Sprintf("`%s`", randStr)
	case "id":
		res = "`BLANK_ID`"
	case "null":
		res = "null"
	case "undefined":
		res = "undefined"
	}

	return res
}

func handleDeclaration(o *OutputMetadata) string {
	lang := o.Settings.Language
	res := ""

	switch lang {
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
				res += fmt.Sprintf("%sconst %s: %s[]", handleExport(o, v.Inline), o.AryName, t)
			}
		case l > 1 && l <= 3:
			if o.CustomType != "" {
				res += fmt.Sprintf("type %s = {\n", o.CustomType)
				for _, field := range o.Fields {
					t := getUnderlyingType(field.fieldType, field.variant)
					res += fmt.Sprintf("%s%s: %s\n", getIndent(&o.Settings, 1), field.name, t)
				}
				res += fmt.Sprintf("}\n\n%sconst %s: %s[]", handleExport(o, v.Inline), o.AryName, o.CustomType)
			} else {
				res += fmt.Sprintf("%sconst %s: { ", handleExport(o, v.Inline), o.AryName)
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
					res += fmt.Sprintf("%s%s: %s;\n", getIndent(&o.Settings, 1), field.name, t)
				}
				res += fmt.Sprintf("}\n\n%sconst %s: %s[]", handleExport(o, v.Inline), o.AryName, o.CustomType)
			} else {
				res += fmt.Sprintf("%sconst %s: {\n", handleExport(o, v.Inline), o.AryName)
				for _, field := range o.Fields {
					t := getUnderlyingType(field.fieldType, field.variant)
					res += fmt.Sprintf("%s%s: %s;\n", getIndent(&o.Settings, 1), field.name, t)
				}
				res += "}[]"
			}
		}
	}

	return res
}

func handleExport(o *OutputMetadata, selected v.Variant) string {
	if o.Settings.Export != selected || o.Settings.Language == v.JSON {
		return ""
	}

	res := ""
	switch o.Settings.Export {
	case v.Inline:
		res += "export "
	case v.ExportDefault:
		res += fmt.Sprintf("\n\nexport default %s;", o.AryName)
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

	if fieldType == "string enum" {
		//- () for string enum in type def
		res := "("
		wordSet := parseStringEnum(strings.Split(string(variant), " "))

		for i, w := range wordSet {
			separator := " | "
			if i == len(wordSet)-1 {
				separator = ""
			}
			res += fmt.Sprintf(`"%s"%s`, w, separator)
		}
		res += ")"
		return res
	}

	for k, v := range underlyingTypes {
		if k == fieldType {
			return v
		}
	}

	return fieldType
}

func parseStringEnum(wordSet []string) []string {
	res := []string{}
	for _, s := range wordSet {
		res = append(res, strings.ReplaceAll(s, "_", " "))
	}

	return res
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

func CreateOutputMetadata(m *Model) *OutputMetadata {
	o := OutputMetadata{}

	//. get data from user interview
	o.AryName = m.Steps[0].Answer.text
	o.Fields = m.Steps[1].Answer.fields
	o.CustomType = m.Steps[3].Answer.text
	if o.CustomType != "" {
		o.CustomType = strings.ToUpper(string(o.CustomType[0])) + o.CustomType[1:]
	}
	l, err := strconv.Atoi(m.Steps[4].Answer.text)
	if err != nil {
		l = 5
	}
	o.Len = l

	if o.AryName == "" {
		o.AryName = "arr"
	}

	o.Settings = m.Settings

	return &o
}

func PrintInterview(o *OutputMetadata) string {
	res := ""

	res += "Array name: "
	res += o.AryName + "\n\n"
	res += "Fields: \n"
	for _, f := range o.Fields {
		res += fmt.Sprintf("%s %s %v \n", f.name, f.fieldType, f.variant)
	}
	res += "\n"
	res += "Custom type: "
	res += o.CustomType + "\n\n"
	res += "Length: "
	res += fmt.Sprint(o.Len)

	return res
}
