package gr11888

type Result struct {
	Query string

	Entries []Entry
}

type Entry struct {
	Title string

	Fields []Field

	SourceURL string
}

type Field struct {
	Name string

	Value string
}
