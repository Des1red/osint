package facebook

import (
	"net/http"
	"net/url"
	"strings"

	"osint/internal/htmlx"
	"osint/internal/platforms/shared/web"
)

func UsernameLookup(
	username string,
) (FacebookResult, error) {
	requestedURL :=
		"https://www.facebook.com/" +
			url.PathEscape(username)

	page, err :=
		web.Fetch(
			requestedURL,
		)

	if err != nil {
		return FacebookResult{},
			err
	}

	finalURL :=
		page.FinalURL

	if page.StatusCode ==
		http.StatusNotFound {

		return FacebookResult{
				Found: false,

				ProfileURL: requestedURL,

				Source: finalURL,
			},
			nil
	}

	if unavailablePage(
		page.Content,
	) {
		return FacebookResult{
				Found: false,

				ProfileURL: requestedURL,

				Source: finalURL,
			},
			nil
	}

	canonicalUsername,
		canonicalURL :=
		canonicalProfile(
			finalURL,
			username,
		)

	result :=
		FacebookResult{
			Username: canonicalUsername,

			ProfileURL: canonicalURL,

			Source: finalURL,
		}

	result.Name =
		htmlx.Meta(
			page.Body,
			"og:title",
		)

	result.Description =
		htmlx.Meta(
			page.Body,
			"og:description",
		)

	if result.Name != "" ||
		result.Description != "" {

		result.Found = true
		result.Accessible = true

		return result,
			nil
	}

	if loginRequired(
		page.Content,
		finalURL,
	) {
		result.Found = false
		result.Accessible = false
		result.LoginRequired = true

		return result,
			nil
	}

	result.Found = false

	return result,
		nil
}

func canonicalProfile(
	finalURL string,
	fallbackUsername string,
) (
	string,
	string,
) {
	parsed, err :=
		url.Parse(
			finalURL,
		)

	if err != nil {
		return fallbackUsername,
			"https://www.facebook.com/" +
				url.PathEscape(
					fallbackUsername,
				)
	}

	host :=
		strings.ToLower(
			parsed.Hostname(),
		)

	if host != "facebook.com" &&
		host != "www.facebook.com" &&
		host != "m.facebook.com" {

		return fallbackUsername,
			"https://www.facebook.com/" +
				url.PathEscape(
					fallbackUsername,
				)
	}

	path :=
		strings.Trim(
			parsed.Path,
			"/",
		)

	if path == "" {
		return fallbackUsername,
			"https://www.facebook.com/" +
				url.PathEscape(
					fallbackUsername,
				)
	}

	parts :=
		strings.Split(
			path,
			"/",
		)

	candidate, err :=
		url.PathUnescape(
			parts[0],
		)

	if err != nil {
		candidate =
			parts[0]
	}

	if reservedPath(
		candidate,
	) {
		return fallbackUsername,
			"https://www.facebook.com/" +
				url.PathEscape(
					fallbackUsername,
				)
	}

	return candidate,
		"https://www.facebook.com/" +
			url.PathEscape(
				candidate,
			) +
			"/"
}

func reservedPath(
	value string,
) bool {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	switch value {

	case "",
		"login",
		"checkpoint",
		"recover",
		"help",
		"privacy",
		"settings",
		"groups",
		"pages",
		"watch",
		"marketplace",
		"gaming",
		"events",
		"messages",
		"notifications",
		"profile.php":

		return true
	}

	return false
}

func loginRequired(
	content string,
	finalURL string,
) bool {
	finalURL =
		strings.ToLower(
			finalURL,
		)

	if strings.Contains(
		finalURL,
		"/login",
	) {
		return true
	}

	return web.ContainsAny(
		content,
		"log in to facebook",
		"log into facebook",
		"please log in",
		"you must log in",
	)
}

func unavailablePage(
	content string,
) bool {
	return web.ContainsAny(
		content,
		"this content isn't available right now",
		"this content isn’t available right now",
		"this page isn't available",
		"this page isn’t available",
		"page not found",
		"the link you followed may be broken",
	)
}
