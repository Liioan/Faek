package model

import (
	"fmt"
	"strconv"
	"strings"

	c "github.com/liioan/faek/internal/configuration"
	"github.com/liioan/faek/internal/data"
	v "github.com/liioan/faek/internal/variants"
)

var predefinedValues = map[string][]string{
	"name":    data.Names,
	"surname": data.Surnames,
	"email":   data.Emails,
	"title":   data.Titles,
	"city":    data.Cities,
	"street":  data.Streets,
	"country": data.Countries,
}

type OutputMetadata struct {
	AryName    string
	Fields     []Property
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
			res += fmt.Sprintf("\"%s\":%s%s", field.getName(), strings.ReplaceAll(field.generateValue(o.Settings), "`", "\""), separator)
		}
		res += "}"
		res += separator
		return res
	}

	l := len(o.Fields)
	switch {
	case l == 1:
		res += fmt.Sprintf("%s%s,\n", getIndent(&o.Settings, 1), o.Fields[0].generateValue(o.Settings))
	case l > 1 && l <= 3:
		res += fmt.Sprintf("%s{", getIndent(&o.Settings, 1))
		res += " "
		for i, field := range o.Fields {
			coma := ","
			if i == l-1 {
				coma = ""
			}
			res += fmt.Sprintf("%s: %s%s ", field.getName(), field.generateValue(o.Settings), coma)
		}
		res += "},\n"
	case l >= 4:
		res += fmt.Sprintf("%s{", getIndent(&o.Settings, 1))
		res += "\n"
		for _, field := range o.Fields {
			res += fmt.Sprintf("%s%s: %s,\n", getIndent(&o.Settings, 2), field.getName(), field.generateValue(o.Settings))
		}

		res += fmt.Sprintf("%s},\n", getIndent(&o.Settings, 1))
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
			t := o.Fields[0].getUnderlyingType()
			if o.CustomType != "" {
				res += fmt.Sprintf("type %s = %s;\n\nconst %s: %s[]", o.CustomType, t, o.AryName, o.CustomType)
			} else {
				res += fmt.Sprintf("%sconst %s: %s[]", handleExport(o, v.Inline), o.AryName, t)
			}
		case l > 1 && l <= 3:
			if o.CustomType != "" {
				res += fmt.Sprintf("type %s = {\n", o.CustomType)
				for _, field := range o.Fields {
					t := field.getUnderlyingType()
					res += fmt.Sprintf("%s%s: %s\n", getIndent(&o.Settings, 1), field.getName(), t)
				}
				res += fmt.Sprintf("}\n\n%sconst %s: %s[]", handleExport(o, v.Inline), o.AryName, o.CustomType)
			} else {
				res += fmt.Sprintf("%sconst %s: { ", handleExport(o, v.Inline), o.AryName)
				for i, field := range o.Fields {
					coma := ","
					if i == l-1 {
						coma = ""
					}
					t := field.getUnderlyingType()
					res += fmt.Sprintf("%s: %s%s ", field.getName(), t, coma)
				}
				res += "}[]"
			}
		case l >= 4:
			if o.CustomType != "" {
				res += fmt.Sprintf("type %s = {\n", o.CustomType)
				for _, field := range o.Fields {
					t := field.getUnderlyingType()
					res += fmt.Sprintf("%s%s: %s;\n", getIndent(&o.Settings, 1), field.getName(), t)
				}
				res += fmt.Sprintf("}\n\n%sconst %s: %s[]", handleExport(o, v.Inline), o.AryName, o.CustomType)
			} else {
				res += fmt.Sprintf("%sconst %s: {\n", handleExport(o, v.Inline), o.AryName)
				for _, field := range o.Fields {
					t := field.getUnderlyingType()
					res += fmt.Sprintf("%s%s: %s;\n", getIndent(&o.Settings, 1), field.getName(), t)
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

// - debug
func PrintInterview(o *OutputMetadata) string {
	res := ""

	res += "Array name: "
	res += o.AryName + "\n\n"
	res += "Fields: \n"
	for _, f := range o.Fields {
		res += fmt.Sprintf("%s %s %v \n", f.getName(), f.getUnderlyingType(), f.getVariant())
	}
	res += "\n"
	res += "Custom type: "
	res += o.CustomType + "\n\n"
	res += "Length: "
	res += fmt.Sprint(o.Len)

	return res
}
