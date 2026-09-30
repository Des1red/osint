package model

type DirectoryResult struct {
	Provider string

	Query string

	Entries []DirectoryEntry
}

type DirectoryEntry struct {
	Title string

	Fields []DirectoryField

	SourceURL string
}

type DirectoryField struct {
	Name string

	Value string
}
