package twitter

import (
	"net/http"
	"net/url"
	"strings"

	"osint/internal/htmlx"
	"osint/internal/platforms/shared/web"
)

func UsernameLookup(
	username string,
) (TwitterResult, error) {
	profileURL :=
		"https://x.com/" +
			url.PathEscape(username)

	page, err :=
		web.Fetch(
			profileURL,
		)

	if err != nil {
		return TwitterResult{},
			err
	}

	if page.StatusCode ==
		http.StatusNotFound {

		return TwitterResult{
				Found: false,

				ProfileURL: profileURL,

				Source: page.FinalURL,
			},
			nil
	}

	finalURL :=
		page.FinalURL

	if loginRequired(
		page.Content,
		finalURL,
	) {
		return TwitterResult{
				Found: true,

				Accessible: false,

				LoginRequired: true,

				Username: username,

				ProfileURL: profileURL,

				Source: finalURL,
			},
			nil
	}

	if unavailableProfile(
		page.Content,
	) {
		return TwitterResult{
				Found: false,

				ProfileURL: profileURL,

				Source: finalURL,
			},
			nil
	}

	result :=
		TwitterResult{
			Found: true,

			Accessible: true,

			Username: username,

			ProfileURL: profileURL,

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

	return result,
		nil
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
		"log in to x",
		"sign in to x",
		"log in to twitter",
		"sign in to twitter",
	)
}

func unavailableProfile(
	content string,
) bool {
	return web.ContainsAny(
		content,
		"this account doesn’t exist",
		"this account doesn't exist",
		"account suspended",
		"something went wrong",
	)
}
