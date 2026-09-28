package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"osint/internal/platforms/shared/web"
)

func UsernameLookup(
	username string,
) (RedditResult, error) {
	profileURL :=
		"https://www.reddit.com/user/" +
			url.PathEscape(username) +
			"/"

	page, err :=
		web.Fetch(
			profileURL,
		)

	if err != nil {
		return RedditResult{},
			err
	}

	finalURL :=
		page.FinalURL

	if page.StatusCode ==
		http.StatusNotFound {

		return RedditResult{
				Found: false,

				ProfileURL: profileURL,

				Source: finalURL,
			},
			nil
	}

	if page.StatusCode !=
		http.StatusOK {

		return RedditResult{},
			fmt.Errorf(
				"HTTP %d",
				page.StatusCode,
			)
	}

	if unavailableProfile(
		page.Content,
	) {
		return RedditResult{
				Found: false,

				ProfileURL: profileURL,

				Source: finalURL,
			},
			nil
	}

	result :=
		RedditResult{
			Username: username,

			ProfileURL: profileURL,

			Source: finalURL,
		}

	result.ProfileDescription =
		extractAbout(
			page.Body,
		)

	result.PostKarma =
		extractNumber(
			page.Body,
			`([0-9,]+)\s+post karma`,
		)

	result.CommentKarma =
		extractNumber(
			page.Body,
			`([0-9,]+)\s+comment karma`,
		)

	result.TotalKarma =
		result.PostKarma +
			result.CommentKarma

	result.RedditAge =
		extractString(
			page.Body,
			`([0-9]+\s+[ymd])\s+Cake day`,
		)

	result.CakeDay =
		extractString(
			page.Body,
			`Cake day:\s*([^<\n]+)`,
		)

	if !confirmedProfile(
		page.Body,
		result,
	) {
		return result,
			nil
	}

	result.Found = true
	result.Accessible = true

	return result,
		nil
}

func confirmedProfile(
	body string,
	result RedditResult,
) bool {
	if result.ProfileDescription != "" {
		return true
	}

	if result.PostKarma > 0 {
		return true
	}

	if result.CommentKarma > 0 {
		return true
	}

	if result.RedditAge != "" {
		return true
	}

	if result.CakeDay != "" {
		return true
	}

	accountPatterns := []string{
		`"id":"t2_`,
		`"accountId":"t2_`,
		`"account_id":"t2_`,
	}

	for _, pattern := range accountPatterns {

		if strings.Contains(
			body,
			pattern,
		) {
			return true
		}
	}

	return false
}

func unavailableProfile(
	content string,
) bool {
	return web.ContainsAny(
		content,
		"page not found",
		"nobody on reddit goes by that name",
		"this account has been suspended",
		"this user has deleted their account",
	)
}

func extractNumber(
	body string,
	pattern string,
) int {
	re :=
		regexp.MustCompile(
			`(?i)` +
				pattern,
		)

	match :=
		re.FindStringSubmatch(
			body,
		)

	if len(match) < 2 {
		return 0
	}

	value :=
		strings.ReplaceAll(
			match[1],
			",",
			"",
		)

	number, err :=
		strconv.Atoi(
			value,
		)

	if err != nil {
		return 0
	}

	return number
}

func extractString(
	body string,
	pattern string,
) string {
	re :=
		regexp.MustCompile(
			`(?i)` +
				pattern,
		)

	match :=
		re.FindStringSubmatch(
			body,
		)

	if len(match) < 2 {
		return ""
	}

	return strings.TrimSpace(
		match[1],
	)
}

func extractAbout(
	body string,
) string {
	re :=
		regexp.MustCompile(
			`(?is)<h2[^>]*>\s*About\s*</h2>\s*.*?<p[^>]*>(.*?)</p>`,
		)

	match :=
		re.FindStringSubmatch(
			body,
		)

	if len(match) < 2 {
		return ""
	}

	value :=
		regexp.MustCompile(
			`<[^>]+>`,
		).ReplaceAllString(
			match[1],
			"",
		)

	return strings.TrimSpace(
		value,
	)
}
