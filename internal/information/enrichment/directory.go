package enrichment

import (
	"fmt"
	"strings"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func printDirectories(
	directories []knowledge.DirectoryResult,
) {
	if len(directories) == 0 {

		return
	}

	fmt.Fprintln(
		output.Writer(),
	)

	fmt.Fprintln(
		output.Writer(),
		"Directory Results",
	)

	fmt.Fprintln(
		output.Writer(),
		"=================",
	)

	for index, directory := range directories {

		last :=
			index ==
				len(directories)-1

		provider :=
			strings.TrimSpace(
				directory.Provider,
			)

		if provider == "" {

			provider =
				"Directory"
		}

		output.TreeItem(
			"",
			last,
			provider,
		)

		printDirectory(
			output.TreeChildPrefix(
				"",
				last,
			),
			directory,
		)
	}
}

func printDirectory(
	prefix string,
	directory knowledge.DirectoryResult,
) {
	query :=
		strings.TrimSpace(
			directory.Query,
		)

	hasEntries :=
		len(directory.Entries) > 0

	if query != "" {

		output.TreeValue(
			prefix,
			!hasEntries,
			"Query",
			query,
		)
	}

	if !hasEntries {

		return
	}

	output.TreeItem(
		prefix,
		true,
		"Entries",
	)

	printDirectoryEntries(
		output.TreeChildPrefix(
			prefix,
			true,
		),
		directory.Entries,
	)
}

func printDirectoryEntries(
	prefix string,
	entries []knowledge.DirectoryEntry,
) {
	for index, entry := range entries {

		last :=
			index ==
				len(entries)-1

		title :=
			strings.TrimSpace(
				entry.Title,
			)

		if title == "" {

			title =
				"Entry"
		}

		output.TreeItem(
			prefix,
			last,
			title,
		)

		printDirectoryEntry(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			entry,
		)
	}
}

func printDirectoryEntry(
	prefix string,
	entry knowledge.DirectoryEntry,
) {
	fields :=
		make(
			[]output.TreeField,
			0,
			len(entry.Fields)+1,
		)

	for _, field := range entry.Fields {

		name :=
			strings.TrimSpace(
				field.Name,
			)

		value :=
			strings.TrimSpace(
				field.Value,
			)

		if name == "" ||
			value == "" {

			continue
		}

		fields =
			append(
				fields,
				output.TreeField{
					Name: name,

					Value: value,
				},
			)
	}

	// sourceURL :=
	// 	strings.TrimSpace(
	// 		entry.SourceURL,
	// 	)

	// if sourceURL != "" {

	// 	fields =
	// 		append(
	// 			fields,
	// 			output.TreeField{
	// 				Name: "Source",

	// 				Value: sourceURL,
	// 			},
	// 		)
	// }

	output.TreeFields(
		prefix,
		fields,
	)
}
