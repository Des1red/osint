package match

import "osint/internal/matcher"

func StrongerMatch(
	left matcher.Level,
	right matcher.Level,
) matcher.Level {
	if matchStrength(
		right,
	) >
		matchStrength(
			left,
		) {

		return right
	}

	return left
}

func matchStrength(
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
