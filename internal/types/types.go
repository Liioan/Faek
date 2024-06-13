package types

type ValidTypes []string

var ValidTypesArray = ValidTypes{"string", "int", "float", "boolean", "img", "date"}

var TypeConversions = map[string]string{
	"short":  "int",
	"double": "float",
	"str":    "string",
	"char":   "string",
	"bool":   "boolean",
}

var TrueTypes = map[string]string{
	"img":   "string",
	"int":   "number",
	"float": "number",
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
