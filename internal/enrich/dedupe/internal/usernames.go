package internal

import (
	"osint/internal/enrich/dedupe/internal/evidence"
	"osint/internal/enrich/dedupe/internal/match"
	"osint/internal/enrich/dedupe/internal/normalize"
	"osint/internal/enrich/model"
)

func UsernameResult(
	values []model.UsernameReference,
) []model.UsernameReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.UsernameReference

	for _, value := range values {

		key :=
			normalize.Text(
				value.Username,
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

		result[index].Match =
			match.StrongerMatch(
				result[index].Match,
				value.Match,
			)
	}

	return result
}
