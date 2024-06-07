package types

type ValidTypes []string

var ValidTypesArray = ValidTypes{"string", "number", "boolean", "img", "strSet"}

var TypeConversions = map[string]string{
	"int":       "number",
	"float":     "number",
	"short":     "number",
	"str":       "string",
	"char":      "string",
	"bool":      "boolean",
	"stringSet": "strSet",
	"ss":        "strSet",
	"strs":      "strSet",
	"strset":    "strSet",
}

var TrueTypes = map[string]string{
	"img":    "string",
	"strSet": "string",
}

func (vt ValidTypes) Contains(item string) bool {
	for _, v := range vt {
		if v == item {
			return true
		}
	}
	return false
}

type Field struct {
	FieldName string
	FieldType string
	Options   []string
}

type Step struct {
	Instruction string
	Answer      string
	IsRepeating bool
	Fields      []Field
	Placeholder string
}

func (s Step) ContainsField(name string) bool {
	for _, field := range s.Fields {
		if field.FieldName == name {
			return true
		}
	}
	return false
}
