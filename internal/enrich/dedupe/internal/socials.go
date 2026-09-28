package internal

import (
	"strings"

	"osint/internal/enrich/dedupe/internal/evidence"
	"osint/internal/enrich/dedupe/internal/match"
	"osint/internal/enrich/dedupe/internal/normalize"
	"osint/internal/enrich/model"
)

func SocialResult(
	values []model.SocialReference,
) []model.SocialReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.SocialReference

	for _, value := range values {

		key :=
			strings.TrimSpace(
				value.URL,
			)

		if key == "" {

			key =
				normalize.Text(
					value.Platform +
						":" +
						value.Username,
				)
		}

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
			result[index].Platform,
		) == "" {

			result[index].Platform =
				value.Platform
		}

		if strings.TrimSpace(
			result[index].Username,
		) == "" {

			result[index].Username =
				value.Username
		}

		if strings.TrimSpace(
			result[index].URL,
		) == "" {

			result[index].URL =
				value.URL
		}

		result[index].Match =
			match.StrongerMatch(
				result[index].Match,
				value.Match,
			)
	}

	return result
}
