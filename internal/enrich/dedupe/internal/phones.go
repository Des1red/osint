package internal

import (
	"strings"
	"unicode"

	"osint/internal/enrich/dedupe/internal/evidence"
	"osint/internal/enrich/model"
)

func PhoneResult(
	values []model.PhoneReference,
) []model.PhoneReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.PhoneReference

	for _, value := range values {

		key :=
			phoneKey(
				value.Phone,
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

func phoneKey(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {

		return ""
	}

	value =
		strings.ReplaceAll(
			value,
			"(0)",
			"",
		)

	var builder strings.Builder

	for _, character := range value {

		if unicode.IsDigit(
			character,
		) {

			builder.WriteRune(
				character,
			)
		}
	}

	key :=
		builder.String()

	key =
		strings.TrimPrefix(
			key,
			"00",
		)

	return key
}
