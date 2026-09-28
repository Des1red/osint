package firsthop

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"osint/internal/enrich/model"
	"osint/internal/enrich/websearch/discovery"
	"osint/internal/httpx"
	"osint/internal/logger"
)

const maxFirstHopConcurrency = 5

const firstHopDebugPreviewLength = 300

func inspect(
	searchAnchor string,
	subjects []model.EvidenceSubject,
	values []discovery.SearchResult,
) []pageEvidence {
	subjects =
		normalizeFirstHopSubjects(
			subjects,
		)

	if len(values) == 0 ||
		len(subjects) == 0 {

		return nil
	}

	pages :=
		make(
			[][]pageEvidence,
			len(values),
		)

	semaphore :=
		make(
			chan struct{},
			maxFirstHopConcurrency,
		)

	var waitGroup sync.WaitGroup

	for index, value := range values {

		waitGroup.Add(
			1,
		)

		go func(
			index int,
			value discovery.SearchResult,
		) {
			defer waitGroup.Done()

			semaphore <- struct{}{}

			defer func() {

				<-semaphore
			}()

			requestedURL :=
				strings.TrimSpace(
					value.URL,
				)

			logger.Debug(
				"FirstHop Fetch",
				"anchor",
				searchAnchor,
				"url",
				requestedURL,
				"engine",
				value.Engine,
				"query",
				value.Query,
			)

			response,
				err :=
				httpx.Get(
					requestedURL,
				)

			if err != nil {

				logger.Debug(
					"FirstHop Fetch Failed",
					"anchor",
					searchAnchor,
					"url",
					requestedURL,
					"error",
					err.Error(),
				)

				return
			}

			finalURL :=
				strings.TrimSpace(
					response.FinalURL,
				)

			if finalURL == "" {

				finalURL =
					requestedURL
			}

			siteChanged :=
				firstHopSiteChanged(
					requestedURL,
					finalURL,
				)

			logger.Debug(
				"FirstHop Response",
				"anchor",
				searchAnchor,
				"requested url",
				requestedURL,
				"final url",
				finalURL,
				"status",
				response.StatusCode,
				"redirected",
				requestedURL !=
					finalURL,
				"site changed",
				siteChanged,
			)

			if response.StatusCode <
				http.StatusOK ||
				response.StatusCode >=
					http.StatusBadRequest {

				logger.Debug(
					"FirstHop Response Rejected",
					"anchor",
					searchAnchor,
					"url",
					requestedURL,
					"final url",
					finalURL,
					"status",
					response.StatusCode,
				)

				return
			}

			var matchedPages []pageEvidence
			var matchedSubjects []string

			//
			// Fetch once, then test the destination
			// independently against every known
			// person in the current investigation
			// graph.
			//
			// The search anchor does not decide who
			// the page is about.
			//
			for _, subject := range subjects {

				if subject.Kind !=
					model.EvidenceAnchorPerson {

					continue
				}

				subjectName :=
					strings.TrimSpace(
						subject.Value,
					)

				if subjectName == "" {

					continue
				}

				page,
					err :=
					parsePage(
						finalURL,
						response.Body,
						subjectName,
					)

				if err != nil {

					logger.Debug(
						"FirstHop Parse Failed",
						"anchor",
						searchAnchor,
						"subject",
						subjectName,
						"url",
						requestedURL,
						"final url",
						finalURL,
						"error",
						err.Error(),
					)

					continue
				}

				rawText :=
					page.Text

				contextText :=
					targetContext(
						subjectName,
						rawText,
					)

				identityMatch :=
					contextText != ""

				logger.Debug(
					"FirstHop Page Evidence",
					"anchor",
					searchAnchor,
					"subject",
					subjectName,
					"requested url",
					requestedURL,
					"page url",
					page.URL,
					"identity match",
					identityMatch,
					"text length",
					len(rawText),
					"text preview",
					firstHopTextPreview(
						rawText,
					),
					"subject context",
					firstHopTextPreview(
						contextText,
					),
					"links",
					page.Links,
					"organizations",
					page.Organizations,
					"locations",
					page.Locations,
					"employment",
					page.Employment,
					"education",
					page.Education,
				)

				if !identityMatch {

					continue
				}

				page.Subject =
					subject

				page.RequestedURL =
					requestedURL

				page.Redirected =
					requestedURL !=
						finalURL

				page.CrossSite =
					siteChanged

				//
				// The evidence text is intentionally
				// subject-local.
				//
				// This lets one fetched page produce
				// a Person1 evidence record and a Person2
				// evidence record without pretending
				// the search anchor owns either one.
				//
				page.Text =
					contextText

				matchedPages =
					append(
						matchedPages,
						page,
					)

				matchedSubjects =
					append(
						matchedSubjects,
						subjectName,
					)
			}

			if len(matchedPages) == 0 {

				reason :=
					"no known subject identity present in page"

				if siteChanged {

					reason =
						"cross-site redirect destination does not contain any known subject"
				}

				logger.Debug(
					"FirstHop Page Skipped",
					"anchor",
					searchAnchor,
					"requested url",
					requestedURL,
					"final url",
					finalURL,
					"reason",
					reason,
				)

				return
			}

			if siteChanged {

				logger.Debug(
					"FirstHop Cross-Site Redirect Accepted",
					"anchor",
					searchAnchor,
					"subjects",
					matchedSubjects,
					"requested url",
					requestedURL,
					"final url",
					finalURL,
					"reason",
					"destination independently contains a known subject",
				)
			}

			pages[index] =
				matchedPages
		}(
			index,
			value,
		)
	}

	waitGroup.Wait()

	var result []pageEvidence

	for _, group := range pages {

		result =
			append(
				result,
				group...,
			)
	}

	return result
}

func firstHopSiteChanged(
	requestedURL string,
	finalURL string,
) bool {
	requested,
		err :=
		url.Parse(
			requestedURL,
		)

	if err != nil {

		return false
	}

	final,
		err :=
		url.Parse(
			finalURL,
		)

	if err != nil {

		return false
	}

	requestedHost :=
		normalizeFirstHopHost(
			requested.Hostname(),
		)

	finalHost :=
		normalizeFirstHopHost(
			final.Hostname(),
		)

	if requestedHost == "" ||
		finalHost == "" {

		return false
	}

	if requestedHost ==
		finalHost {

		return false
	}

	//
	// Treat ordinary subdomain changes as the
	// same site.
	//
	// Examples:
	//
	// gr.linkedin.com -> linkedin.com
	// m.facebook.com  -> facebook.com
	//
	if strings.HasSuffix(
		requestedHost,
		"."+finalHost,
	) ||
		strings.HasSuffix(
			finalHost,
			"."+requestedHost,
		) {

		return false
	}

	return true
}

func normalizeFirstHopHost(
	value string,
) string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	value =
		strings.TrimPrefix(
			value,
			"www.",
		)

	value =
		strings.TrimPrefix(
			value,
			"m.",
		)

	return value
}

func firstHopTextPreview(
	value string,
) string {
	value =
		strings.Join(
			strings.Fields(
				strings.TrimSpace(
					value,
				),
			),
			" ",
		)

	if value == "" {

		return ""
	}

	runes :=
		[]rune(
			value,
		)

	if len(runes) <=
		firstHopDebugPreviewLength {

		return value
	}

	return string(
		runes[:firstHopDebugPreviewLength],
	) +
		"..."
}
