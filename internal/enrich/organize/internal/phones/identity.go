package phones

import (
	"regexp"
	"strings"

	"osint/internal/enrich/model"
)

type identityCandidate struct {
	Owner owner

	Name string

	Variants []string
}

func identityCandidates(
	input model.Input,
	people []model.PersonReference,
) []identityCandidate {
	var result []identityCandidate

	target :=
		strings.TrimSpace(
			input.FullName,
		)

	if target != "" {

		result =
			append(
				result,
				identityCandidate{
					Owner: owner{
						Kind: ownerTarget,
					},

					Name: target,

					Variants: nameVariants(
						target,
					),
				},
			)
	}

	for index, person := range people {

		name :=
			strings.TrimSpace(
				person.Name,
			)

		if name == "" {

			continue
		}

		if target != "" &&
			sameIdentityName(
				target,
				name,
			) {

			continue
		}

		result =
			append(
				result,
				identityCandidate{
					Owner: owner{
						Kind: ownerPerson,

						PersonIndex: index,
					},

					Name: name,

					Variants: nameVariants(
						name,
					),
				},
			)
	}

	return result
}

func nameVariants(
	value string,
) []string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {

		return nil
	}

	seen :=
		make(
			map[string]struct{},
		)

	var result []string

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

			key :=
				strings.ToLower(
					value,
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
					value,
				)
		}

	add(
		value,
	)

	fields :=
		strings.Fields(
			value,
		)

	if len(fields) >= 2 {

		first :=
			fields[0]

		last :=
			fields[len(fields)-1]

		add(
			first +
				" " +
				last,
		)

		add(
			last +
				" " +
				first,
		)
	}

	return result
}

func sameIdentityName(
	left string,
	right string,
) bool {
	leftFields :=
		strings.Fields(
			strings.ToLower(
				strings.TrimSpace(
					left,
				),
			),
		)

	rightFields :=
		strings.Fields(
			strings.ToLower(
				strings.TrimSpace(
					right,
				),
			),
		)

	if len(leftFields) == 0 ||
		len(rightFields) == 0 {

		return false
	}

	if strings.Join(
		leftFields,
		" ",
	) ==
		strings.Join(
			rightFields,
			" ",
		) {

		return true
	}

	if len(leftFields) < 2 ||
		len(rightFields) < 2 {

		return false
	}

	return leftFields[0] ==
		rightFields[len(rightFields)-1] &&
		leftFields[len(leftFields)-1] ==
			rightFields[0]
}

func identityScore(
	text string,
	phoneStart int,
	candidate identityCandidate,
) int {
	best :=
		0

	for _, variant := range candidate.Variants {

		pattern :=
			namePattern(
				variant,
			)

		if pattern == nil {

			continue
		}

		matches :=
			pattern.FindAllStringIndex(
				text,
				-1,
			)

		for _, match := range matches {

			if len(match) != 2 {

				continue
			}

			score :=
				scoreOccurrence(
					text,
					match[0],
					match[1],
					phoneStart,
				)

			if score >
				best {

				best =
					score
			}
		}
	}

	return best
}

func namePattern(
	value string,
) *regexp.Regexp {
	fields :=
		strings.Fields(
			strings.TrimSpace(
				value,
			),
		)

	if len(fields) == 0 {

		return nil
	}

	parts :=
		make(
			[]string,
			0,
			len(fields),
		)

	for _, field := range fields {

		parts =
			append(
				parts,
				regexp.QuoteMeta(
					field,
				),
			)
	}

	pattern :=
		`(?i)` +
			strings.Join(
				parts,
				`\s+`,
			)

	compiled,
		err :=
		regexp.Compile(
			pattern,
		)

	if err != nil {

		return nil
	}

	return compiled
}

func scoreOccurrence(
	text string,
	nameStart int,
	nameEnd int,
	phoneStart int,
) int {
	//
	// We deliberately require the identity to
	// precede the phone.
	//
	// This prevents:
	//
	// Phone: +30 company-number
	// Person1 +30 personal-number
	//
	// from assigning the company number to Person1
	// just because his name appears later.
	//
	if nameEnd >
		phoneStart {

		return 0
	}

	gap :=
		phoneStart -
			nameEnd

	//
	// Very close identity + phone:
	//
	// Person1 +30...
	// Person2: +44...
	//
	if gap <= 32 {

		return 1000 -
			gap
	}

	beforeStart :=
		nameStart -
			64

	if beforeStart < 0 {

		beforeStart =
			0
	}

	before :=
		strings.ToLower(
			text[beforeStart:nameStart],
		)

	between :=
		strings.ToLower(
			text[nameEnd:phoneStart],
		)

	//
	// Named contact blocks:
	//
	// Contact: Person1
	// Address: ...
	// Phone: ...
	//
	if gap <= 260 &&
		containsContactMarker(
			before,
		) {

		return 850 -
			gap
	}

	//
	// Slightly looser but still explicit:
	//
	// Person1
	// Mobile: +30...
	//
	if gap <= 120 &&
		containsPhoneMarker(
			between,
		) {

		return 750 -
			gap
	}

	return 0
}

func containsContactMarker(
	value string,
) bool {
	return strings.Contains(
		value,
		"contact",
	)
}

func containsPhoneMarker(
	value string,
) bool {
	markers :=
		[]string{
			"phone",
			"mobile",
			"telephone",
			"tel:",
			"tel ",
			"whatsapp",
		}

	for _, marker := range markers {

		if strings.Contains(
			value,
			marker,
		) {

			return true
		}
	}

	return false
}
