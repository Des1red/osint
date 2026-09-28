package web

import "strings"

func ContainsAny(
	content string,
	checks ...string,
) bool {
	for _, check := range checks {

		if strings.Contains(
			content,
			check,
		) {
			return true
		}
	}

	return false
}
