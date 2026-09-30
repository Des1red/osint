package firsthop

import (
	"strings"

	identitymatch "osint/internal/enrich/match"
	"osint/internal/enrich/model"
	"osint/internal/enrich/websearch/discovery"
)

func matchingSubjectResults(
	subjects []model.EvidenceSubject,
	values []discovery.SearchResult,
) []discovery.SearchResult {
	subjects =
		normalizeFirstHopSubjects(
			subjects,
		)

	if len(subjects) == 0 {

		return nil
	}

	seen :=
		make(
			map[string]struct{},
		)

	var result []discovery.SearchResult

	for _, value := range values {

		if len(
			matchingSearchResultSubjects(
				value,
				subjects,
			),
		) == 0 {

			continue
		}

		resultURL :=
			strings.TrimSpace(
				value.URL,
			)

		if resultURL == "" {

			continue
		}

		key :=
			strings.ToLower(
				resultURL,
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

func matchingSearchResultSubjects(
	value discovery.SearchResult,
	subjects []model.EvidenceSubject,
) []model.EvidenceSubject {
	return matchingSubjects(
		subjects,
		value.Title,
		value.Snippet,
		value.URL,
	)
}

func matchingSubjects(
	subjects []model.EvidenceSubject,
	values ...string,
) []model.EvidenceSubject {
	subjects =
		normalizeFirstHopSubjects(
			subjects,
		)

	var result []model.EvidenceSubject

	for _, subject := range subjects {

		if subject.Kind !=
			model.EvidenceAnchorPerson {

			continue
		}

		if !identitymatch.IdentityPresent(
			subject.Value,
			values...,
		) {

			continue
		}

		result =
			append(
				result,
				subject,
			)
	}

	return result
}

func normalizeFirstHopSubjects(
	values []model.EvidenceSubject,
) []model.EvidenceSubject {
	seen :=
		make(
			map[string]struct{},
		)

	var result []model.EvidenceSubject

	for _, value := range values {

		subjectValue :=
			strings.TrimSpace(
				value.Value,
			)

		if subjectValue == "" {

			continue
		}

		key :=
			string(value.Kind) +
				":" +
				strings.ToLower(
					strings.Join(
						strings.Fields(
							subjectValue,
						),
						" ",
					),
				)

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		value.Value =
			subjectValue

		result =
			append(
				result,
				value,
			)
	}

	return result
}
