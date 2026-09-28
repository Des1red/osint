package enrichment

import (
	"fmt"
	"strings"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

const maxPrintedPeopleSearchResults = 5

func printPeople(
	people []knowledge.PersonReference,
) {
	if len(people) == 0 {

		return
	}

	fmt.Fprintln(
		output.Writer(),
	)

	fmt.Fprintln(
		output.Writer(),
		"People",
	)

	fmt.Fprintln(
		output.Writer(),
		"------",
	)

	for index, person := range people {

		name :=
			strings.TrimSpace(
				person.Name,
			)

		if name == "" {

			continue
		}

		last :=
			index ==
				len(people)-1

		output.TreeItem(
			"",
			last,
			name,
		)

		printPersonDetails(
			output.TreeChildPrefix(
				"",
				last,
			),
			person,
		)
	}
}

func printPersonDetails(
	prefix string,
	person knowledge.PersonReference,
) {
	fields :=
		[]output.TreeField{}

	if person.RelationVerified {

		fields =
			append(
				fields,
				output.TreeField{
					Name: "Relationship",

					Value: person.Relation,
				},
				output.TreeField{
					Name: "Relationship Verified",

					Value: "true",
				},
			)

		if output.Full() {

			if strings.TrimSpace(
				person.RelationSource,
			) != "" {

				fields =
					append(
						fields,
						output.TreeField{
							Name: "Relationship Evidence Source",

							Value: person.RelationSource,
						},
					)
			}

			relationEvidence :=
				evidenceSummary(
					person.RelationEvidence,
				)

			if relationEvidence != "" {

				fields =
					append(
						fields,
						output.TreeField{
							Name: "Relationship Evidence",

							Value: relationEvidence,
						},
					)
			}
		}

	} else {

		fields =
			append(
				fields,
				output.TreeField{
					Name: "Relationship",

					Value: "Unverified",
				},
				output.TreeField{
					Name: "Relationship Verified",

					Value: "false",
				},
			)
	}

	if output.Full() {

		if strings.TrimSpace(
			person.Source,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Discovered From",

						Value: person.Source,
					},
				)
		}

		personEvidence :=
			evidenceSummary(
				person.Evidence,
			)

		if personEvidence != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Evidence",

						Value: personEvidence,
					},
				)
		}
	}

	sections :=
		personSections(
			person,
		)

	hasSearchResults :=
		len(person.SearchResults) > 0

	hasChildren :=
		len(sections) > 0 ||
			hasSearchResults

	filtered :=
		make(
			[]output.TreeField,
			0,
			len(fields),
		)

	for _, field := range fields {

		if strings.TrimSpace(
			field.Value,
		) == "" {

			continue
		}

		filtered =
			append(
				filtered,
				field,
			)
	}

	for index, field := range filtered {

		last :=
			index ==
				len(filtered)-1 &&
				!hasChildren

		output.TreeValue(
			prefix,
			last,
			field.Name,
			field.Value,
		)
	}

	for index, section := range sections {

		last :=
			index ==
				len(sections)-1 &&
				!hasSearchResults

		output.TreeItem(
			prefix,
			last,
			section.Title,
		)

		section.Print(
			output.TreeChildPrefix(
				prefix,
				last,
			),
		)
	}

	if !hasSearchResults {

		return
	}

	output.TreeItem(
		prefix,
		true,
		"WebSearch",
	)

	printPeopleSearchResults(
		output.TreeChildPrefix(
			prefix,
			true,
		),
		person.SearchResults,
	)
}

func personSections(
	person knowledge.PersonReference,
) []treeSection {
	var sections []treeSection

	if len(person.Usernames) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Usernames",

					Print: func(
						prefix string,
					) {
						printUsernames(
							prefix,
							person.Usernames,
						)
					},
				},
			)
	}

	if len(person.RelatedAccounts) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Related Accounts",

					Print: func(
						prefix string,
					) {
						printRelatedAccounts(
							prefix,
							person.RelatedAccounts,
						)
					},
				},
			)
	}

	if len(person.Emails) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Emails",

					Print: func(
						prefix string,
					) {
						printEmails(
							prefix,
							person.Emails,
						)
					},
				},
			)
	}

	if len(person.Socials) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Social Profiles",

					Print: func(
						prefix string,
					) {
						printSocials(
							prefix,
							person.Socials,
						)
					},
				},
			)
	}

	if len(person.Phones) > 0 {

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
							person.Phones,
						)
					},
				},
			)
	}

	if len(person.Locations) > 0 {

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
							person.Locations,
						)
					},
				},
			)
	}

	if len(person.Organizations) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Organizations",

					Print: func(
						prefix string,
					) {
						printOrganizations(
							prefix,
							person.Organizations,
						)
					},
				},
			)
	}

	if len(person.Employment) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Employment",

					Print: func(
						prefix string,
					) {
						printEmployment(
							prefix,
							person.Employment,
						)
					},
				},
			)
	}

	if len(person.Education) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Education",

					Print: func(
						prefix string,
					) {
						printEducation(
							prefix,
							person.Education,
						)
					},
				},
			)
	}

	if len(person.Links) > 0 {

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
							person.Links,
						)
					},
				},
			)
	}

	return sections
}

func printPeopleSearchResults(
	prefix string,
	results []knowledge.RelativeSearchResult,
) {
	limit :=
		len(results)

	if limit >
		maxPrintedPeopleSearchResults {

		limit =
			maxPrintedPeopleSearchResults
	}

	for index := 0; index < limit; index++ {

		result :=
			results[index]

		title :=
			strings.TrimSpace(
				result.Title,
			)

		if title == "" {

			title =
				strings.TrimSpace(
					result.URL,
				)
		}

		if title == "" {

			continue
		}

		last :=
			index ==
				limit-1

		output.TreeItem(
			prefix,
			last,
			title,
		)

		fields :=
			[]output.TreeField{}

		if output.Full() {

			if strings.TrimSpace(
				result.Engine,
			) != "" {

				fields =
					append(
						fields,
						output.TreeField{
							Name: "Engine",

							Value: result.Engine,
						},
					)
			}

			if strings.TrimSpace(
				result.Query,
			) != "" {

				fields =
					append(
						fields,
						output.TreeField{
							Name: "Query",

							Value: result.Query,
						},
					)
			}

			if strings.TrimSpace(
				result.EvidenceID,
			) != "" {

				fields =
					append(
						fields,
						output.TreeField{
							Name: "Evidence",

							Value: result.EvidenceID,
						},
					)
			}
		}

		if strings.TrimSpace(
			result.URL,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "URL",

						Value: result.URL,
					},
				)
		}

		if strings.TrimSpace(
			result.Snippet,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Description",

						Value: result.Snippet,
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
