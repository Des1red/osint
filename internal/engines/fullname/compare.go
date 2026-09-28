package fullname

import (
	"strings"
	"unicode"
)

func matchCandidates(
	fullName string,
	candidates []Candidate,
) []Match {
	firstName,
		lastName :=
		nameParts(
			fullName,
		)

	if firstName == "" ||
		lastName == "" {

		return nil
	}

	var matches []Match

	seen :=
		make(
			map[string]struct{},
		)

	for _, candidate := range candidates {
		if !nameMatches(
			candidate.Name,
			firstName,
			lastName,
		) {
			continue
		}

		match :=
			Match{
				Platform: candidate.Platform,

				Username: candidate.Username,

				Name: candidate.Name,

				ProfileURL: candidate.ProfileURL,
			}

		key :=
			candidateKey(
				match.Platform,
				match.Username,
				match.ProfileURL,
			)

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		matches =
			append(
				matches,
				match,
			)
	}

	return matches
}

func nameParts(
	fullName string,
) (
	string,
	string,
) {
	parts :=
		strings.Fields(
			normalizeName(
				fullName,
			),
		)

	if len(parts) < 2 {
		return "",
			""
	}

	firstName :=
		normalizeNamePart(
			parts[0],
		)

	lastName :=
		normalizeNamePart(
			parts[len(parts)-1],
		)

	return firstName,
		lastName
}

func nameMatches(
	value string,
	firstName string,
	lastName string,
) bool {
	parts :=
		strings.Fields(
			normalizeName(
				value,
			),
		)

	if len(parts) == 0 {
		return false
	}

	hasFirstName := false
	hasLastName := false

	for _, part := range parts {
		part =
			normalizeNamePart(
				part,
			)

		if part == firstName {
			hasFirstName = true
		}

		if part == lastName {
			hasLastName = true
		}
	}

	return hasFirstName &&
		hasLastName
}

func normalizeName(
	value string,
) string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	if value == "" {
		return ""
	}

	var builder strings.Builder

	for _, char := range value {

		switch {

		case unicode.IsLetter(
			char,
		):

			builder.WriteRune(
				char,
			)

		case unicode.IsDigit(
			char,
		):

			builder.WriteRune(
				char,
			)

		default:
			builder.WriteRune(
				' ',
			)
		}
	}

	return strings.Join(
		strings.Fields(
			builder.String(),
		),
		" ",
	)
}

func normalizeNamePart(
	value string,
) string {
	value =
		normalizeName(
			value,
		)

	return strings.ReplaceAll(
		value,
		" ",
		"",
	)
}

func candidateKey(
	platform string,
	username string,
	profileURL string,
) string {
	platform =
		strings.ToLower(
			strings.TrimSpace(
				platform,
			),
		)

	profileURL =
		strings.ToLower(
			strings.TrimSpace(
				profileURL,
			),
		)

	if profileURL != "" {
		return platform +
			"|" +
			profileURL
	}

	return platform +
		"|" +
		strings.ToLower(
			strings.TrimSpace(
				username,
			),
		)
}
