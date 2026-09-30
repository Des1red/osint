package gr11888

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var mapContentPattern = regexp.MustCompile(
	`content:\s*"([^"]*)"`,
)

var mapLatitudePattern = regexp.MustCompile(
	`lat:\s*"([^"]*)"`,
)

var mapLongitudePattern = regexp.MustCompile(
	`lng:\s*"([^"]*)"`,
)

func parseEntry(
	body string,
	sourceURL string,
) (
	Entry,
	error,
) {
	document, err :=
		goquery.NewDocumentFromReader(
			strings.NewReader(
				body,
			),
		)

	if err != nil {

		return Entry{},
			err
	}

	result :=
		Entry{
			SourceURL: strings.TrimSpace(
				sourceURL,
			),
		}

	//
	// Name.
	//
	result.Title =
		strings.TrimSpace(
			document.Find(
				"h1",
			).
				First().
				Text(),
		)

	//
	// Address and coordinates.
	//
	mapButton :=
		document.Find(
			`button[aria-label="Κουμπί_Διεύθυνσης"]`,
		).
			First()

	if mapButton.Length() > 0 {

		onclick,
			exists :=
			mapButton.Attr(
				"onclick",
			)

		if exists {

			appendPatternField(
				&result,
				"Address",
				onclick,
				mapContentPattern,
			)

			appendPatternField(
				&result,
				"Latitude",
				onclick,
				mapLatitudePattern,
			)

			appendPatternField(
				&result,
				"Longitude",
				onclick,
				mapLongitudePattern,
			)
		}
	}

	//
	// Telephone numbers.
	//
	document.Find(
		`a[href^="tel:"]`,
	).Each(
		func(
			_ int,
			selection *goquery.Selection,
		) {
			href,
				exists :=
				selection.Attr(
					"href",
				)

			if !exists {

				return
			}

			phone :=
				strings.TrimSpace(
					strings.TrimPrefix(
						href,
						"tel:",
					),
				)

			if phone == "" {

				return
			}

			result.Fields =
				append(
					result.Fields,
					Field{
						Name: "Phone",

						Value: phone,
					},
				)
		},
	)

	return result,
		nil
}

func appendPatternField(
	entry *Entry,
	name string,
	value string,
	pattern *regexp.Regexp,
) {
	match :=
		pattern.FindStringSubmatch(
			value,
		)

	if len(match) < 2 {

		return
	}

	fieldValue :=
		strings.TrimSpace(
			match[1],
		)

	if fieldValue == "" {

		return
	}

	entry.Fields =
		append(
			entry.Fields,
			Field{
				Name: name,

				Value: fieldValue,
			},
		)
}

var recordPathPattern = regexp.MustCompile(
	`^/white-pages/\d+/$`,
)

func parseRecordURLs(
	body string,
) (
	[]string,
	error,
) {
	document, err :=
		goquery.NewDocumentFromReader(
			strings.NewReader(
				body,
			),
		)

	if err != nil {

		return nil,
			err
	}

	seen :=
		make(
			map[string]struct{},
		)

	var result []string

	document.Find(
		`a[href]`,
	).Each(
		func(
			_ int,
			selection *goquery.Selection,
		) {
			href,
				exists :=
				selection.Attr(
					"href",
				)

			if !exists {

				return
			}

			href =
				strings.TrimSpace(
					href,
				)

			if !recordPathPattern.MatchString(
				href,
			) {

				return
			}

			recordURL :=
				"https://www.11888.gr" +
					href

			if _, exists :=
				seen[recordURL]; exists {

				return
			}

			seen[recordURL] =
				struct{}{}

			result =
				append(
					result,
					recordURL,
				)
		},
	)

	return result,
		nil
}
