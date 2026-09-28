package extract

import (
	"strings"

	"osint/internal/enrich/model"
	"osint/internal/matcher"
)

var usernameMonthSuffixes = []string{
	"september",
	"february",
	"november",
	"december",
	"january",
	"october",
	"august",
	"march",
	"april",
	"june",
	"july",
	"sept",
	"jan",
	"feb",
	"mar",
	"apr",
	"jun",
	"jul",
	"aug",
	"sep",
	"oct",
	"nov",
	"dec",
	"may",
}

func extractUsernames(
	input model.Input,
) []model.UsernameReference {
	var usernames []model.UsernameReference

	//
	// Explicit username evidence supplied by
	// another collector.
	//
	for _, value := range input.Usernames {

		username :=
			normalizeUsername(
				value.Username,
			)

		if !validUsername(
			username,
		) {

			continue
		}

		match :=
			matcher.Classify(
				input.RootUsername,
				username,
			)

		//
		// Plain repetition of the root username
		// is not new information.
		//
		if match ==
			matcher.Exact {

			continue
		}

		usernames =
			append(
				usernames,
				model.UsernameReference{
					Username: username,

					Source: value.Source,

					Match: match,

					Evidence: wholeEvidenceReference(
						value.EvidenceID,
					),
				},
			)
	}

	//
	// Username occurrences discovered inside
	// free text.
	//
	for _, value := range input.Text {

		//
		// Preserve the original text exactly.
		//
		// Email ranges are skipped rather than
		// removed so EvidenceReference offsets
		// continue pointing into Evidence.Text.
		//
		emailRanges :=
			emailPattern.FindAllStringIndex(
				value.Text,
				-1,
			)

		matches :=
			usernamePattern.FindAllStringIndex(
				value.Text,
				-1,
			)

		for _, matchIndex := range matches {

			if len(matchIndex) != 2 {

				continue
			}

			start :=
				matchIndex[0]

			end :=
				matchIndex[1]

			if start < 0 ||
				end < start ||
				end > len(value.Text) {

				continue
			}

			if overlapsAnyRange(
				start,
				end,
				emailRanges,
			) {

				continue
			}

			discovered :=
				value.Text[start:end]

			trailing :=
				value.Text[end:]

			username :=
				cleanExtractedUsername(
					discovered,
					trailing,
				)

			if !validUsername(
				username,
			) {

				continue
			}

			//
			// The extracted token may have been
			// shortened because rendered website
			// UI glued metadata to the handle.
			//
			// Example:
			//
			//     @Person2 1Dec
			//
			// becomes:
			//
			//     @Person2 1
			//
			// Evidence should end at the actual
			// username, not at the removed UI
			// metadata.
			//
			evidenceEnd :=
				start +
					1 +
					len(username)

			if evidenceEnd >
				end {

				evidenceEnd =
					end
			}

			match :=
				matcher.Classify(
					input.RootUsername,
					username,
				)

			if match ==
				matcher.Exact {

				continue
			}

			usernames =
				append(
					usernames,
					model.UsernameReference{
						Username: username,

						Source: value.Source,

						Match: match,

						Evidence: evidenceReference(
							value.EvidenceID,
							start,
							evidenceEnd,
						),
					},
				)
		}
	}

	return uniqueUsernames(
		usernames,
	)
}

func cleanExtractedUsername(
	value string,
	trailing string,
) string {
	username :=
		normalizeUsername(
			value,
		)

	if username == "" {

		return ""
	}

	//
	// Rendered pages sometimes collapse adjacent
	// DOM elements:
	//
	//     @Person2 1 Joined December 2015
	//
	// can become:
	//
	//     @Person2 1Joined December 2015
	//
	// or:
	//
	//     @Person2 1Dec 2, 2015
	//
	// The regexp correctly sees one continuous
	// token, so clean those known UI suffixes
	// afterwards using the following text as
	// context.
	//

	//
	// Sentence punctuation can also be consumed
	// because "." is a valid username character.
	//
	if strings.HasSuffix(
		username,
		".",
	) &&
		usernameTrailingBoundary(
			trailing,
		) {

		username =
			strings.TrimSuffix(
				username,
				".",
			)
	}

	lower :=
		strings.ToLower(
			username,
		)

	monthSuffix :=
		usernameAttachedMonthSuffix(
			lower,
		)

	if monthSuffix != "" &&
		usernameTrailingStartsDigit(
			trailing,
		) {

		username =
			username[:len(username)-
				len(monthSuffix)]

		lower =
			strings.ToLower(
				username,
			)
	}

	//
	// Handles such as:
	//
	//     Person2 1JoinedDec
	//
	// first lose "Dec", leaving:
	//
	//     Person2 1Joined
	//
	if strings.HasSuffix(
		lower,
		"joined",
	) &&
		(usernameTrailingStartsMonth(
			trailing,
		) ||
			usernameTrailingStartsDigit(
				trailing,
			)) {

		username =
			username[:len(username)-
				len("joined")]
	}

	return normalizeUsername(
		username,
	)
}

func usernameAttachedMonthSuffix(
	value string,
) string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	for _, suffix := range usernameMonthSuffixes {

		if len(value) <=
			len(suffix) {

			continue
		}

		if strings.HasSuffix(
			value,
			suffix,
		) {

			return suffix
		}
	}

	return ""
}

func usernameTrailingStartsDigit(
	value string,
) bool {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {

		return false
	}

	return value[0] >= '0' &&
		value[0] <= '9'
}

func usernameTrailingStartsMonth(
	value string,
) bool {
	fields :=
		strings.Fields(
			strings.TrimSpace(
				value,
			),
		)

	if len(fields) == 0 {

		return false
	}

	word :=
		strings.ToLower(
			strings.Trim(
				fields[0],
				`"'“”‘’.,;:()[]{}-`,
			),
		)

	switch word {

	case "january",
		"february",
		"march",
		"april",
		"may",
		"june",
		"july",
		"august",
		"september",
		"october",
		"november",
		"december",
		"jan",
		"feb",
		"mar",
		"apr",
		"jun",
		"jul",
		"aug",
		"sep",
		"sept",
		"oct",
		"nov",
		"dec":

		return true
	}

	return false
}

func usernameTrailingBoundary(
	value string,
) bool {
	if value == "" {

		return true
	}

	switch value[0] {

	case ' ',
		'\t',
		'\n',
		'\r',
		',',
		';',
		':',
		'!',
		'?':

		return true
	}

	return false
}

func overlapsAnyRange(
	start int,
	end int,
	ranges [][]int,
) bool {
	for _, value := range ranges {

		if len(value) != 2 {

			continue
		}

		rangeStart :=
			value[0]

		rangeEnd :=
			value[1]

		if start <
			rangeEnd &&
			end >
				rangeStart {

			return true
		}
	}

	return false
}

func validUsername(
	value string,
) bool {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {

		return false
	}

	hasAlphaNumeric :=
		false

	for _, char := range value {

		switch {

		case char >= 'a' &&
			char <= 'z':

			hasAlphaNumeric =
				true

		case char >= 'A' &&
			char <= 'Z':

			hasAlphaNumeric =
				true

		case char >= '0' &&
			char <= '9':

			hasAlphaNumeric =
				true

		case char == '.',
			char == '_',
			char == '-':

		default:

			return false
		}
	}

	return hasAlphaNumeric
}

func genericUsernameCandidate(
	value string,
) bool {
	value =
		normalizeUsername(
			value,
		)

	if !validUsername(
		value,
	) {

		return false
	}

	//
	// Numeric profile IDs are valid identifiers
	// on some platforms, especially Facebook.
	//
	// Keep them as SocialReference usernames,
	// but do not duplicate them into the generic
	// Usernames section.
	//
	for _, char := range value {

		if char >= 'a' &&
			char <= 'z' {

			return true
		}

		if char >= 'A' &&
			char <= 'Z' {

			return true
		}
	}

	return false
}

func uniqueUsernames(
	values []model.UsernameReference,
) []model.UsernameReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.UsernameReference

	for _, value := range values {

		username :=
			strings.TrimSpace(
				value.Username,
			)

		if username == "" {

			continue
		}

		key :=
			strings.ToLower(
				username +
					":" +
					value.Source,
			)

		if index,
			exists :=
			indexes[key]; exists {

			result[index].Evidence =
				mergeEvidence(
					result[index].Evidence,
					value.Evidence,
				)

			result[index].Match =
				strongerUsernameMatch(
					result[index].Match,
					value.Match,
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

func strongerUsernameMatch(
	left matcher.Level,
	right matcher.Level,
) matcher.Level {
	if usernameMatchStrength(
		right,
	) >
		usernameMatchStrength(
			left,
		) {

		return right
	}

	return left
}

func usernameMatchStrength(
	value matcher.Level,
) int {
	switch value {

	case matcher.Exact:

		return 4

	case matcher.Close:

		return 3

	case matcher.Broad:

		return 2

	case matcher.None:

		return 1

	default:

		return 0
	}
}
