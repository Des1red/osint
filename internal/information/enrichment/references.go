package enrichment

import (
	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func printReferences(
	result knowledge.EnrichmentResult,
) {
	var sections []treeSection

	if len(result.Links) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "External Links",

					Print: func(
						prefix string,
					) {
						printLinks(
							prefix,
							result.Links,
						)
					},
				},
			)
	}

	if len(
		result.UnattributedLinks,
	) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Unattributed External Links",

					Print: func(
						prefix string,
					) {
						printLinks(
							prefix,
							result.UnattributedLinks,
						)
					},
				},
			)
	}

	printTreeSection(
		"References",
		"----------",
		sections,
	)
}

func printLinks(
	prefix string,
	links []knowledge.ExternalLink,
) {
	for index, link := range links {

		last :=
			index ==
				len(links)-1

		output.TreeItem(
			prefix,
			last,
			link.URL,
		)

		fields :=
			[]output.TreeField{}

		if link.ResolvedURL != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Resolved",

						Value: link.ResolvedURL,
					},
				)
		}

		if link.Category != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Category",

						Value: link.Category,
					},
				)
		}

		fields =
			appendProvenanceFields(
				fields,
				link.Source,
				link.Evidence,
			)

		output.TreeFields(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			fields,
		)
	}
}
