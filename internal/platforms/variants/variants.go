package variants

import (
	"strings"
	"unicode"

	"osint/internal/matcher"
)

const MaxCandidates = 12

type Candidate struct {
	Username string

	Match matcher.Level

	Source string
}

func Generate(
	root string,
) []Candidate {
	root =
		normalize(
			root,
		)

	if root == "" {
		return nil
	}

	var result []Candidate

	seen :=
		make(
			map[string]struct{},
		)

	add :=
		func(
			username string,
			source string,
		) {
			if len(result) >=
				MaxCandidates {

				return
			}

			username =
				normalize(
					username,
				)

			if username == "" {
				return
			}

			if strings.EqualFold(
				username,
				root,
			) {
				return
			}

			if !valid(
				username,
			) {
				return
			}

			key :=
				strings.ToLower(
					username,
				)

			if _, exists :=
				seen[key]; exists {

				return
			}

			match :=
				matcher.Classify(
					root,
					username,
				)

			if match ==
				matcher.None {

				return
			}

			seen[key] =
				struct{}{}

			result =
				append(
					result,
					Candidate{
						Username: username,

						Match: match,

						Source: source,
					},
				)
		}

	compact :=
		removeSeparators(
			root,
		)

	//
	// Root already contains separators.
	//
	if compact != root {

		add(
			compact,
			"separator removed",
		)

		for _, separator := range []rune{
			'_',
			'.',
			'-',
		} {

			add(
				replaceSeparators(
					root,
					separator,
				),
				"separator variation",
			)
		}
	}

	//
	// Insert separators at likely username
	// boundaries.
	//
	base :=
		compact

	boundaries :=
		candidateBoundaries(
			base,
		)

	for _, boundary := range boundaries {

		for _, separator := range []rune{
			'_',
			'.',
			'-',
		} {

			add(
				insertSeparator(
					base,
					boundary,
					separator,
				),
				"separator variation",
			)
		}
	}

	//
	// Leading/trailing separator variants.
	//
	add(
		base+"_",
		"separator variation",
	)

	add(
		"_"+base,
		"separator variation",
	)

	//
	// Conservative broader variants.
	//
	add(
		base+"1",
		"numeric suffix",
	)

	add(
		base+"01",
		"numeric suffix",
	)

	add(
		"real"+base,
		"common prefix",
	)

	add(
		"its"+base,
		"common prefix",
	)

	return result
}

func candidateBoundaries(
	value string,
) []int {
	if len(value) < 2 {
		return nil
	}

	var result []int

	seen :=
		make(
			map[int]struct{},
		)

	add :=
		func(
			position int,
		) {
			if position <= 0 ||
				position >= len(value) {

				return
			}

			if _, exists :=
				seen[position]; exists {

				return
			}

			seen[position] =
				struct{}{}

			result =
				append(
					result,
					position,
				)
		}

	//
	// Prefer transitions such as:
	//
	// j0|kelo
	// des1|red
	//
	for i := 1; i < len(value); i++ {

		left :=
			rune(
				value[i-1],
			)

		right :=
			rune(
				value[i],
			)

		leftDigit :=
			unicode.IsDigit(
				left,
			)

		rightDigit :=
			unicode.IsDigit(
				right,
			)

		if leftDigit !=
			rightDigit {

			add(
				i,
			)
		}
	}

	//
	// Fall back to the middle of the username.
	//
	if len(result) == 0 {

		add(
			len(value) / 2,
		)
	}

	//
	// Keep generation bounded.
	//
	if len(result) > 2 {
		result =
			result[:2]
	}

	return result
}

func insertSeparator(
	value string,
	position int,
	separator rune,
) string {
	if position <= 0 ||
		position >= len(value) {

		return value
	}

	return value[:position] +
		string(separator) +
		value[position:]
}

func removeSeparators(
	value string,
) string {
	var builder strings.Builder

	for _, char := range value {

		switch char {

		case '.',
			'_',
			'-':

			continue

		default:
			builder.WriteRune(
				char,
			)
		}
	}

	return builder.String()
}

func replaceSeparators(
	value string,
	separator rune,
) string {
	var builder strings.Builder

	for _, char := range value {

		switch char {

		case '.',
			'_',
			'-':

			builder.WriteRune(
				separator,
			)

		default:
			builder.WriteRune(
				char,
			)
		}
	}

	return builder.String()
}

func normalize(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	value =
		strings.TrimPrefix(
			value,
			"@",
		)

	return value
}

func valid(
	value string,
) bool {
	if value == "" {
		return false
	}

	for _, char := range value {

		switch {

		case char >= 'a' &&
			char <= 'z':

		case char >= 'A' &&
			char <= 'Z':

		case char >= '0' &&
			char <= '9':

		case char == '.',
			char == '_',
			char == '-':

		default:
			return false
		}
	}

	return true
}
