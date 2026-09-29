package brave

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const braveResultSelector = `
	div.snippet[data-type="web"],
	[data-type="web"][data-pos]
`

func parseBraveResults(
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
		braveResultSelector,
	).Each(
		func(
			_ int,
			selection *goquery.Selection,
		) {
			linkSelection :=
				findBraveResultLink(
					selection,
				)

			if linkSelection == nil ||
				linkSelection.Length() == 0 {

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
				cleanBraveURL(
					link,
				)

			if !publicBraveResultURL(
				link,
			) {
				return
			}

			title :=
				findBraveTitle(
					selection,
					linkSelection,
				)

			if title == "" {
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

						Snippet: findBraveSnippet(
							selection,
						),

						Engine: "Brave",
					},
				)
		},
	)

	return result,
		nil
}

func findBraveResultLink(
	selection *goquery.Selection,
) *goquery.Selection {
	if selection == nil ||
		selection.Length() == 0 {

		return nil
	}

	//
	// Prefer links surrounding the actual
	// result title.
	//
	selectors :=
		[]string{
			`a[href]:has(.search-snippet-title)`,
			`a[href]:has(.snippet-title)`,
			`.result-content a[href]`,
			`a[href]`,
		}

	for _, selector := range selectors {

		var found *goquery.Selection

		selection.Find(
			selector,
		).EachWithBreak(
			func(
				_ int,
				link *goquery.Selection,
			) bool {
				href,
					exists :=
					link.Attr(
						"href",
					)

				if !exists {
					return true
				}

				href =
					cleanBraveURL(
						href,
					)

				if !publicBraveResultURL(
					href,
				) {
					return true
				}

				found =
					link

				return false
			},
		)

		if found != nil &&
			found.Length() > 0 {

			return found
		}
	}

	return nil
}

func findBraveTitle(
	selection *goquery.Selection,
	linkSelection *goquery.Selection,
) string {
	selectors :=
		[]string{
			`.search-snippet-title`,
			`.snippet-title`,
			`[class*="snippet-title"]`,
			`h2`,
			`h3`,
		}

	for _, selector := range selectors {

		title :=
			selection.Find(
				selector,
			).First()

		if title.Length() == 0 {
			continue
		}

		//
		// Brave may keep the complete,
		// non-truncated title in the title
		// attribute.
		//
		if value, exists :=
			title.Attr(
				"title",
			); exists {

			value =
				strings.TrimSpace(
					value,
				)

			if value != "" {
				return value
			}
		}

		value :=
			strings.TrimSpace(
				title.Text(),
			)

		if value != "" {
			return value
		}
	}

	if linkSelection != nil &&
		linkSelection.Length() > 0 {

		value :=
			strings.TrimSpace(
				linkSelection.Text(),
			)

		if value != "" {
			return value
		}
	}

	return ""
}

func findBraveSnippet(
	selection *goquery.Selection,
) string {
	if selection == nil ||
		selection.Length() == 0 {

		return ""
	}

	selectors :=
		[]string{
			`.generic-snippet .content`,
			`.generic-snippet`,
			`.snippet-content`,
			`[data-testid="snippet-description"]`,
			`.description`,
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

func cleanBraveURL(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	//
	// Resolve scheme-relative URLs.
	//
	if strings.HasPrefix(
		value,
		"//",
	) {

		value =
			"https:" +
				value
	}

	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return ""
	}

	//
	// Relative Brave Search URLs are internal
	// navigation, not actual result URLs.
	//
	if parsed.Scheme == "" &&
		parsed.Host == "" {

		return ""
	}

	return strings.TrimSpace(
		value,
	)
}

func publicBraveResultURL(
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

	//
	// Never emit Brave Search navigation
	// itself as a discovered result.
	//
	if host == "search.brave.com" ||
		strings.HasSuffix(
			host,
			".search.brave.com",
		) {

		return false
	}

	return true
}
