package enrichment

import (
	"strings"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func printProfessional(
	result knowledge.EnrichmentResult,
) {
	var sections []treeSection

	if len(result.Organizations) > 0 {

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
							result.Organizations,
						)
					},
				},
			)
	}

	if len(result.Employment) > 0 {

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
							result.Employment,
						)
					},
				},
			)
	}

	if len(result.Education) > 0 {

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
							result.Education,
						)
					},
				},
			)
	}

	printTreeSection(
		"Professional",
		"------------",
		sections,
	)
}

func printOrganizations(
	prefix string,
	organizations []knowledge.OrganizationReference,
) {
	for index, organization := range organizations {

		last :=
			index ==
				len(organizations)-1

		output.TreeItem(
			prefix,
			last,
			organization.Organization,
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

					Value: organization.Source,
				},
			},
		)
	}
}

func printEmployment(
	prefix string,
	employment []knowledge.EmploymentReference,
) {
	for index, job := range employment {

		last :=
			index ==
				len(employment)-1

		title :=
			strings.TrimSpace(
				job.Title,
			)

		if title == "" {
			title =
				"Employment"
		}

		output.TreeItem(
			prefix,
			last,
			title,
		)

		fields :=
			[]output.TreeField{
				{
					Name: "Organization",

					Value: job.Organization,
				},
				{
					Name: "Start Date",

					Value: job.StartDate,
				},
				{
					Name: "End Date",

					Value: job.EndDate,
				},
			}

		if job.Current {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Current",

						Value: "true",
					},
				)
		}

		fields =
			append(
				fields,
				output.TreeField{
					Name: "Summary",

					Value: job.Summary,
				},
			)

		if output.Full() {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Discovered From",

						Value: job.Source,
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

func printEducation(
	prefix string,
	education []knowledge.EducationReference,
) {
	for index, school := range education {

		last :=
			index ==
				len(education)-1

		name :=
			strings.TrimSpace(
				school.School,
			)

		if name == "" {
			name =
				"Education"
		}

		output.TreeItem(
			prefix,
			last,
			name,
		)

		fields :=
			[]output.TreeField{
				{
					Name: "Degrees",

					Value: strings.Join(
						school.Degrees,
						", ",
					),
				},
				{
					Name: "Majors",

					Value: strings.Join(
						school.Majors,
						", ",
					),
				},
				{
					Name: "Start Date",

					Value: school.StartDate,
				},
				{
					Name: "End Date",

					Value: school.EndDate,
				},
				{
					Name: "Summary",

					Value: school.Summary,
				},
			}

		if output.Full() {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Discovered From",

						Value: school.Source,
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
