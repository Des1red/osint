package accounts

import (
	"net/url"
	"strings"
)

func platformFromURL(
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
		strings.ToLower(
			strings.TrimSpace(
				parsed.Hostname(),
			),
		)

	host =
		strings.TrimPrefix(
			host,
			"www.",
		)

	host =
		strings.TrimPrefix(
			host,
			"m.",
		)

	switch {

	case host ==
		"instagram.com" ||
		strings.HasSuffix(
			host,
			".instagram.com",
		):

		return "instagram"

	case host ==
		"facebook.com" ||
		strings.HasSuffix(
			host,
			".facebook.com",
		):

		return "facebook"

	case host ==
		"x.com" ||
		host ==
			"twitter.com":

		return "twitter"

	case host ==
		"linkedin.com" ||
		strings.HasSuffix(
			host,
			".linkedin.com",
		):

		return "linkedin"

	case host ==
		"reddit.com" ||
		strings.HasSuffix(
			host,
			".reddit.com",
		):

		return "reddit"

	case host ==
		"tiktok.com" ||
		strings.HasSuffix(
			host,
			".tiktok.com",
		):

		return "tiktok"

	case host ==
		"youtube.com" ||
		strings.HasSuffix(
			host,
			".youtube.com",
		):

		return "youtube"

	case host ==
		"github.com":

		return "github"

	case host ==
		"gitlab.com":

		return "gitlab"
	}

	return ""
}
