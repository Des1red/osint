package firsthop

import (
	"bytes"
	"net/url"
	"osint/internal/enrich/model"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type pageEvidence struct {
	Subject model.EvidenceSubject

	RequestedURL string

	URL string

	Redirected bool

	CrossSite bool

	Text string

	Links []string

	Organizations []string

	Locations []string

	Employment []employmentEvidence

	Education []educationEvidence
}

func parsePage(
	finalURL string,
	body []byte,
	fullName string,
) (
	pageEvidence,
	error,
) {
	document, err :=
		goquery.NewDocumentFromReader(
			bytes.NewReader(
				body,
			),
		)

	if err != nil {
		return pageEvidence{},
			err
	}

	title :=
		strings.TrimSpace(
			document.Find(
				"title",
			).First().Text(),
		)

	description :=
		pageDescription(
			document,
		)

	links,
		references :=
		pageLinks(
			document,
			finalURL,
		)

	//
	// Read JSON-LD before removing scripts.
	//
	structured :=
		structuredEvidence(
			document,
			fullName,
			finalURL,
		)

	links =
		append(
			links,
			structured.URLs...,
		)

	references =
		append(
			references,
			structured.References...,
		)

	links =
		uniquePageValues(
			links,
		)

	references =
		uniquePageValues(
			references,
		)

	document.Find(
		"script, style, noscript, template, svg",
	).Remove()

	bodyText :=
		strings.Join(
			strings.Fields(
				document.Find(
					"body",
				).Text(),
			),
			" ",
		)

	textParts :=
		[]string{
			title,
			description,
			bodyText,
		}

	//
	// mailto, tel, JSON-LD email and JSON-LD
	// telephone values become ordinary text
	// evidence for the shared extractor.
	//
	textParts =
		append(
			textParts,
			references...,
		)

	return pageEvidence{
			URL: finalURL,

			Text: strings.TrimSpace(
				strings.Join(
					textParts,
					" ",
				),
			),

			Links: links,

			Organizations: structured.Organizations,

			Locations: structured.Locations,

			Employment: structured.Employment,

			Education: structured.Education,
		},
		nil
}

func pageDescription(
	document *goquery.Document,
) string {
	if document == nil {
		return ""
	}

	selectors :=
		[]string{
			`meta[name="description"]`,
			`meta[property="og:description"]`,
			`meta[name="twitter:description"]`,
		}

	for _, selector := range selectors {

		value,
			exists :=
			document.Find(
				selector,
			).First().Attr(
				"content",
			)

		if !exists {
			continue
		}

		value =
			strings.TrimSpace(
				value,
			)

		if value != "" {
			return value
		}
	}

	return ""
}

func pageLinks(
	document *goquery.Document,
	baseURL string,
) (
	[]string,
	[]string,
) {
	if document == nil {
		return nil,
			nil
	}

	base,
		err :=
		url.Parse(
			baseURL,
		)

	if err != nil {
		return nil,
			nil
	}

	baseHost :=
		strings.ToLower(
			base.Hostname(),
		)

	seenLinks :=
		make(
			map[string]struct{},
		)

	seenReferences :=
		make(
			map[string]struct{},
		)

	var links []string
	var references []string

	addReference :=
		func(
			value string,
		) {
			value =
				strings.TrimSpace(
					value,
				)

			if value == "" {
				return
			}

			key :=
				strings.ToLower(
					value,
				)

			if _, exists :=
				seenReferences[key]; exists {

				return
			}

			seenReferences[key] =
				struct{}{}

			references =
				append(
					references,
					value,
				)
		}

	document.Find(
		"a[href]",
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

			if href == "" {
				return
			}

			lower :=
				strings.ToLower(
					href,
				)

			if strings.HasPrefix(
				lower,
				"mailto:",
			) {

				value :=
					strings.TrimSpace(
						href[len("mailto:"):],
					)

				if decoded,
					err :=
					url.QueryUnescape(
						value,
					); err == nil {

					value =
						decoded
				}

				if index :=
					strings.Index(
						value,
						"?",
					); index >= 0 {

					value =
						value[:index]
				}

				addReference(
					value,
				)

				return
			}

			if strings.HasPrefix(
				lower,
				"tel:",
			) {

				value :=
					strings.TrimSpace(
						href[len("tel:"):],
					)

				if decoded,
					err :=
					url.QueryUnescape(
						value,
					); err == nil {

					value =
						decoded
				}

				addReference(
					value,
				)

				return
			}

			parsed,
				err :=
				url.Parse(
					href,
				)

			if err != nil {
				return
			}

			resolved :=
				base.ResolveReference(
					parsed,
				)

			scheme :=
				strings.ToLower(
					resolved.Scheme,
				)

			if scheme != "http" &&
				scheme != "https" {

				return
			}

			host :=
				strings.ToLower(
					resolved.Hostname(),
				)

			if host == "" {
				return
			}

			//
			// Still one hop only.
			//
			// Internal navigation is ignored.
			// External URLs become evidence but
			// are never fetched by FirstHop.
			//
			if host == baseHost {
				return
			}

			resolved.Fragment =
				""

			value :=
				resolved.String()

			key :=
				strings.ToLower(
					value,
				)

			if _, exists :=
				seenLinks[key]; exists {

				return
			}

			seenLinks[key] =
				struct{}{}

			links =
				append(
					links,
					value,
				)
		},
	)

	return links,
		references
}

func uniquePageValues(
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
