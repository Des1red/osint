package discovery

import (
	"net"
	"net/url"
	"strings"
)

func cleanURL(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	//
	// DuckDuckGo redirect URL.
	//
	if strings.HasPrefix(
		value,
		"/l/?",
	) {

		parsed, err :=
			url.Parse(
				value,
			)

		if err != nil {
			return ""
		}

		target :=
			parsed.Query().
				Get(
					"uddg",
				)

		if target != "" {

			decoded, err :=
				url.QueryUnescape(
					target,
				)

			if err == nil {
				value =
					decoded
			}
		}
	}

	return strings.TrimSpace(
		value,
	)
}

func publicHTTPURL(
	value string,
) bool {
	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return false
	}

	switch strings.ToLower(
		parsed.Scheme,
	) {

	case "http",
		"https":

	default:
		return false
	}

	host :=
		strings.ToLower(
			strings.TrimSpace(
				parsed.Hostname(),
			),
		)

	if host == "" {
		return false
	}

	if host == "duckduckgo.com" ||
		strings.HasSuffix(
			host,
			".duckduckgo.com",
		) {

		return false
	}

	return true
}

func hostFromURL(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	parsed, err :=
		url.Parse(
			value,
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

	return host
}

func canonicalURL(
	value string,
) string {
	value =
		cleanURL(
			value,
		)

	if value == "" {
		return ""
	}

	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return ""
	}

	scheme :=
		strings.ToLower(
			strings.TrimSpace(
				parsed.Scheme,
			),
		)

	if scheme != "http" &&
		scheme != "https" {

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

	if host == "" {
		return ""
	}

	port :=
		parsed.Port()

	//
	// Default ports do not distinguish pages.
	//
	if scheme == "http" &&
		port == "80" {

		port = ""
	}

	if scheme == "https" &&
		port == "443" {

		port = ""
	}

	parsed.Scheme =
		scheme

	if port != "" {

		parsed.Host =
			net.JoinHostPort(
				host,
				port,
			)

	} else {

		parsed.Host =
			host
	}

	//
	// Fragments never identify a different
	// server-side resource.
	//
	parsed.Fragment =
		""

	//
	// Remove common tracking parameters while
	// keeping parameters that may actually
	// identify page content.
	//
	query :=
		parsed.Query()

	for key := range query {

		if trackingParameter(
			key,
		) {

			query.Del(
				key,
			)
		}
	}

	parsed.RawQuery =
		query.Encode()

	//
	// Treat:
	//
	// example.com/person
	//
	// and:
	//
	// example.com/person/
	//
	// as the same resource.
	//
	if parsed.Path == "/" {

		parsed.Path =
			""

	} else {

		parsed.Path =
			strings.TrimSuffix(
				parsed.Path,
				"/",
			)
	}

	return parsed.String()
}

func trackingParameter(
	value string,
) bool {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	if strings.HasPrefix(
		value,
		"utm_",
	) {
		return true
	}

	switch value {

	case "fbclid",
		"gclid",
		"dclid",
		"msclkid",
		"mc_cid",
		"mc_eid",
		"igshid":

		return true
	}

	return false
}
