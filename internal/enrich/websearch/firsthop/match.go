package firsthop

import (
	"strings"
	"unicode"

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
			matchingSubjects(
				searchResultText(
					value,
				),
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

func matchingSubjects(
	value string,
	subjects []model.EvidenceSubject,
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

		firstName,
			lastName,
			ok :=
			fullNameTokens(
				subject.Value,
			)

		if !ok {

			continue
		}

		if !containsNameTokens(
			value,
			firstName,
			lastName,
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

func searchResultText(
	value discovery.SearchResult,
) string {
	return strings.TrimSpace(
		value.Title +
			" " +
			value.Snippet +
			" " +
			value.URL,
	)
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

func fullNameTokens(
	fullName string,
) (
	string,
	string,
	bool,
) {
	parts :=
		tokenize(
			fullName,
		)

	if len(parts) < 2 {

		return "",
			"",
			false
	}

	firstName :=
		parts[0]

	lastName :=
		parts[len(parts)-1]

	if firstName == "" ||
		lastName == "" {

		return "",
			"",
			false
	}

	return firstName,
		lastName,
		true
}

func containsNameTokens(
	value string,
	firstName string,
	lastName string,
) bool {
	tokens :=
		tokenize(
			value,
		)

	foundFirst :=
		false

	foundLast :=
		false

	for _, token := range tokens {

		if token ==
			firstName {

			foundFirst =
				true
		}

		if token ==
			lastName {

			foundLast =
				true
		}

		if foundFirst &&
			foundLast {

			return true
		}
	}

	return false
}

func tokenize(
	value string,
) []string {
	var builder strings.Builder

	for _, character := range strings.ToLower(
		value,
	) {

		if unicode.IsLetter(
			character,
		) ||
			unicode.IsNumber(
				character,
			) {

			builder.WriteRune(
				character,
			)

			continue
		}

		builder.WriteRune(
			' ',
		)
	}

	return strings.Fields(
		builder.String(),
	)
}
