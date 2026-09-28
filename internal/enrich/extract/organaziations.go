package extract

import (
	"strings"

	"osint/internal/enrich/model"
)

func extractOrganizations(
	input model.Input,
) []model.OrganizationReference {
	var organizations []model.OrganizationReference

	for _, value := range input.Organizations {

		organization :=
			strings.TrimSpace(
				value.Organization,
			)

		if organization == "" {
			continue
		}

		organizations =
			append(
				organizations,
				model.OrganizationReference{
					Organization: organization,

					Source: value.Source,

					Evidence: wholeEvidenceReference(
						value.EvidenceID,
					),
				},
			)
	}

	return uniqueOrganizations(
		organizations,
	)
}

func uniqueOrganizations(
	values []model.OrganizationReference,
) []model.OrganizationReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.OrganizationReference

	for _, value := range values {

		organization :=
			strings.TrimSpace(
				value.Organization,
			)

		if organization == "" {
			continue
		}

		key :=
			strings.ToLower(
				organization +
					":" +
					value.Source,
			)

		if index,
			exists :=
			indexes[key]; exists {

			result[index].Evidence =
				mergeEvidence(
					result[index].Evidence,
					value.Evidence,
				)

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
