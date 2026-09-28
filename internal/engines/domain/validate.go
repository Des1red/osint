package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var domainLabelPattern = regexp.MustCompile(
	`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`,
)

func normalizeDomain(
	value string,
) string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	value =
		strings.TrimSuffix(
			value,
			".",
		)

	return value
}

func validate(
	domain string,
) error {
	if domain == "" {

		return fmt.Errorf(
			"domain is empty",
		)
	}

	if len(domain) > 253 {

		return fmt.Errorf(
			"domain exceeds maximum length",
		)
	}

	if strings.Contains(
		domain,
		"://",
	) {

		return fmt.Errorf(
			"domain must not contain a URL scheme",
		)
	}

	if strings.ContainsAny(
		domain,
		"/ \t\r\n",
	) {

		return fmt.Errorf(
			"invalid domain format",
		)
	}

	labels :=
		strings.Split(
			domain,
			".",
		)

	if len(labels) < 2 {

		return fmt.Errorf(
			"domain must contain a hostname and top-level domain",
		)
	}

	for _, label := range labels {

		if label == "" {

			return fmt.Errorf(
				"domain contains an empty label",
			)
		}

		if len(label) > 63 {

			return fmt.Errorf(
				"domain label exceeds maximum length",
			)
		}

		if !domainLabelPattern.MatchString(
			label,
		) {

			return fmt.Errorf(
				"invalid domain label: %s",
				label,
			)
		}
	}

	return nil
}
