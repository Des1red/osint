package extract

import (
	"net/url"
	"strings"
)

func classifyLink(
	value string,
) string {
	if value == "" {
		return ""
	}

	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return "External Website"
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

	case host == "linkedin.com" ||
		strings.HasSuffix(
			host,
			".linkedin.com",
		):

		return "Professional"

	case host == "facebook.com" ||
		strings.HasSuffix(
			host,
			".facebook.com",
		),
		host == "instagram.com" ||
			strings.HasSuffix(
				host,
				".instagram.com",
			),
		host == "twitter.com",
		host == "x.com",
		host == "reddit.com" ||
			strings.HasSuffix(
				host,
				".reddit.com",
			),
		host == "tiktok.com" ||
			strings.HasSuffix(
				host,
				".tiktok.com",
			),
		host == "youtube.com" ||
			strings.HasSuffix(
				host,
				".youtube.com",
			),
		host == "youtu.be":

		return "Social"
	}

	switch host {

	case "etsy.com":
		return "Marketplace"

	case "ebay.com":
		return "Marketplace"

	case "amazon.com":
		return "Marketplace"

	case "linktr.ee",
		"beacons.ai",
		"bio.link",
		"solo.to":

		return "Link Aggregator"

	case "github.com",
		"gitlab.com":

		return "Developer"

	case "medium.com",
		"substack.com":

		return "Publishing"

	case "twitch.tv":
		return "Streaming"

	case "spotify.com",
		"soundcloud.com":

		return "Music"

	default:
		return "External Website"
	}
}
