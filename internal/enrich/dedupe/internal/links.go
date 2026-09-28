package internal

import (
	"strings"

	"osint/internal/enrich/dedupe/internal/evidence"
	"osint/internal/enrich/model"
)

func LinkResult(
	values []model.ExternalLink,
) []model.ExternalLink {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.ExternalLink

	for _, value := range values {

		//
		// Deliberately use the ORIGINAL URL.
		//
		// Two different source URLs may resolve
		// to the same destination while still
		// representing distinct references.
		//
		key :=
			strings.TrimSpace(
				value.URL,
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
			result[index].ResolvedURL,
		) == "" {

			result[index].ResolvedURL =
				value.ResolvedURL
		}

		if strings.TrimSpace(
			result[index].Category,
		) == "" {

			result[index].Category =
				value.Category
		}
	}

	return result
}
