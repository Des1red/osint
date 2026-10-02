package directory

import (
	"strings"

	"osint/internal/enrich/directory/gr11888"
	"osint/internal/enrich/model"
	"osint/internal/logger"
)

func Enrich(
	input model.Input,
) (
	[]model.DirectoryResult,
	error,
) {
	surname :=
		strings.TrimSpace(
			input.Surname,
		)

	if surname == "" {

		return nil,
			nil
	}

	logger.HeaderStart(
		"Directory Enrichment",
	)

	defer logger.HeaderEnd(
		"Directory Enrichment",
	)

	logger.Info(
		"Searching 11888...",
	)

	result,
		err :=
		gr11888.Search(
			surname,
		)

	if err != nil {

		return nil,
			err
	}

	if len(result.Entries) == 0 {

		logger.Info(
			"No 11888 directory results found.",
		)

		return nil,
			nil
	}

	logger.Info(
		"11888 directory search complete.",
	)

	return []model.DirectoryResult{
			toDirectoryResult(
				result,
			),
		},
		nil
}

func toDirectoryResult(
	value gr11888.Result,
) model.DirectoryResult {
	result :=
		model.DirectoryResult{
			Provider: "Greece",

			Query: value.Query,

			Entries: make(
				[]model.DirectoryEntry,
				0,
				len(value.Entries),
			),
		}

	for _, entry := range value.Entries {

		item :=
			model.DirectoryEntry{
				Title: entry.Title,

				SourceURL: entry.SourceURL,

				Fields: make(
					[]model.DirectoryField,
					0,
					len(entry.Fields),
				),
			}

		for _, field := range entry.Fields {

			item.Fields =
				append(
					item.Fields,
					model.DirectoryField{
						Name: field.Name,

						Value: field.Value,
					},
				)
		}

		result.Entries =
			append(
				result.Entries,
				item,
			)
	}

	return result
}
