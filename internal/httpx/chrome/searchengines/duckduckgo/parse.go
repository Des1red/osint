package duckduckgo

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const duckDuckGoResultSelector = `
	article[data-testid="result"],
	div[data-testid="result"],
	div.result,
	div.web-result
`

func parseDuckDuckGoResults(
	html string,
) (
	[]Result,
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

	var result []Result

	seen :=
		make(
			map[string]struct{},
		)

	document.Find(
		duckDuckGoResultSelector,
	).Each(
		func(
			_ int,
			selection *goquery.Selection,
		) {
			linkSelection :=
				findDuckDuckGoResultLink(
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
				cleanDuckDuckGoURL(
					link,
				)

			if !publicDuckDuckGoResultURL(
				link,
			) {
				return
			}

			key :=
				strings.ToLower(
					strings.TrimSpace(
						link,
					),
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
					Result{
						Title: title,

						URL: link,

						Snippet: findDuckDuckGoSnippet(
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

func findDuckDuckGoResultLink(
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

func findDuckDuckGoSnippet(
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

func cleanDuckDuckGoURL(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return ""
	}

	isRedirect :=
		strings.HasPrefix(
			value,
			"/l/?",
		)

	host :=
		strings.ToLower(
			strings.TrimSpace(
				parsed.Hostname(),
			),
		)

	if (host == "duckduckgo.com" ||
		strings.HasSuffix(
			host,
			".duckduckgo.com",
		)) &&
		parsed.Path == "/l/" {

		isRedirect =
			true
	}

	if isRedirect {

		target :=
			parsed.Query().
				Get(
					"uddg",
				)

		if target != "" {

			decoded, err :=
				url.QueryUnescape(
					target,
				)

			if err == nil {
				value =
					decoded
			}
		}
	}

	return strings.TrimSpace(
		value,
	)
}

func publicDuckDuckGoResultURL(
	value string,
) bool {
	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return false
	}

	switch strings.ToLower(
		parsed.Scheme,
	) {

	case "http",
		"https":

	default:
		return false
	}

	host :=
		strings.ToLower(
			strings.TrimSpace(
				parsed.Hostname(),
			),
		)

	if host == "" {
		return false
	}

	if host == "duckduckgo.com" ||
		strings.HasSuffix(
			host,
			".duckduckgo.com",
		) {

		return false
	}

	return true
}
