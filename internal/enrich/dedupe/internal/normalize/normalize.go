package normalize

import (
	"strings"
)

func Text(
	value string,
) string {
	return strings.ToLower(
		strings.Join(
			strings.Fields(
				strings.TrimSpace(
					value,
				),
			),
			" ",
		),
	)
}
