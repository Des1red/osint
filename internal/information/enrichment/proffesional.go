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

	if len(
		result.UnattributedOrganizations,
	) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Unattributed Organizations",

					Print: func(
						prefix string,
					) {
						printOrganizations(
							prefix,
							result.UnattributedOrganizations,
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

	if len(
		result.UnattributedEmployment,
	) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Unattributed Employment",

					Print: func(
						prefix string,
					) {
						printEmployment(
							prefix,
							result.UnattributedEmployment,
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

	if len(
		result.UnattributedEducation,
	) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Unattributed Education",

					Print: func(
						prefix string,
					) {
						printEducation(
							prefix,
							result.UnattributedEducation,
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

		fields :=
			appendProvenanceFields(
				nil,
				organization.Source,
				organization.Evidence,
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
			[]output.TreeField{}

		if strings.TrimSpace(
			job.Organization,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Organization",

						Value: job.Organization,
					},
				)
		}

		if strings.TrimSpace(
			job.StartDate,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Start Date",

						Value: job.StartDate,
					},
				)
		}

		if strings.TrimSpace(
			job.EndDate,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "End Date",

						Value: job.EndDate,
					},
				)
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

		if strings.TrimSpace(
			job.Summary,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Summary",

						Value: job.Summary,
					},
				)
		}

		fields =
			appendProvenanceFields(
				fields,
				job.Source,
				job.Evidence,
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
			[]output.TreeField{}

		degrees :=
			strings.Join(
				school.Degrees,
				", ",
			)

		if strings.TrimSpace(
			degrees,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Degrees",

						Value: degrees,
					},
				)
		}

		majors :=
			strings.Join(
				school.Majors,
				", ",
			)

		if strings.TrimSpace(
			majors,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Majors",

						Value: majors,
					},
				)
		}

		if strings.TrimSpace(
			school.StartDate,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Start Date",

						Value: school.StartDate,
					},
				)
		}

		if strings.TrimSpace(
			school.EndDate,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "End Date",

						Value: school.EndDate,
					},
				)
		}

		if strings.TrimSpace(
			school.Summary,
		) != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Summary",

						Value: school.Summary,
					},
				)
		}

		fields =
			appendProvenanceFields(
				fields,
				school.Source,
				school.Evidence,
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
