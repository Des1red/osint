package instagram

import (
	"net/http"
	"net/url"
	"strings"

	"osint/internal/htmlx"
	"osint/internal/platforms/shared/web"
)

func UsernameLookup(
	username string,
) (InstagramResult, error) {
	profileURL :=
		"https://www.instagram.com/" +
			url.PathEscape(username) +
			"/"

	page, err :=
		web.Fetch(
			profileURL,
		)

	if err != nil {
		return InstagramResult{},
			err
	}

	finalURL :=
		page.FinalURL

	if page.StatusCode ==
		http.StatusNotFound {

		return InstagramResult{
				Found: false,

				ProfileURL: profileURL,

				Source: finalURL,
			},
			nil
	}

	if unavailableProfile(
		page.Content,
	) {
		return InstagramResult{
				Found: false,

				ProfileURL: profileURL,

				Source: finalURL,
			},
			nil
	}

	result :=
		InstagramResult{
			Username: username,

			ProfileURL: profileURL,

			Source: finalURL,
		}

	if loginRequired(
		finalURL,
	) {
		result.Found = true
		result.Accessible = false
		result.LoginRequired = true

		return result,
			nil
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

	if profileMetadata(
		result.Name,
		result.Description,
		username,
	) {
		result.Found = true
		result.Accessible = true

		return result,
			nil
	}

	result.Found = false
	result.Accessible = false

	return result,
		nil
}

func loginRequired(
	finalURL string,
) bool {
	finalURL =
		strings.ToLower(
			finalURL,
		)

	return strings.Contains(
		finalURL,
		"/accounts/login",
	)
}

func unavailableProfile(
	content string,
) bool {
	return web.ContainsAny(
		content,
		"sorry, this page isn't available",
		"the link you followed may be broken",
		"page not found",
	)
}

func profileMetadata(
	name string,
	description string,
	username string,
) bool {
	name =
		strings.TrimSpace(
			name,
		)

	description =
		strings.TrimSpace(
			description,
		)

	username =
		strings.ToLower(
			strings.TrimSpace(
				username,
			),
		)

	if strings.EqualFold(
		name,
		"Instagram",
	) &&
		description == "" {

		return false
	}

	handle :=
		"@" +
			username

	if username != "" {
		if strings.Contains(
			strings.ToLower(
				name,
			),
			handle,
		) {
			return true
		}

		if strings.Contains(
			strings.ToLower(
				description,
			),
			handle,
		) {
			return true
		}
	}

	return false
}
