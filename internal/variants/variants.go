package variants

type Variant string

type VariantData struct {
	Key   Variant
	Value string
}

type VariantSet string

const (
	DateSet     VariantSet = "date"
	ImgSet      VariantSet = "img"
	OutputSet   VariantSet = "output"
	LanguageSet VariantSet = "lang"
)

// - generator options
const (
	HorizontalImg     Variant = "300x500"
	VerticalImg       Variant = "500x300"
	ProfilePictureImg Variant = "100x100"
	ArticleImg        Variant = "600x400"
	Banner            Variant = "600x240"
	Custom            Variant = "custom"
)

var ImgVariants = []VariantData{
	{HorizontalImg, "Horizontal (default) 300x500"},
	{VerticalImg, "Vertical 500x300"},
	{ProfilePictureImg, "Profile Picture 100x100"},
	{ArticleImg, "Article 600x400"},
	{Banner, "Banner 600x240"},
	{Custom, "Custom"},
}

const (
	DateTime   Variant = "dateTime"
	Timestamp  Variant = "timestamp"
	Day        Variant = "day"
	Month      Variant = "month"
	Year       Variant = "year"
	DateObject Variant = "object"
)

var DateVariants = []VariantData{
	{DateTime, "dateTime: e.g. 27.02.2024"},
	{Timestamp, "timestamp: e.g. 1718051654"},
	{Day, "day: 1-31"},
	{Month, "month: 1-12"},
	{Year, "year: current year (-10)"},
	{DateObject, "object: new Date()"},
}

//- configuration options

const (
	Terminal Variant = "terminal"
	File     Variant = "file"
)

var OutputVariants = []VariantData{
	{Terminal, "In terminal"},
	{File, "Output file"},
}

const (
	TypeScript Variant = "typescript"
	JavaScript Variant = "javascript"
	JSON       Variant = "json"
	Config     Variant = "config"
)

var LanguageVariants = []VariantData{
	{TypeScript, "TypeScript"},
	{JavaScript, "JavaScript"},
	{JSON, "JSON"},
}
