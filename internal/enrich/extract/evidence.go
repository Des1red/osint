package extract

import (
	"strconv"
	"strings"

	"osint/internal/enrich/model"
)

func evidenceReference(
	evidenceID string,
	start int,
	end int,
) []model.EvidenceReference {
	evidenceID =
		strings.TrimSpace(
			evidenceID,
		)

	if evidenceID == "" {

		return nil
	}

	if start < 0 {

		start =
			0
	}

	if end < start {

		end =
			start
	}

	return []model.EvidenceReference{
		{
			EvidenceID: evidenceID,

			Start: start,

			End: end,
		},
	}
}

func wholeEvidenceReference(
	evidenceID string,
) []model.EvidenceReference {
	return evidenceReference(
		evidenceID,
		0,
		0,
	)
}

func mergeEvidence(
	left []model.EvidenceReference,
	right []model.EvidenceReference,
) []model.EvidenceReference {
	seen :=
		make(
			map[string]struct{},
		)

	result :=
		make(
			[]model.EvidenceReference,
			0,
			len(left)+len(right),
		)

	add :=
		func(
			value model.EvidenceReference,
		) {
			evidenceID :=
				strings.TrimSpace(
					value.EvidenceID,
				)

			if evidenceID == "" {

				return
			}

			key :=
				evidenceID +
					":" +
					strconv.Itoa(
						value.Start,
					) +
					":" +
					strconv.Itoa(
						value.End,
					)

			if _, exists :=
				seen[key]; exists {

				return
			}

			seen[key] =
				struct{}{}

			result =
				append(
					result,
					model.EvidenceReference{
						EvidenceID: evidenceID,

						Start: value.Start,

						End: value.End,
					},
				)
		}

	for _, value := range left {

		add(
			value,
		)
	}

	for _, value := range right {

		add(
			value,
		)
	}

	return result
}
