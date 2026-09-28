package discovery

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const resultSelector = `
	article[data-testid="result"],
	div[data-testid="result"],
	div.result,
	div.web-result
`

func parseResults(
	html string,
) (
	[]SearchResult,
	error,
) {
	document, err :=
		goquery.NewDocumentFromReader(
			strings.NewReader(
				html,
			),
		)

	if err != nil {
		return nil,
			err
	}

	var result []SearchResult

	seen :=
		make(
			map[string]struct{},
		)

	document.Find(
		resultSelector,
	).Each(
		func(
			_ int,
			selection *goquery.Selection,
		) {
			linkSelection :=
				findResultLink(
					selection,
				)

			if linkSelection == nil ||
				linkSelection.Length() == 0 {

				return
			}

			title :=
				strings.TrimSpace(
					linkSelection.Text(),
				)

			if title == "" {
				return
			}

			link,
				exists :=
				linkSelection.Attr(
					"href",
				)

			if !exists {
				return
			}

			link =
				cleanURL(
					link,
				)

			if link == "" {
				return
			}

			if !publicHTTPURL(
				link,
			) {
				return
			}

			key :=
				canonicalURL(
					link,
				)

			if key == "" {
				return
			}

			if _, exists :=
				seen[key]; exists {

				return
			}

			seen[key] =
				struct{}{}

			result =
				append(
					result,
					SearchResult{
						Title: title,

						URL: link,

						Snippet: findSnippet(
							selection,
						),

						Engine: "DuckDuckGo",
					},
				)
		},
	)

	return result,
		nil
}

func findResultLink(
	selection *goquery.Selection,
) *goquery.Selection {
	if selection == nil ||
		selection.Length() == 0 {

		return nil
	}

	selectors :=
		[]string{
			`a[data-testid="result-title-a"]`,
			`h2 a[href]`,
			`a.result__a[href]`,
			`a[data-testid="result-extras-url-link"]`,
		}

	for _, selector := range selectors {

		link :=
			selection.Find(
				selector,
			).First()

		if link.Length() == 0 {
			continue
		}

		href,
			exists :=
			link.Attr(
				"href",
			)

		if !exists ||
			strings.TrimSpace(
				href,
			) == "" {

			continue
		}

		return link
	}

	return nil
}

func findSnippet(
	selection *goquery.Selection,
) string {
	if selection == nil ||
		selection.Length() == 0 {

		return ""
	}

	selectors :=
		[]string{
			`[data-result="snippet"]`,
			`[data-testid="result-snippet"]`,
			`.result__snippet`,
			`.snippet`,
			".OgdwYG6KE2qthn9XQWFC",
		}

	for _, selector := range selectors {

		value :=
			strings.TrimSpace(
				selection.Find(
					selector,
				).First().Text(),
			)

		if value != "" {
			return value
		}
	}

	return ""
}
