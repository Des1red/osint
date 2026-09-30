package match

import (
	"strings"
	"unicode"
)

// IdentityPresent checks whether a person's
// identity occurs as one coherent name inside
// any supplied value.
//
// First and last name must occur together.
func IdentityPresent(
	identity string,
	values ...string,
) bool {
	for _, value := range values {

		valueTokens :=
			Tokens(
				value,
			)

		if len(valueTokens) < 2 {

			continue
		}

		identityTokens :=
			Tokens(
				identity,
			)

		if len(identityTokens) < 2 {

			return false
		}

		//
		// Prefer the complete identity.
		//
		if containsSequence(
			valueTokens,
			identityTokens,
		) {

			return true
		}

		//
		// First + last are stable identity
		// anchors and may also appear reversed.
		//
		for index := 0; index < len(valueTokens)-1; index++ {

			if IdentityPair(
				identity,
				valueTokens[index],
				valueTokens[index+1],
			) {

				return true
			}
		}
	}

	return false
}

func IdentityPair(
	identity string,
	left string,
	right string,
) bool {
	firstName,
		lastName,
		ok :=
		IdentityAnchors(
			identity,
		)

	if !ok {

		return false
	}

	forward :=
		left ==
			firstName &&
			right ==
				lastName

	reverse :=
		left ==
			lastName &&
			right ==
				firstName

	return forward ||
		reverse
}

func IdentityAnchors(
	identity string,
) (
	string,
	string,
	bool,
) {
	tokens :=
		Tokens(
			identity,
		)

	if len(tokens) < 2 {

		return "",
			"",
			false
	}

	firstName :=
		tokens[0]

	lastName :=
		tokens[len(tokens)-1]

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

func Tokens(
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

func containsSequence(
	value []string,
	sequence []string,
) bool {
	if len(sequence) == 0 ||
		len(value) <
			len(sequence) {

		return false
	}

	for start := 0; start <=
		len(value)-len(sequence); start++ {

		matched :=
			true

		for index := range sequence {

			if value[start+index] !=
				sequence[index] {

				matched =
					false

				break
			}
		}

		if matched {

			return true
		}
	}

	return false
}
