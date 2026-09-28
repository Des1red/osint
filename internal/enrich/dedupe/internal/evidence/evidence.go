package evidence

import (
	"strings"

	"osint/internal/enrich/model"
)

type referenceKey struct {
	EvidenceID string

	Start int

	End int
}

func Merge(
	groups ...[]model.EvidenceReference,
) []model.EvidenceReference {
	seen :=
		make(
			map[referenceKey]struct{},
		)

	var result []model.EvidenceReference

	for _, group := range groups {

		for _, value := range group {

			evidenceID :=
				strings.TrimSpace(
					value.EvidenceID,
				)

			if evidenceID == "" {

				continue
			}

			key :=
				referenceKey{
					EvidenceID: evidenceID,

					Start: value.Start,

					End: value.End,
				}

			if _, exists :=
				seen[key]; exists {

				continue
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
	}

	return result
}

func PreferSource(
	current string,
	incoming string,
	currentEvidence []model.EvidenceReference,
	incomingEvidence []model.EvidenceReference,
) string {
	current =
		strings.TrimSpace(
			current,
		)

	incoming =
		strings.TrimSpace(
			incoming,
		)

	if incoming == "" {

		return current
	}

	if current == "" {

		return incoming
	}

	//
	// If the existing duplicate had no
	// provenance but this occurrence does,
	// prefer the source associated with the
	// evidence-bearing occurrence.
	//
	if len(currentEvidence) == 0 &&
		len(incomingEvidence) > 0 {

		return incoming
	}

	return current
}
