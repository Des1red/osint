package matcher

import (
	"strings"
)

type Level string

const (
	Exact Level = "exact"
	Close Level = "close"
	Broad Level = "broad"
	None  Level = "none"
)

func Classify(
	root string,
	candidate string,
) Level {
	root =
		normalize(
			root,
		)

	candidate =
		normalize(
			candidate,
		)

	if root == "" ||
		candidate == "" {

		return None
	}

	if strings.EqualFold(
		root,
		candidate,
	) {
		return Exact
	}

	rootCompact :=
		compact(
			root,
		)

	candidateCompact :=
		compact(
			candidate,
		)

	if rootCompact == "" ||
		candidateCompact == "" {

		return None
	}

	//
	// Same username after removing separators.
	//
	// j0kelo
	// j0_kelo
	// j0.kelo
	// j0-kelo
	//
	if rootCompact ==
		candidateCompact {

		return Close
	}

	//
	// Candidate contains the full root username
	// with only a small amount of extra text.
	//
	if len(rootCompact) >= 4 &&
		strings.Contains(
			candidateCompact,
			rootCompact,
		) {

		difference :=
			len(candidateCompact) -
				len(rootCompact)

		if difference > 0 &&
			difference <= 4 {

			return Broad
		}
	}

	//
	// Candidate may be a slightly shortened
	// version of the root.
	//
	if len(candidateCompact) >= 4 &&
		strings.Contains(
			rootCompact,
			candidateCompact,
		) {

		difference :=
			len(rootCompact) -
				len(candidateCompact)

		if difference > 0 &&
			difference <= 2 {

			return Broad
		}
	}

	//
	// Small spelling variations.
	//
	distance :=
		editDistance(
			rootCompact,
			candidateCompact,
		)

	shorter :=
		len(rootCompact)

	if len(candidateCompact) <
		shorter {

		shorter =
			len(candidateCompact)
	}

	if shorter >= 5 &&
		distance == 1 {

		return Broad
	}

	if shorter >= 7 &&
		distance <= 2 {

		return Broad
	}

	return None
}

func Label(
	level Level,
) string {
	switch level {

	case Exact:
		return "Exact"

	case Close:
		return "Close"

	case Broad:
		return "Broad"

	default:
		return "None"
	}
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

	return strings.ToLower(
		value,
	)
}

func compact(
	value string,
) string {
	value =
		normalize(
			value,
		)

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

func editDistance(
	left string,
	right string,
) int {
	if left == right {
		return 0
	}

	if left == "" {
		return len(right)
	}

	if right == "" {
		return len(left)
	}

	previous :=
		make(
			[]int,
			len(right)+1,
		)

	current :=
		make(
			[]int,
			len(right)+1,
		)

	for index := range previous {

		previous[index] =
			index
	}

	for i := 1; i <= len(left); i++ {

		current[0] =
			i

		for j := 1; j <= len(right); j++ {

			cost := 0

			if left[i-1] !=
				right[j-1] {

				cost = 1
			}

			deletion :=
				previous[j] +
					1

			insertion :=
				current[j-1] +
					1

			substitution :=
				previous[j-1] +
					cost

			current[j] =
				minimum(
					deletion,
					insertion,
					substitution,
				)
		}

		previous,
			current =
			current,
			previous
	}

	return previous[len(right)]
}

func minimum(
	values ...int,
) int {
	minimum :=
		values[0]

	for _, value := range values[1:] {

		if value <
			minimum {

			minimum =
				value
		}
	}

	return minimum
}
