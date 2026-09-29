package social

import (
	"net/http"
	"net/url"
	"strings"

	"osint/internal/enrich/model"
	"osint/internal/httpx"
	"osint/internal/logger"
)

type postMention struct {
	Username string

	Start int

	End int

	Kind string
}

type postDocument struct {
	URL string

	Title string

	Text string

	Mentions []postMention
}

type cachedPost struct {
	Document postDocument

	OK bool
}

func (state *builder) post(
	key string,
	postURL string,
) (
	postDocument,
	bool,
) {
	if cached,
		exists :=
		state.posts[key]; exists {

		return cached.Document,
			cached.OK
	}

	document,
		ok :=
		fetchInstagramPost(
			postURL,
		)

	state.posts[key] =
		cachedPost{
			Document: document,

			OK: ok,
		}

	return document,
		ok
}

func fetchInstagramPost(
	postURL string,
) (
	postDocument,
	bool,
) {
	response,
		err :=
		httpx.Get(
			postURL,
		)

	if err != nil {

		logger.Debug(
			"Social Post Fetch Failed",
			"url",
			postURL,
			"error",
			err.Error(),
		)

		return postDocument{},
			false
	}

	if response.StatusCode <
		http.StatusOK ||
		response.StatusCode >=
			http.StatusBadRequest {

		logger.Debug(
			"Social Post Response Rejected",
			"url",
			postURL,
			"status",
			response.StatusCode,
		)

		return postDocument{},
			false
	}

	finalURL :=
		strings.TrimSpace(
			response.FinalURL,
		)

	if finalURL == "" {

		finalURL =
			postURL
	}

	//
	// Reject redirects to login/home pages.
	//
	if !instagramPostURL(
		finalURL,
	) {

		logger.Debug(
			"Social Post Redirect Rejected",
			"requested",
			postURL,
			"final",
			finalURL,
		)

		return postDocument{},
			false
	}

	title,
		text,
		mentions,
		ok :=
		parseInstagramPost(
			response.Body,
		)

	if !ok {

		logger.Debug(
			"Social Post No Associations",
			"url",
			finalURL,
		)

		return postDocument{},
			false
	}

	logger.Debug(
		"Social Post Parsed",
		"url",
		finalURL,
		"title",
		title,
		"mentions",
		mentions,
	)

	return postDocument{
			URL: finalURL,

			Title: title,

			Text: text,

			Mentions: mentions,
		},
		true
}

func instagramPostFromLink(
	link model.ExternalLink,
) (
	string,
	bool,
) {
	values :=
		[]string{
			link.URL,
			link.ResolvedURL,
		}

	for _, value := range values {

		value =
			strings.TrimSpace(
				value,
			)

		if !instagramPostURL(
			value,
		) {

			continue
		}

		return value,
			true
	}

	return "",
		false
}

func instagramPostURL(
	value string,
) bool {
	parsed,
		err :=
		url.Parse(
			strings.TrimSpace(
				value,
			),
		)

	if err != nil {

		return false
	}

	host :=
		normalizedHost(
			parsed.Hostname(),
		)

	if host !=
		"instagram.com" {

		return false
	}

	parts :=
		strings.Split(
			strings.Trim(
				parsed.Path,
				"/",
			),
			"/",
		)

	if len(parts) <
		2 {

		return false
	}

	switch strings.ToLower(
		parts[0],
	) {

	case "p",
		"reel",
		"reels":

		return strings.TrimSpace(
			parts[1],
		) != ""
	}

	return false
}

func postKey(
	value string,
) string {
	parsed,
		err :=
		url.Parse(
			strings.TrimSpace(
				value,
			),
		)

	if err != nil {

		return ""
	}

	host :=
		normalizedHost(
			parsed.Hostname(),
		)

	path :=
		strings.TrimRight(
			strings.TrimSpace(
				parsed.Path,
			),
			"/",
		)

	if host == "" ||
		path == "" {

		return ""
	}

	return host +
		path
}

func normalizedHost(
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
