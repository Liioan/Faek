package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	c "github.com/liioan/faek/internal/configuration"
	"github.com/liioan/faek/internal/data"
	"github.com/liioan/faek/internal/utils"
	v "github.com/liioan/faek/internal/variants"
)

type Property interface {
	getName() string
	getType() string
	getUnderlyingType() string
	generateValue(c.Settings) string
	getVariant() v.Variant
	setVariant(v.Variant)
	clone(...string) Property
}

type PlaceholderProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p PlaceholderProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p PlaceholderProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p PlaceholderProperty) getName() string {
	return p.name
}

func (p PlaceholderProperty) getUnderlyingType() string {
	return ""
}

func (p PlaceholderProperty) getType() string {
	return p.propertyType
}

func (p PlaceholderProperty) getVariant() v.Variant {
	return p.variant
}

func (p PlaceholderProperty) generateValue(settings c.Settings) string {
	return ""
}

// ~~~~~~~~~~~~~ Object ~~~~~~~~~~~~~

type ObjectProperty struct {
	name            string
	propertyType    string
	nextLevel       int
	innerProperties []Property
	variant         v.Variant
}

func (p ObjectProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p ObjectProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p ObjectProperty) getName() string {
	return p.name
}

func (p ObjectProperty) getUnderlyingType() string {
	return p.propertyType
}

func (p ObjectProperty) getType() string {
	return p.propertyType
}

func (p ObjectProperty) getVariant() v.Variant {
	return p.variant
}

// TODO nested objects
func (p ObjectProperty) generateValue(settings c.Settings) string {
	res := "{\n"
	for _, innerProp := range p.innerProperties {
		res += fmt.Sprintf("%s%s: %s,\n", getIndent(&settings, 2), innerProp.getName(), innerProp.generateValue(settings))
	}
	res += "},\n"
	return res
}

// ~~~~~~~~~~~~~ String ~~~~~~~~~~~~~

type StringProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p StringProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p StringProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p StringProperty) getName() string {
	return p.name
}

func (p StringProperty) getUnderlyingType() string {
	return p.propertyType
}

func (p StringProperty) getType() string {
	return p.propertyType
}

func (p StringProperty) getVariant() v.Variant {
	return p.variant
}

func (p StringProperty) generateValue(settings c.Settings) string {
	variant := string(p.variant)
	if len(predefinedValues[variant]) > 0 {
		values := predefinedValues[variant]
		return fmt.Sprintf("`%s`", values[utils.Random(0, len(values)-1)])
	}
	length := 39 // lorem(39) -> Lorem ipsum, dolor sit amet consectetur

	if p.variant == "content" {
		length = len(data.Content) - 1
	}

	text := data.Content

	return fmt.Sprintf("`%s`", text[0:length])
}

// ~~~~~~~~~~~~~ Number ~~~~~~~~~~~~~

type NumberProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p NumberProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p NumberProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p NumberProperty) getName() string {
	return p.name
}

func (p NumberProperty) getUnderlyingType() string {
	return "number"
}

func (p NumberProperty) getType() string {
	return p.propertyType
}

func (p NumberProperty) getVariant() v.Variant {
	return p.variant
}

func (p NumberProperty) generateValue(settings c.Settings) string {
	variant := string(p.variant)

	min := 0
	max := 100
	numRange := strings.Split(variant, " ")

	if len(numRange) == 1 {
		max = utils.ParseInt(numRange[0], max)
	} else if len(numRange) >= 2 {
		min = utils.ParseInt(numRange[0], min)
		max = utils.ParseInt(numRange[1], max)
	}

	return fmt.Sprint(utils.Random(min, max))
}

// ~~~~~~~~~~~~~ Boolean ~~~~~~~~~~~~~

type BooleanProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p BooleanProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p BooleanProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p BooleanProperty) getName() string {
	return p.name
}

func (p BooleanProperty) getUnderlyingType() string {
	return p.propertyType
}

func (p BooleanProperty) getType() string {
	return p.propertyType
}

func (p BooleanProperty) getVariant() v.Variant {
	return p.variant
}

func (p BooleanProperty) generateValue(settings c.Settings) string {
	if utils.Random(0, 100) >= 50 {
		return "true"
	} else {
		return "false"
	}
}

// ~~~~~~~~~~~~~ Image ~~~~~~~~~~~~~

type ImageProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p ImageProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p ImageProperty) getName() string {
	return p.name
}

func (p ImageProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p ImageProperty) getUnderlyingType() string {
	return "string"
}

func (p ImageProperty) getType() string {
	return p.propertyType
}

func (p ImageProperty) getVariant() v.Variant {
	return p.variant
}

func (p ImageProperty) generateValue(settings c.Settings) string {
	variant := string(p.variant)

	dimensions := strings.Split(variant, "x")
	width := dimensions[0]
	height := dimensions[1]

	return fmt.Sprintf("`https://unsplash.it/%s/%s`", width, height)
}

// ~~~~~~~~~~~~~ Date ~~~~~~~~~~~~~

var underlyingDateTypes = map[v.Variant]string{
	"dateTime":      "string",
	"timestamp":     "number",
	"day":           "number",
	"month":         "number",
	"year":          "number",
	"date":          "number",
	"object":        "Date",
	"personal data": "string",
}

type DateProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p DateProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p DateProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p DateProperty) getName() string {
	return p.name
}

func (p DateProperty) getUnderlyingType() string {
	for k, v := range underlyingDateTypes {
		if k == p.getVariant() {
			return v
		}
	}
	return p.propertyType
}

func (p DateProperty) getType() string {
	return p.propertyType
}

func (p DateProperty) getVariant() v.Variant {
	return p.variant
}

func (p DateProperty) generateValue(settings c.Settings) string {
	variant := p.variant

	YEAR_IN_DAYS := 365
	YEAR_IN_MONTHS := 12
	MONTH_IN_DAYS := 31
	TEN_YEARS := 10

	switch variant {
	case v.DateTime:
		return fmt.Sprintf("`%s`", time.Now().AddDate(0, 0, -1*rand.Intn(YEAR_IN_DAYS+1)).Format("2.1.2006"))
	case v.Timestamp:
		return fmt.Sprintf("%d", time.Now().AddDate(0, 0, -1*rand.Intn(YEAR_IN_DAYS+1)).Unix()*1000) // unix time to js timestamp
	case v.Day:
		return fmt.Sprintf("%d", time.Now().AddDate(0, 0, -1*rand.Intn(MONTH_IN_DAYS+1)).Day())
	case v.Month:
		return fmt.Sprintf("%d", time.Now().AddDate(0, -1*rand.Intn(YEAR_IN_MONTHS+1), 0).Month())
	case v.Year:
		return fmt.Sprintf("%d", time.Now().AddDate(-1*rand.Intn(TEN_YEARS+1), 0, 0).Year())
	case v.DateObject:
		return "new Date()"
	default:
		return fmt.Sprintf("`%s`", time.Now().AddDate(0, 0, -1*rand.Intn(YEAR_IN_DAYS+1)).Format("2.1.2006"))
	}
}

// ~~~~~~~~~~~~~ ID ~~~~~~~~~~~~~

type IdProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p IdProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p IdProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p IdProperty) getName() string {
	return p.name
}

func (p IdProperty) getUnderlyingType() string {
	return "string"
}

func (p IdProperty) getType() string {
	return p.propertyType
}

func (p IdProperty) getVariant() v.Variant {
	return p.variant
}

func (p IdProperty) generateValue(settings c.Settings) string {
	switch p.variant {
	case v.UUID:
		uuid, err := utils.UUIDv4()
		if err != nil {
			return "`60a0c89e-84c9-41c8-a731-31f6e0a31138`"
		}
		return fmt.Sprintf("`%s`", uuid)

	case v.NanoID:
		id, err := utils.NanoID()
		if err != nil {
			return "`V1StGXR8_Z5jdHi6B-myT`"
		}
		return fmt.Sprintf("`%s`", id)
	default:
		return fmt.Sprint(utils.Random(1_000_000, 9_999_999))
	}
}

// ~~~~~~~~~~~~~ Enum ~~~~~~~~~~~~~

type EnumProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p EnumProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p EnumProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p EnumProperty) getName() string {
	return p.name
}

func (p EnumProperty) getUnderlyingType() string {
	res := "("
	wordSet := parseStringEnum(strings.Split(string(p.getVariant()), " "))

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

func (p EnumProperty) getType() string {
	return p.propertyType
}

func (p EnumProperty) getVariant() v.Variant {
	return p.variant
}

func (p EnumProperty) generateValue(settings c.Settings) string {
	variant := string(p.variant)

	wordSet := strings.Split(data.Content[0:39], " ")
	v := strings.Split(variant, " ")
	if len(v) != 0 {
		wordSet = v
	}
	wordSet = parseStringEnum(wordSet)
	randStr := wordSet[utils.Random(0, len(wordSet)-1)]

	return fmt.Sprintf("`%s`", randStr)
}

// ~~~~~~~~~~~~~ Null, Undefined ~~~~~~~~~~~~~

type NullProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p NullProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p NullProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p NullProperty) getName() string {
	return p.name
}

func (p NullProperty) getUnderlyingType() string {
	return "null"
}

func (p NullProperty) getType() string {
	return p.propertyType
}

func (p NullProperty) getVariant() v.Variant {
	return p.variant
}

func (p NullProperty) generateValue(settings c.Settings) string {
	return "null"
}

type UndefinedProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p UndefinedProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p UndefinedProperty) getName() string {
	return p.name

}
func (p UndefinedProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p UndefinedProperty) getUnderlyingType() string {
	return "undefined"
}

func (p UndefinedProperty) getType() string {
	return p.propertyType
}

func (p UndefinedProperty) getVariant() v.Variant {
	return p.variant
}

func (p UndefinedProperty) generateValue(settings c.Settings) string {
	return "null"
}

// ~~~~~ Personal Data ~~~~~~
type PersonalDataProperty struct {
	name         string
	propertyType string
	variant      v.Variant
}

func (p PersonalDataProperty) clone(other ...string) Property {
	return CreateProperty(p.name, p.propertyType, other...)
}

func (p PersonalDataProperty) setVariant(variant v.Variant) {
	p.variant = variant
}

func (p PersonalDataProperty) getName() string {
	return p.name
}

func (p PersonalDataProperty) getUnderlyingType() string {
	if p.variant == "phone number" {
		return "number"
	}
	return "string"
}

func (p PersonalDataProperty) getType() string {
	return p.propertyType
}

func (p PersonalDataProperty) getVariant() v.Variant {
	if len(predefinedValues[string(p.variant)]) != 0 || p.variant == "phone number" {
		return p.variant
	}
	return "zip-code"
}

func (p PersonalDataProperty) generateValue(settings c.Settings) string {
	variant := string(p.variant)

	if len(predefinedValues[variant]) > 0 {
		values := predefinedValues[variant]
		return fmt.Sprintf("`%s`", values[utils.Random(0, len(values)-1)])
	}

	if p.variant == "phone number" {
		res := ""
		for range 9 {
			res += fmt.Sprint(utils.Random(0, 9))
		}
		return res
	}

	//- zip-code
	res := "`"

	for _, char := range variant {
		switch char {
		case '-':
			res += "-"
		case '_':
			res += " "
		case 'd':
			res += fmt.Sprint(utils.Random(0, 9))
		case 'a':
			res += strings.ToLower(string(utils.Alphabet[utils.Random(0, 25)]))
		case 'A':
			res += string(utils.Alphabet[utils.Random(0, 25)])
		}
	}

	res += "`"

	return res
}

func CreateProperty(name, propType string, other ...string) Property {
	var variant v.Variant
	if len(other) == 1 {
		variant = v.Variant(other[0])
	}

	var res Property
	switch propType {
	case "string":
		res = StringProperty{name: name, propertyType: propType, variant: variant}
	case "number":
		res = NumberProperty{name: name, propertyType: propType, variant: variant}
	case "boolean":
		res = BooleanProperty{name: name, propertyType: propType, variant: variant}
	case "date":
		res = DateProperty{name: name, propertyType: propType, variant: variant}
	case "img":
		res = ImageProperty{name: name, propertyType: propType, variant: variant}
	case "id":
		res = IdProperty{name: name, propertyType: propType, variant: variant}
	case "string enum":
		res = EnumProperty{name: name, propertyType: propType, variant: variant}
	case "null":
		res = NullProperty{name: name, propertyType: propType, variant: variant}
	case "undefined":
		res = UndefinedProperty{name: name, propertyType: propType, variant: variant}
	case "personal data":
		res = PersonalDataProperty{name: name, propertyType: propType, variant: variant}
	}
	return res
}
