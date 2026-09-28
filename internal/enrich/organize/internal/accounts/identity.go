package accounts

import (
	"strings"
	"unicode"
)

func identityPresent(
	value string,
	identity string,
) bool {
	valueTokens :=
		textTokens(
			value,
		)

	identityTokens :=
		textTokens(
			identity,
		)

	if len(valueTokens) == 0 ||
		len(identityTokens) == 0 {

		return false
	}

	//
	// When a person has middle names or
	// initials, use first and last components
	// as the stable identity anchors.
	//
	required :=
		identityTokens

	if len(identityTokens) > 2 {

		required =
			[]string{
				identityTokens[0],

				identityTokens[len(identityTokens)-1],
			}
	}

	available :=
		make(
			map[string]struct{},
			len(valueTokens),
		)

	for _, token := range valueTokens {

		available[token] =
			struct{}{}
	}

	for _, token := range required {

		if _, exists :=
			available[token]; !exists {

			return false
		}
	}

	return true
}

func textTokens(
	value string,
) []string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	if value == "" {

		return nil
	}

	return strings.FieldsFunc(
		value,
		func(
			character rune,
		) bool {

			return !unicode.IsLetter(
				character,
			) &&
				!unicode.IsDigit(
					character,
				)
		},
	)
}

func urlSearchText(
	value string,
) string {
	replacer :=
		strings.NewReplacer(
			"/",
			" ",

			"-",
			" ",

			"_",
			" ",

			".",
			" ",

			"%20",
			" ",
		)

	return strings.ToLower(
		replacer.Replace(
			value,
		),
	)
}
