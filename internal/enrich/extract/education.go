package extract

import (
	"strings"

	"osint/internal/enrich/model"
)

func extractEducation(
	input model.Input,
) []model.EducationReference {
	var result []model.EducationReference

	for _, value := range input.Education {

		school :=
			strings.TrimSpace(
				value.School,
			)

		if school == "" {
			continue
		}

		result =
			append(
				result,
				model.EducationReference{
					School: school,

					Degrees: value.Degrees,

					Majors: value.Majors,

					StartDate: strings.TrimSpace(
						value.StartDate,
					),

					EndDate: strings.TrimSpace(
						value.EndDate,
					),

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

	return uniqueEducation(
		result,
	)
}

func uniqueEducation(
	values []model.EducationReference,
) []model.EducationReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.EducationReference

	for _, value := range values {

		key :=
			strings.ToLower(
				strings.TrimSpace(
					value.School,
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

			if len(result[index].Degrees) == 0 &&
				len(value.Degrees) > 0 {

				result[index].Degrees =
					value.Degrees
			}

			if len(result[index].Majors) == 0 &&
				len(value.Majors) > 0 {

				result[index].Majors =
					value.Majors
			}

			if strings.TrimSpace(
				result[index].Summary,
			) == "" {

				result[index].Summary =
					value.Summary
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
