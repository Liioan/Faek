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
	ArrName    string
	Fields     []Property
	CustomType string
	Len        int

	Settings c.Settings
}

func (m *Model) generateOutput() string {
	res := ""

	outputMetadata := CreateOutputMetadata(m)

	res += handleDeclaration(outputMetadata)

	if outputMetadata.Settings.Language != v.JSON &&
		outputMetadata.Len == 1 {
		res += fmt.Sprintf(" = %s;", outputMetadata.Fields[0].generateValue(outputMetadata.Settings))
		return res
	}

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

func handleObject(outputMetadata *OutputMetadata, iteration int) string {
	res := ""

	if outputMetadata.Settings.Language == v.JSON {
		separator := ","
		if iteration == outputMetadata.Len-1 {
			separator = ""
		}
		res += "{"
		for i, field := range outputMetadata.Fields {
			separator := ","
			if i == len(outputMetadata.Fields)-1 {
				separator = ""
			}
			res += fmt.Sprintf("\"%s\":%s%s", field.getName(), strings.ReplaceAll(field.generateValue(outputMetadata.Settings), "`", "\""), separator)
		}
		res += "}"
		res += separator
		return res
	}

	fieldsLen := len(outputMetadata.Fields)
	switch {
	case fieldsLen == 1:
		res += fmt.Sprintf("%s%s,\n", getIndent(&outputMetadata.Settings, 1), outputMetadata.Fields[0].generateValue(outputMetadata.Settings))
	case fieldsLen > 1 && fieldsLen <= 3:
		res += fmt.Sprintf("%s{", getIndent(&outputMetadata.Settings, 1))
		res += " "
		for i, field := range outputMetadata.Fields {
			coma := ","
			if i == fieldsLen-1 {
				coma = ""
			}
			res += fmt.Sprintf("%s: %s%s ", field.getName(), field.generateValue(outputMetadata.Settings), coma)
		}
		res += "},\n"
	case fieldsLen >= 4:
		res += fmt.Sprintf("%s{", getIndent(&outputMetadata.Settings, 1))
		res += "\n"
		for _, field := range outputMetadata.Fields {
			res += fmt.Sprintf("%s%s: %s,\n", getIndent(&outputMetadata.Settings, 2), field.getName(), field.generateValue(outputMetadata.Settings))
		}

		res += fmt.Sprintf("%s},\n", getIndent(&outputMetadata.Settings, 1))
	}
	return res
}

func handleDeclaration(outputMetadata *OutputMetadata) string {
	lang := outputMetadata.Settings.Language
	res := ""

	switch lang {
	case v.JavaScript:
		res += fmt.Sprintf("const %s", outputMetadata.ArrName)
	case v.TypeScript:
		fieldsLen := len(outputMetadata.Fields)
		if outputMetadata.Len == 1 {
			res += fmt.Sprintf("const %s: %s", outputMetadata.ArrName, outputMetadata.Fields[0].getType())
			return res
		}
		switch {
		case fieldsLen == 1:
			t := outputMetadata.Fields[0].getUnderlyingType()
			if outputMetadata.CustomType != "" {
				res += fmt.Sprintf("type %s = %s;\n\nconst %s: %s[]", outputMetadata.CustomType, t, outputMetadata.ArrName, outputMetadata.CustomType)
			} else {
				res += fmt.Sprintf("%sconst %s: %s[]", handleExport(outputMetadata, v.Inline), outputMetadata.ArrName, t)
			}
		case fieldsLen > 1 && fieldsLen <= 3:
			if outputMetadata.CustomType != "" {
				res += fmt.Sprintf("type %s = {\n", outputMetadata.CustomType)
				for _, field := range outputMetadata.Fields {
					t := field.getUnderlyingType()
					res += fmt.Sprintf("%s%s: %s\n", getIndent(&outputMetadata.Settings, 1), field.getName(), t)
				}
				res += fmt.Sprintf("}\n\n%sconst %s: %s[]", handleExport(outputMetadata, v.Inline), outputMetadata.ArrName, outputMetadata.CustomType)
			} else {
				res += fmt.Sprintf("%sconst %s: { ", handleExport(outputMetadata, v.Inline), outputMetadata.ArrName)
				for i, field := range outputMetadata.Fields {
					coma := ","
					if i == fieldsLen-1 {
						coma = ""
					}
					t := field.getUnderlyingType()
					res += fmt.Sprintf("%s: %s%s ", field.getName(), t, coma)
				}
				res += "}[]"
			}
		case fieldsLen >= 4:
			if outputMetadata.CustomType != "" {
				res += fmt.Sprintf("type %s = {\n", outputMetadata.CustomType)
				for _, field := range outputMetadata.Fields {
					t := field.getUnderlyingType()
					res += fmt.Sprintf("%s%s: %s;\n", getIndent(&outputMetadata.Settings, 1), field.getName(), t)
				}
				res += fmt.Sprintf("}\n\n%sconst %s: %s[]", handleExport(outputMetadata, v.Inline), outputMetadata.ArrName, outputMetadata.CustomType)
			} else {
				res += fmt.Sprintf("%sconst %s: {\n", handleExport(outputMetadata, v.Inline), outputMetadata.ArrName)
				for _, field := range outputMetadata.Fields {
					t := field.getUnderlyingType()
					res += fmt.Sprintf("%s%s: %s;\n", getIndent(&outputMetadata.Settings, 1), field.getName(), t)
				}
				res += "}[]"
			}
		}
	}

	return res
}

func handleExport(outputMetadata *OutputMetadata, selected v.Variant) string {
	if outputMetadata.Settings.Export != selected || outputMetadata.Settings.Language == v.JSON {
		return ""
	}

	res := ""
	switch outputMetadata.Settings.Export {
	case v.Inline:
		res += "export "
	case v.ExportDefault:
		res += fmt.Sprintf("\n\nexport default %s;", outputMetadata.ArrName)
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
	outputMetadata := OutputMetadata{}

	//. get data from user interview
	outputMetadata.ArrName = m.Steps[0].Answer.text
	outputMetadata.Fields = m.Steps[1].Answer.fields
	outputMetadata.CustomType = m.Steps[3].Answer.text
	if outputMetadata.CustomType != "" {
		outputMetadata.CustomType = strings.ToUpper(string(outputMetadata.CustomType[0])) + outputMetadata.CustomType[1:]
	}
	l, err := strconv.Atoi(m.Steps[4].Answer.text)
	if err != nil {
		l = 5
	}
	outputMetadata.Len = l

	if outputMetadata.ArrName == "" {
		outputMetadata.ArrName = "arr"
	}

	outputMetadata.Settings = m.Settings

	return &outputMetadata
}

// - debug
func PrintInterview(outputMetadata *OutputMetadata) string {
	res := ""

	res += "Array name: "
	res += outputMetadata.ArrName + "\n\n"
	res += "Fields: \n"
	for _, f := range outputMetadata.Fields {
		res += fmt.Sprintf("%s %s %v \n", f.getName(), f.getUnderlyingType(), f.getVariant())
	}
	res += "\n"
	res += "Custom type: "
	res += outputMetadata.CustomType + "\n\n"
	res += "Length: "
	res += fmt.Sprint(outputMetadata.Len)

	return res
}
