package generator

import (
	"faek/internal/constance"
	"faek/internal/data"
	"faek/internal/types"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Generator struct {
	ArrName        string
	CustomType     bool
	CustomTypeName string
	Fields         []types.Field
	Length         int
}

func InitGenerator() Generator {
	return Generator{ArrName: "arr", CustomType: false, Length: 5}
}

func (generator Generator) GenerateOutput(m *Model) string {

	//. array name
	if len(m.Steps[0].Answer) > 0 {
		generator.ArrName = strings.Fields(m.Steps[0].Answer)[0]
	}

	//. fields
	generator.Fields = m.Steps[1].Fields

	//. custom type
	if len(m.Steps[2].Answer) > 0 {
		generator.CustomType = true
		customType := m.Steps[2].Answer
		customType = strings.ToUpper(string(customType[0])) + customType[1:]
		generator.CustomTypeName = strings.Fields(customType)[0]
	}

	//. array length
	length, _ := strconv.Atoi(m.Steps[3].Answer)
	if length > 0 {
		generator.Length = length
	}

	//- generating output
	outputStr := ""

	if len(generator.Fields) == 1 {
		fieldType := generator.Fields[0].FieldType

		field := generator.Fields[0].FieldName
		fieldOptions := generator.Fields[0].Options

		hasFakeType := false
		for fakeType := range types.TrueTypes {
			if fieldType == fakeType {
				hasFakeType = true
			}
		}

		if hasFakeType {
			if generator.CustomType {
				outputStr += fmt.Sprintf("type %s = %s;\n\n", generator.CustomTypeName, types.TrueTypes[fieldType])
			} else {
				generator.CustomTypeName = types.TrueTypes[fieldType]
			}
		} else {
			if generator.CustomType {
				outputStr += fmt.Sprintf("type %s = %s;\n\n", generator.CustomTypeName, fieldType)
			} else {
				generator.CustomTypeName = fieldType
			}

		}

		outputStr += fmt.Sprintf("const %s: %s[] = [\n", generator.ArrName, generator.CustomTypeName)
		for i := 0; i < generator.Length; i++ {
			outputStr += insertData(field, fieldType, len(generator.Fields), fieldOptions)
		}
		outputStr += "];\n"

		return outputStr
	}

	if generator.CustomType {
		//. type declaration
		outputStr += fmt.Sprintf("type %s = {\n", generator.CustomTypeName)
		for _, field := range generator.Fields {
			fieldType := field.FieldType
			hasFakeType := false
			for fakeType := range types.TrueTypes {
				if fieldType == fakeType {
					hasFakeType = true
				}
			}

			if hasFakeType {
				fieldType = types.TrueTypes[fieldType]
			}
			fieldName := field.FieldName
			outputStr += fmt.Sprintf("  %s: %s;\n", fieldName, fieldType)
		}
		outputStr += "};\n\n"

		//. arr declaration
		outputStr += fmt.Sprintf("const %s: %s[] = [\n", generator.ArrName, generator.CustomTypeName)
	} else {
		outputStr += fmt.Sprintf("const %s: { ", generator.ArrName)
		for _, field := range generator.Fields {
			fieldType := field.FieldType
			hasFakeType := false
			for fakeType := range types.TrueTypes {
				if fieldType == fakeType {
					hasFakeType = true
				}
			}
			if hasFakeType {
				fieldType = types.TrueTypes[fieldType]
			}
			fieldName := field.FieldName
			outputStr += fmt.Sprintf("%s: %s; ", fieldName, fieldType)
		}
		outputStr += "}[] = [\n"
	}

	fieldAmount := len(generator.Fields)

	for i := 0; i < generator.Length; i++ {
		outputStr += "  { "
		for _, field := range generator.Fields {
			fieldType := field.FieldType
			fieldName := field.FieldName
			fieldOptions := field.Options
			data := insertData(fieldName, fieldType, fieldAmount, fieldOptions)
			outputStr += data
		}
		if fieldAmount >= constance.LONG_OBJ {
			outputStr += "\n  "
		}
		outputStr += "},\n"
	}

	outputStr += "];\n"

	return outputStr
}

func insertData(field string, fieldType string, fieldAmount int, fieldOptions []string) string {
	//! this is very dirty, but I`m a pepega, and this works
	const itemAmount = 20
	recognizedFields := map[string][]string{
		"name":      data.Names,
		"author":    data.Names,
		"surname":   data.Surnames,
		"lastName":  data.Surnames,
		"last_name": data.Surnames,
		"email":     data.Emails,
		"title":     data.Titles,
		"content":   data.Content,
	}

	imageTypes := map[string]string{
		"default":  data.HorizontalImg,
		"vertical": data.VerticalImg,
		"profile":  data.ProfileImg,
		"article":  data.ArticleImg,
		"banner":   data.BannerImg,
	}

	data := ""

	if fieldAmount >= constance.LONG_OBJ {
		data += "\n    "
	}

	switch fieldType {
	case "string":
		if fieldAmount == 1 {
			if recognizedFields[field] != nil {
				randItem := recognizedFields[field][rand.Intn(itemAmount)]
				data += fmt.Sprintf("  `%s`,\n", randItem)
			} else {
				if len(fieldOptions) > 0 {
					userString := ""
					for i, word := range fieldOptions {
						userString += word
						if i < len(fieldOptions)-1 {
							userString += " "
						}
					}
					data += fmt.Sprintf("  `%s`,\n", userString)
				} else {
					data += fmt.Sprintf("  `%s`,\n", "lorem ipsum dolor sit amet")
				}
			}
		} else {

			if recognizedFields[field] != nil {
				randItem := recognizedFields[field][rand.Intn(itemAmount)]
				data += fmt.Sprintf("%s: `%s`, ", field, randItem)
			} else {
				if len(fieldOptions) > 0 {
					userString := ""
					for i, word := range fieldOptions {
						userString += word
						if i < len(fieldOptions)-1 {
							userString += " "
						}
					}
					data += fmt.Sprintf("%s: `%s`, ", field, userString)
				} else {
					data += fmt.Sprintf("%s: `%s`, ", field, "lorem ipsum dolor sit amet")
				}
			}
		}
	case "number":
		var number int
		switch len(fieldOptions) {
		case 0:
			number = rand.Intn(101)
		case 1:
			MaxNum, err := strconv.Atoi(fieldOptions[0])
			if err != nil {
				MaxNum = 100
			}
			number = rand.Intn(MaxNum + 1)
		case 2:
			LowNum, err := strconv.Atoi(fieldOptions[0])
			if err != nil {
				LowNum = 0
			}
			MaxNum, err := strconv.Atoi(fieldOptions[1])
			if err != nil {
				MaxNum = 100
			}
			number = rand.Intn((MaxNum-LowNum)+1) + LowNum
		}
		if fieldAmount == 1 {
			data += fmt.Sprintf("  %d,\n", number)
		} else {
			data += fmt.Sprintf("%s: %d, ", field, number)
		}
	case "boolean":
		boolean := false
		if rand.Intn(101) >= 50 {
			boolean = true
		}
		if fieldAmount == 1 {
			data += fmt.Sprintf("  %t,\n", boolean)
		} else {
			data += fmt.Sprintf("%s: %t, ", field, boolean)
		}
	case "img":
		switch len(fieldOptions) {
		case 0:
			if fieldAmount == 1 {
				data += fmt.Sprintf("  `%s`,\n", imageTypes["default"])
			} else {
				data += fmt.Sprintf("%s: `%s`,", field, imageTypes["default"])
			}
		case 1:
			for typeName, imgType := range imageTypes {
				if fieldOptions[0] == typeName {
					if fieldAmount == 1 {
						data += fmt.Sprintf(" `%s`,\n", imgType)
					} else {
						data += fmt.Sprintf("%s: `%s`,", field, imgType)
					}
				}
			}
		case 2:
			if fieldAmount == 1 {
				data += fmt.Sprintf("  `unsplash.it/%s/%s`,\n", fieldOptions[0], fieldOptions[1])
			} else {
				data += fmt.Sprintf("%s: `https://unsplash.it/%s/%s`,", field, fieldOptions[0], fieldOptions[1])
			}
		}
	case "strSet":
		randomWord := ""
		switch len(fieldOptions) {
		case 0:
			randomWord = "lorem"
		default:
			randomWord = strings.Replace(fieldOptions[rand.Intn(len(fieldOptions))], "_", " ", -1)
		}
		if fieldAmount == 1 {
			data += fmt.Sprintf("  `%s`,\n", randomWord)
		} else {
			data += fmt.Sprintf("%s: `%s`,", field, randomWord)
		}
	}

	return data
}
