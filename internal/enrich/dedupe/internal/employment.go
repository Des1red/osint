package internal

import (
	"strings"

	"osint/internal/enrich/dedupe/internal/evidence"
	"osint/internal/enrich/dedupe/internal/normalize"
	"osint/internal/enrich/model"
)

func EmploymentResult(
	values []model.EmploymentReference,
) []model.EmploymentReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.EmploymentReference

	for _, value := range values {

		key :=
			normalize.Text(
				value.Title,
			) +
				":" +
				normalize.Text(
					value.Organization,
				) +
				":" +
				normalize.Text(
					value.StartDate,
				) +
				":" +
				normalize.Text(
					value.EndDate,
				)

		if key == ":::" {

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
	}

	return result
}
