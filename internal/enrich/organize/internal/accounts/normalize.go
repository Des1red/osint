package accounts

import (
	"strings"
)

func normalizeUsername(
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

func normalizePlatform(
	value string,
) string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	switch value {

	case "x":
		return "twitter"
	}

	return value
}
