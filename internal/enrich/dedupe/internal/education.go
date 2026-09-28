package internal

import (
	"strings"

	"osint/internal/enrich/dedupe/internal/evidence"
	"osint/internal/enrich/dedupe/internal/normalize"
	"osint/internal/enrich/model"
)

func EducationResult(
	values []model.EducationReference,
) []model.EducationReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.EducationReference

	for _, value := range values {

		key :=
			normalize.Text(
				value.School,
			) +
				":" +
				normalize.Text(
					value.StartDate,
				) +
				":" +
				normalize.Text(
					value.EndDate,
				)

		if key == "::" {

			continue
		}

		index,
			exists :=
			indexes[key]

		if !exists {

			indexes[key] =
				len(result)

			result =
				append(
					result,
					value,
				)

			continue
		}

		source :=
			evidence.PreferSource(
				result[index].Source,
				value.Source,
				result[index].Evidence,
				value.Evidence,
			)

		result[index].Evidence =
			evidence.Merge(
				result[index].Evidence,
				value.Evidence,
			)

		result[index].Source =
			source

		result[index].Degrees =
			mergeStringValues(
				result[index].Degrees,
				value.Degrees,
			)

		result[index].Majors =
			mergeStringValues(
				result[index].Majors,
				value.Majors,
			)

		if strings.TrimSpace(
			result[index].Summary,
		) == "" {

			result[index].Summary =
				value.Summary
		}
	}

	return result
}

func mergeStringValues(
	groups ...[]string,
) []string {
	seen :=
		make(
			map[string]struct{},
		)

	var result []string

	for _, group := range groups {

		for _, value := range group {

			value =
				strings.TrimSpace(
					value,
				)

			if value == "" {

				continue
			}

			key :=
				normalize.Text(
					value,
				)

			if _, exists :=
				seen[key]; exists {

				continue
			}

			seen[key] =
				struct{}{}

			result =
				append(
					result,
					value,
				)
		}
	}

	return result
}
