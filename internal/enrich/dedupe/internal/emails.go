package internal

import (
	"strings"

	"osint/internal/enrich/dedupe/internal/evidence"
	"osint/internal/enrich/dedupe/internal/normalize"
	"osint/internal/enrich/model"
)

func EmailResult(
	values []model.EmailReference,
) []model.EmailReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.EmailReference

	for _, value := range values {

		key :=
			normalize.Text(
				value.Email,
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

		if strings.TrimSpace(
			result[index].Type,
		) == "" &&
			strings.TrimSpace(
				value.Type,
			) != "" {

			result[index].Type =
				value.Type
		}
	}

	return result
}
