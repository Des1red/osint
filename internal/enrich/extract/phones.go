package extract

import (
	"strings"
	"unicode"

	"osint/internal/enrich/model"
)

func extractPhones(
	input model.Input,
) []model.PhoneReference {
	var phones []model.PhoneReference

	for _, value := range input.Text {

		matches :=
			phonePattern.FindAllStringIndex(
				value.Text,
				-1,
			)

		for _, match := range matches {

			if len(match) != 2 {

				continue
			}

			start :=
				match[0]

			end :=
				match[1]

			if start < 0 ||
				end < start ||
				end > len(value.Text) {

				continue
			}

			phone :=
				strings.TrimSpace(
					value.Text[start:end],
				)

			if !validPhoneCandidate(
				phone,
			) {

				continue
			}

			phones =
				append(
					phones,
					model.PhoneReference{
						Phone: phone,

						Source: value.Source,

						Evidence: evidenceReference(
							value.EvidenceID,
							start,
							end,
						),
					},
				)
		}
	}

	return uniquePhones(
		phones,
	)
}

func validPhoneCandidate(
	value string,
) bool {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {

		return false
	}

	//
	// Search snippets commonly contain dates
	// and year ranges that satisfy the loose
	// phone-number character pattern.
	//
	// Examples:
	//
	// 31.03.2026
	// 2026-03-31
	// 2003 - 2004
	//
	if dateLikePhonePattern.MatchString(
		value,
	) ||
		yearRangePhonePattern.MatchString(
			value,
		) {

		return false
	}

	digits :=
		phoneDigits(
			value,
		)

	//
	// E.164 permits at most 15 digits.
	//
	// Seven digits remains a reasonable lower
	// bound because public listings can expose
	// local numbers without a country code.
	//
	if len(digits) < 7 ||
		len(digits) > 15 {

		return false
	}

	//
	// Reject sequences which only look like
	// phone numbers because individual digits
	// are separated by spaces or punctuation.
	//
	// Example:
	//
	// 1 2 3 4 5 6 7 8 9
	//
	// Real phone numbers normally contain at
	// least one multi-digit group:
	//
	// +30 6932038409
	// +44 (0) 7983437715
	// 22870 22084
	//
	groups :=
		phoneDigitGroupPattern.FindAllString(
			value,
			-1,
		)

	if len(groups) >= 5 {

		allSingleDigit :=
			true

		for _, group := range groups {

			if len(group) > 1 {

				allSingleDigit =
					false

				break
			}
		}

		if allSingleDigit {

			return false
		}
	}

	return true
}

func uniquePhones(
	values []model.PhoneReference,
) []model.PhoneReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.PhoneReference

	for _, value := range values {

		digits :=
			phoneDigits(
				value.Phone,
			)

		if digits == "" {

			continue
		}

		//
		// Keep the existing extract-level
		// semantics:
		//
		// same phone + same source
		// = same extracted fact.
		//
		//  repeated
		// occurrences merge their provenance
		// instead of losing it.
		//
		key :=
			digits +
				":" +
				value.Source

		if index,
			exists :=
			indexes[key]; exists {

			result[index].Evidence =
				mergeEvidence(
					result[index].Evidence,
					value.Evidence,
				)

			continue
		}

		indexes[key] =
			len(result)

		result =
			append(
				result,
				value,
			)
	}

	return result
}

func phoneDigits(
	value string,
) string {
	var builder strings.Builder

	for _, character := range value {

		if unicode.IsDigit(
			character,
		) {

			builder.WriteRune(
				character,
			)
		}
	}

	return builder.String()
}
