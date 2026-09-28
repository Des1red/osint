package fullname

import (
	"strings"
	"unicode"
)

func prepareCandidates(
	fullName string,
) []string {
	parts :=
		strings.Fields(
			strings.TrimSpace(
				fullName,
			),
		)

	if len(parts) < 2 {
		return nil
	}

	firstName :=
		strings.TrimSpace(
			parts[0],
		)

	lastName :=
		strings.TrimSpace(
			parts[len(parts)-1],
		)

	if firstName == "" ||
		lastName == "" {

		return nil
	}

	separators :=
		[]string{
			"",
			".",
			"_",
			"-",
		}

	var result []string

	seen :=
		make(
			map[string]struct{},
		)

	add :=
		func(
			value string,
		) {
			value =
				strings.TrimSpace(
					value,
				)

			if value == "" {
				return
			}

			if _, exists :=
				seen[value]; exists {

				return
			}

			seen[value] =
				struct{}{}

			result =
				append(
					result,
					value,
				)
		}

	addVariants :=
		func(
			left string,
			right string,
		) {
			left =
				strings.TrimSpace(
					left,
				)

			right =
				strings.TrimSpace(
					right,
				)

			for _, separator := range separators {

				candidate :=
					left +
						separator +
						right

				//
				// all lowercase
				//
				add(
					strings.ToLower(
						candidate,
					),
				)

				//
				// ALL CAPITAL
				//
				add(
					strings.ToUpper(
						candidate,
					),
				)

				//
				// Only first letter capital.
				//
				add(
					capitalizeFirst(
						candidate,
					),
				)
			}
		}

	//
	// First name + last name.
	//
	addVariants(
		firstName,
		lastName,
	)

	//
	// Last name + first name.
	//
	addVariants(
		lastName,
		firstName,
	)

	return result
}

func capitalizeFirst(
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

	runes :=
		[]rune(
			value,
		)

	runes[0] =
		unicode.ToUpper(
			runes[0],
		)

	return string(
		runes,
	)
}
