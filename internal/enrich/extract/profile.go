package extract

import (
	"net/url"
	"strings"
)

// DirectSocialProfileFromURL returns a social
// account only when the supplied URL itself is
// a profile URL.
//
// SocialFromURL is intentionally broader:
//
//	/@person/video/123
//
// can still reveal:
//
//	@person
//
// which is useful OSINT.
//
// That broader relationship must not establish
// that @person is owned by the investigated
// identity.
func DirectSocialProfileFromURL(
	value string,
) (
	string,
	string,
	bool,
) {
	platform,
		username,
		ok :=
		SocialFromURL(
			value,
		)

	if !ok {

		return "",
			"",
			false
	}

	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {

		return "",
			"",
			false
	}

	path :=
		strings.Trim(
			parsed.Path,
			"/",
		)

	if path == "" {

		return "",
			"",
			false
	}

	parts :=
		strings.Split(
			path,
			"/",
		)

	switch platform {

	case "github",
		"gitlab":

		//
		// /username
		//
		// not:
		//
		// /username/repository
		//
		if len(parts) != 1 {

			return "",
				"",
				false
		}

	case "instagram":

		//
		// /username
		//
		// Instagram content routes such as
		// /p/... and /reel/... are already
		// rejected by SocialFromURL.
		//
		if len(parts) != 1 {

			return "",
				"",
				false
		}

	case "twitter":

		//
		// /username
		//
		// not:
		//
		// /username/status/123
		//
		if len(parts) != 1 {

			return "",
				"",
				false
		}

	case "facebook":

		//
		// /username
		//
		// Content below a username must not
		// establish ownership of that username.
		//
		if len(parts) != 1 {

			return "",
				"",
				false
		}

	case "reddit":

		//
		// /user/username
		// /u/username
		//
		if len(parts) != 2 {

			return "",
				"",
				false
		}

	case "youtube":

		//
		// /@username
		//
		if len(parts) != 1 {

			return "",
				"",
				false
		}

	case "tiktok":

		//
		// /@username
		//
		// not:
		//
		// /@username/video/123
		//
		if len(parts) != 1 {

			return "",
				"",
				false
		}

	case "linkedin":

		//
		// /in/username
		//
		if len(parts) != 2 {

			return "",
				"",
				false
		}

	default:

		return "",
			"",
			false
	}

	return platform,
		username,
		true
}
