package extract

import (
	"strings"

	"osint/internal/enrich/model"
)

func extractEmployment(
	input model.Input,
) []model.EmploymentReference {
	var result []model.EmploymentReference

	for _, value := range input.Employment {

		title :=
			strings.TrimSpace(
				value.Title,
			)

		organization :=
			strings.TrimSpace(
				value.Organization,
			)

		if title == "" &&
			organization == "" {

			continue
		}

		result =
			append(
				result,
				model.EmploymentReference{
					Title: title,

					Organization: organization,

					StartDate: strings.TrimSpace(
						value.StartDate,
					),

					EndDate: strings.TrimSpace(
						value.EndDate,
					),

					Current: value.Current,

					Summary: strings.TrimSpace(
						value.Summary,
					),

					Source: value.Source,

					Evidence: wholeEvidenceReference(
						value.EvidenceID,
					),
				},
			)
	}

	return uniqueEmployment(
		result,
	)
}

func uniqueEmployment(
	values []model.EmploymentReference,
) []model.EmploymentReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.EmploymentReference

	for _, value := range values {

		key :=
			strings.ToLower(
				strings.TrimSpace(
					value.Title,
				) +
					":" +
					strings.TrimSpace(
						value.Organization,
					) +
					":" +
					strings.TrimSpace(
						value.StartDate,
					) +
					":" +
					strings.TrimSpace(
						value.EndDate,
					) +
					":" +
					strings.TrimSpace(
						value.Source,
					),
			)

		if index,
			exists :=
			indexes[key]; exists {

			result[index].Evidence =
				mergeEvidence(
					result[index].Evidence,
					value.Evidence,
				)

			if strings.TrimSpace(
				result[index].Summary,
			) == "" {

				result[index].Summary =
					value.Summary
			}

			if !result[index].Current &&
				value.Current {

				result[index].Current =
					true
			}

			continue
		}

		indexes[key] =
			len(result)

		result =
			append(
				result,
				value,
			)
	}

	return result
}
