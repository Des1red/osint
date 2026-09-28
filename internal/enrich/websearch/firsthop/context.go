package firsthop

import (
	"strings"
)

const targetContextRadius = 45

type targetContextRange struct {
	Start int

	End int
}

func targetContext(
	fullName string,
	text string,
) string {
	text =
		strings.TrimSpace(
			text,
		)

	if text == "" {
		return ""
	}

	firstName,
		lastName,
		ok :=
		fullNameTokens(
			fullName,
		)

	if !ok {
		return ""
	}

	fields :=
		strings.Fields(
			text,
		)

	if len(fields) == 0 {
		return ""
	}

	type indexedToken struct {
		Value string

		Field int
	}

	var tokens []indexedToken

	for fieldIndex, field := range fields {

		for _, token := range tokenize(
			field,
		) {

			tokens =
				append(
					tokens,
					indexedToken{
						Value: token,

						Field: fieldIndex,
					},
				)
		}
	}

	if len(tokens) < 2 {
		return ""
	}

	var ranges []targetContextRange

	for index := 0; index < len(tokens)-1; index++ {

		left :=
			tokens[index]

		right :=
			tokens[index+1]

		forward :=
			left.Value ==
				firstName &&
				right.Value ==
					lastName

		reverse :=
			left.Value ==
				lastName &&
				right.Value ==
					firstName

		if !forward &&
			!reverse {

			continue
		}

		start :=
			left.Field -
				targetContextRadius

		if start < 0 {

			start =
				0
		}

		end :=
			right.Field +
				targetContextRadius +
				1

		if end >
			len(fields) {

			end =
				len(fields)
		}

		if len(ranges) > 0 {

			last :=
				&ranges[len(ranges)-1]

			if start <=
				last.End {

				if end >
					last.End {

					last.End =
						end
				}

				continue
			}
		}

		ranges =
			append(
				ranges,
				targetContextRange{
					Start: start,

					End: end,
				},
			)
	}

	if len(ranges) == 0 {
		return ""
	}

	var parts []string

	for _, value := range ranges {

		parts =
			append(
				parts,
				strings.Join(
					fields[value.Start:value.End],
					" ",
				),
			)
	}

	return strings.TrimSpace(
		strings.Join(
			parts,
			" ",
		),
	)
}
