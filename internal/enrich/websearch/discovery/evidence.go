package discovery

import (
	"strings"

	"osint/internal/enrich/model"
	"osint/internal/enrich/provenance"
)

func SearchEvidence(
	input model.Input,
	item SearchResult,
) model.Evidence {
	title :=
		strings.TrimSpace(
			item.Title,
		)

	snippet :=
		strings.TrimSpace(
			item.Snippet,
		)

	text :=
		searchEvidenceText(
			title,
			snippet,
		)

	value :=
		model.Evidence{
			Anchor: provenance.Anchor(
				input,
			),

			Source: searchEvidenceSource(
				item,
			),

			Query: strings.TrimSpace(
				item.Query,
			),

			Engine: strings.TrimSpace(
				item.Engine,
			),

			Title: title,

			//
			// This MUST be exactly the same text
			// passed to TextInput.
			//
			Text: text,

			URL: strings.TrimSpace(
				item.URL,
			),
		}

	value.ID =
		provenance.ID(
			"search",
			value,
		)

	return value
}

func searchEvidenceText(
	title string,
	snippet string,
) string {
	title =
		strings.TrimSpace(
			title,
		)

	snippet =
		strings.TrimSpace(
			snippet,
		)

	if title == "" {

		return snippet
	}

	if snippet == "" {

		return title
	}

	titleFields :=
		strings.Fields(
			title,
		)

	snippetFields :=
		strings.Fields(
			snippet,
		)

	if len(titleFields) == 0 {

		return snippet
	}

	if len(snippetFields) == 0 {

		return title
	}

	maxOverlap :=
		len(titleFields)

	if len(snippetFields) <
		maxOverlap {

		maxOverlap =
			len(snippetFields)
	}

	//
	// Search engines often repeat the trailing
	// entity from the title at the beginning of
	// the description.
	//
	// Example:
	//
	// Title:
	//     Η εταιρεία – some company
	//
	// Snippet:
	//     Some company Way was founded by...
	//
	// Blind concatenation produces:
	//
	//     Some company Some company
	//
	// Find the longest word-for-word overlap
	// between the end of the title and the
	// beginning of the snippet and include that
	// text only once.
	//
	overlap :=
		0

	for size :=
		maxOverlap; size > 0; size-- {

		titleStart :=
			len(titleFields) -
				size

		matched :=
			true

		for index := 0; index < size; index++ {

			if !strings.EqualFold(
				titleFields[titleStart+
					index],
				snippetFields[index],
			) {

				matched =
					false

				break
			}
		}

		if !matched {

			continue
		}

		overlap =
			size

		break
	}

	if overlap == 0 {

		return strings.TrimSpace(
			title +
				" " +
				snippet,
		)
	}

	if overlap >=
		len(snippetFields) {

		return title
	}

	return strings.TrimSpace(
		title +
			" " +
			strings.Join(
				snippetFields[overlap:],
				" ",
			),
	)
}

func searchEvidenceSource(
	item SearchResult,
) string {
	source :=
		"WebSearch Discovery"

	engine :=
		strings.TrimSpace(
			item.Engine,
		)

	if engine != "" {

		source +=
			" [" +
				engine +
				"]"
	}

	query :=
		strings.TrimSpace(
			item.Query,
		)

	if query != "" {

		source +=
			" | Query: " +
				query
	}

	return source
}
