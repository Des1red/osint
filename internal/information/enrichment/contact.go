package enrichment

import (
	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func printContactAndLocation(
	result knowledge.EnrichmentResult,
) {
	var sections []treeSection

	if len(result.Phones) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Phones",

					Print: func(
						prefix string,
					) {
						printPhones(
							prefix,
							result.Phones,
						)
					},
				},
			)
	}

	if len(result.UnattributedPhones) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Unattributed Phones",

					Print: func(
						prefix string,
					) {
						printPhones(
							prefix,
							result.UnattributedPhones,
						)
					},
				},
			)
	}

	if len(result.Locations) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Locations",

					Print: func(
						prefix string,
					) {
						printLocations(
							prefix,
							result.Locations,
						)
					},
				},
			)
	}

	printTreeSection(
		"Contact & Location",
		"------------------",
		sections,
	)
}

func printPhones(
	prefix string,
	phones []knowledge.PhoneReference,
) {
	for index, phone := range phones {

		last :=
			index ==
				len(phones)-1

		output.TreeItem(
			prefix,
			last,
			phone.Phone,
		)

		if !output.Full() {

			continue
		}

		fields :=
			[]output.TreeField{
				{
					Name: "Discovered From",

					Value: phone.Source,
				},
			}

		evidence :=
			evidenceSummary(
				phone.Evidence,
			)

		if evidence != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Evidence",

						Value: evidence,
					},
				)
		}

		output.TreeFields(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			fields,
		)
	}
}

func printLocations(
	prefix string,
	locations []knowledge.LocationReference,
) {
	for index, location := range locations {

		last :=
			index ==
				len(locations)-1

		output.TreeItem(
			prefix,
			last,
			location.Location,
		)

		if !output.Full() {

			continue
		}

		output.TreeFields(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			[]output.TreeField{
				{
					Name: "Discovered From",

					Value: location.Source,
				},
			},
		)
	}
}
