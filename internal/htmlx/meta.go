package htmlx

import (
	"html"
	"strings"
)

func Meta(
	body string,
	property string,
) string {
	body =
		strings.TrimSpace(
			body,
		)

	property =
		strings.TrimSpace(
			property,
		)

	if body == "" ||
		property == "" {

		return ""
	}

	patterns := []string{
		`property="` +
			property +
			`" content="`,

		`property='` +
			property +
			`' content='`,

		`name="` +
			property +
			`" content="`,

		`name='` +
			property +
			`' content='`,
	}

	for _, pattern := range patterns {

		start :=
			strings.Index(
				body,
				pattern,
			)

		if start == -1 {
			continue
		}

		start +=
			len(pattern)

		quote :=
			byte('"')

		if strings.Contains(
			pattern,
			"'",
		) {
			quote = '\''
		}

		end :=
			strings.IndexByte(
				body[start:],
				quote,
			)

		if end == -1 {
			continue
		}

		value :=
			strings.TrimSpace(
				body[start : start+end],
			)

		if value == "" {
			continue
		}

		return html.UnescapeString(
			value,
		)
	}

	return ""
}
