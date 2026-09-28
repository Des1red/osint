package tiktok

import (
	"net/http"
	"net/url"
	"strings"

	"osint/internal/htmlx"
	"osint/internal/platforms/shared/web"
)

func UsernameLookup(
	username string,
) (TikTokResult, error) {
	profileURL :=
		"https://www.tiktok.com/@" +
			url.PathEscape(username)

	page, err :=
		web.Fetch(
			profileURL,
		)

	if err != nil {
		return TikTokResult{},
			err
	}

	if page.StatusCode ==
		http.StatusNotFound {

		return TikTokResult{
				Found: false,

				ProfileURL: profileURL,

				Source: page.FinalURL,
			},
			nil
	}

	finalURL :=
		page.FinalURL

	result :=
		TikTokResult{
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

	if result.Name != "" ||
		result.Description != "" {

		result.Found = true
		result.Accessible = true

		return result,
			nil
	}

	if unavailableProfile(
		page.Content,
	) {
		result.Found = false

		return result,
			nil
	}

	if loginRequired(
		page.Content,
		finalURL,
	) {
		result.Found = true
		result.Accessible = false
		result.LoginRequired = true

		return result,
			nil
	}

	result.Found = true
	result.Accessible = false

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
		"log in to tiktok",
		"sign in to tiktok",
		"login to tiktok",
	)
}

func unavailableProfile(
	content string,
) bool {
	return web.ContainsAny(
		content,
		"couldn't find this account",
		"could not find this account",
		"account not found",
		"page not available",
	)
}
