package options

type Option string

type OptionData struct {
	Key   Option
	Value string
}

type OptionSet string

const (
	DateSet     OptionSet = "date"
	ImgSet      OptionSet = "img"
	OutputSet   OptionSet = "output"
	LanguageSet OptionSet = "lang"
)

// - generator options
const (
	HorizontalImg     Option = "300x500"
	VerticalImg       Option = "500x300"
	ProfilePictureImg Option = "100x100"
	ArticleImg        Option = "600x400"
	Banner            Option = "600x240"
	Custom            Option = "custom"
)

var ImgOptions = []OptionData{
	{HorizontalImg, "Horizontal (default) 300x500"},
	{VerticalImg, "Vertical 500x300"},
	{ProfilePictureImg, "Profile Picture 100x100"},
	{ArticleImg, "Article 600x400"},
	{Banner, "Banner 600x240"},
	{Custom, "Custom"},
}

const (
	DateTime   Option = "dateTime"
	Timestamp  Option = "timestamp"
	Day        Option = "day"
	Month      Option = "month"
	Year       Option = "year"
	DateObject Option = "object"
)

var DateOptions = []OptionData{
	{DateTime, "dateTime: e.g. 27.02.2024"},
	{Timestamp, "timestamp: e.g. 1718051654"},
	{Day, "day: 0-31"},
	{Month, "month: 0-12"},
	{Year, "year: current year"},
	{DateObject, "object: new Date()"},
}

//- configuration options

const (
	Terminal Option = "terminal"
	File     Option = "file"
)

var OutputOptions = []OptionData{
	{Terminal, "In terminal"},
	{File, "Output file"},
}

const (
	TypeScript Option = "typescript"
	JavaScript Option = "javascript"
)

var LanguageOptions = []OptionData{
	{TypeScript, "TypeScript"},
	{JavaScript, "JavaScript"},
}
