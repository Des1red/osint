package internal

import (
	"osint/internal/enrich/dedupe/internal/evidence"
	"osint/internal/enrich/dedupe/internal/normalize"
	"osint/internal/enrich/model"
)

func LocationResult(
	values []model.LocationReference,
) []model.LocationReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.LocationReference

	for _, value := range values {

		key :=
			normalize.Text(
				value.Location,
			)

		if key == "" {

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
	}

	return result
}
