package social

import (
	"strconv"
	"strings"

	"osint/internal/enrich/model"
)

func mergeEvidenceReferences(
	left []model.EvidenceReference,
	right []model.EvidenceReference,
) []model.EvidenceReference {
	result :=
		append(
			[]model.EvidenceReference(nil),
			left...,
		)

	seen :=
		make(
			map[string]struct{},
		)

	for _, value := range result {

		seen[evidenceReferenceKey(
			value,
		)] =
			struct{}{}
	}

	for _, value := range right {

		key :=
			evidenceReferenceKey(
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

	return result
}

func evidenceReferenceKey(
	value model.EvidenceReference,
) string {
	return strings.TrimSpace(
		value.EvidenceID,
	) +
		":" +
		strconv.Itoa(
			value.Start,
		) +
		":" +
		strconv.Itoa(
			value.End,
		)
}

func appendUniqueString(
	values []string,
	value string,
) []string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {

		return values
	}

	for _, existing := range values {

		if strings.EqualFold(
			strings.TrimSpace(
				existing,
			),
			value,
		) {

			return values
		}
	}

	return append(
		values,
		value,
	)
}
