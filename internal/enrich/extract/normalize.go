package extract

import (
	"net/url"
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

	return value
}

func normalizeURL(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	value =
		strings.Trim(
			value,
			`"'`,
		)

	value =
		strings.TrimRight(
			value,
			".,;:!?)",
		)

	if value == "" {
		return ""
	}

	if strings.HasPrefix(
		value,
		"//",
	) {
		value =
			"https:" +
				value
	}

	if strings.HasPrefix(
		value,
		"www.",
	) {
		value =
			"https://" +
				value
	}

	if !strings.HasPrefix(
		value,
		"http://",
	) &&
		!strings.HasPrefix(
			value,
			"https://",
		) {
		value =
			"https://" +
				value
	}

	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return ""
	}

	if parsed.Hostname() == "" {
		return ""
	}

	parsed.Fragment = ""

	return parsed.String()
}

func uniqueStrings(
	values []string,
) []string {
	seen :=
		make(
			map[string]struct{},
		)

	var result []string

	for _, value := range values {

		value =
			strings.TrimSpace(
				value,
			)

		if value == "" {
			continue
		}

		key :=
			strings.ToLower(
				value,
			)

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		result =
			append(
				result,
				value,
			)
	}

	return result
}
