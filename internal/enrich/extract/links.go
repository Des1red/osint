package extract

import (
	"net/url"
	"regexp"
	"strings"

	"osint/internal/enrich/model"
)

var linkPattern = regexp.MustCompile(
	`(?i)(?:https?://|www\.)[^\s<>"']+|(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z]{2,}(?:/[^\s<>"']*)?`,
)

func extractLinks(
	input model.Input,
) (
	[]model.ExternalLink,
	[]model.SocialReference,
) {
	var links []model.ExternalLink
	var socials []model.SocialReference

	//
	// Explicit URLs already supplied by
	// another collector are trusted as URL
	// evidence.
	//
	for _, value := range input.URLs {

		processLink(
			value.URL,
			value.Source,
			wholeEvidenceReference(
				value.EvidenceID,
			),
			&links,
			&socials,
		)
	}

	//
	// URLs discovered inside free text need
	// additional validation because dotted
	// usernames and ordinary dotted text can
	// resemble bare domains.
	//
	for _, value := range input.Text {

		//
		// Keep the original text intact so byte
		// offsets remain valid.
		//
		// Email ranges are skipped instead of
		// removing email text first.
		//
		emailRanges :=
			emailPattern.FindAllStringIndex(
				value.Text,
				-1,
			)

		matches :=
			linkPattern.FindAllStringIndex(
				value.Text,
				-1,
			)

		for _, matchIndex := range matches {

			if len(matchIndex) != 2 {

				continue
			}

			start :=
				matchIndex[0]

			end :=
				matchIndex[1]

			if start < 0 ||
				end < start ||
				end > len(value.Text) {

				continue
			}

			if overlapsAnyRange(
				start,
				end,
				emailRanges,
			) {

				continue
			}

			if start > 0 &&
				value.Text[start-1] == '@' {

				continue
			}

			match :=
				value.Text[start:end]

			lower :=
				strings.ToLower(
					strings.TrimSpace(
						match,
					),
				)

			explicitURL :=
				strings.HasPrefix(
					lower,
					"http://",
				) ||
					strings.HasPrefix(
						lower,
						"https://",
					) ||
					strings.HasPrefix(
						lower,
						"www.",
					)

			if !explicitURL {

				if bareDomainNoise(
					match,
				) {

					continue
				}

				normalized :=
					normalizeURL(
						match,
					)

				if normalized == "" {

					continue
				}

				parsed,
					err :=
					url.Parse(
						normalized,
					)

				if err != nil {

					continue
				}

				if !publicHost(
					parsed.Hostname(),
				) {

					continue
				}
			}

			processLink(
				match,
				value.Source,
				evidenceReference(
					value.EvidenceID,
					start,
					end,
				),
				&links,
				&socials,
			)
		}
	}

	links =
		uniqueLinks(
			links,
		)

	socials =
		uniqueSocials(
			socials,
		)

	return links,
		socials
}

func processLink(
	value string,
	source string,
	evidenceRefs []model.EvidenceReference,
	links *[]model.ExternalLink,
	socials *[]model.SocialReference,
) {
	value =
		canonicalReferenceURL(
			value,
		)

	if value == "" {
		return
	}

	platform,
		username,
		ok :=
		SocialFromURL(
			value,
		)

	if ok {
		profileURL :=
			canonicalSocialURL(
				platform,
				username,
			)

		if profileURL == "" {
			profileURL =
				value
		}

		*socials =
			append(
				*socials,
				model.SocialReference{
					Platform: platform,

					Username: username,

					URL: profileURL,

					Source: source,

					Evidence: evidenceRefs,
				},
			)

		return
	}

	resolvedURL :=
		""

	resolved,
		err :=
		resolveURL(
			value,
		)

	if err == nil &&
		resolved != value {

		resolvedURL =
			resolved

		platform,
			username,
			ok =
			SocialFromURL(
				resolved,
			)

		if ok {
			profileURL :=
				canonicalSocialURL(
					platform,
					username,
				)

			if profileURL == "" {
				profileURL =
					resolved
			}

			*socials =
				append(
					*socials,
					model.SocialReference{
						Platform: platform,

						Username: username,

						URL: profileURL,

						Source: source,

						Evidence: evidenceRefs,
					},
				)

			return
		}
	}

	categoryURL :=
		value

	if resolvedURL != "" {
		categoryURL =
			resolvedURL
	}

	*links =
		append(
			*links,
			model.ExternalLink{
				URL: value,

				ResolvedURL: resolvedURL,

				Category: classifyLink(
					categoryURL,
				),

				Source: source,

				Evidence: evidenceRefs,
			},
		)
}

func SocialFromURL(
	value string,
) (
	string,
	string,
	bool,
) {
	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return "", "", false
	}

	host :=
		strings.ToLower(
			parsed.Hostname(),
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

	if strings.HasSuffix(
		host,
		".linkedin.com",
	) {
		host =
			"linkedin.com"
	}

	path :=
		strings.Trim(
			parsed.Path,
			"/",
		)

	if path == "" {
		return "", "", false
	}

	parts :=
		strings.Split(
			path,
			"/",
		)

	switch host {

	case "github.com":
		return socialResult(
			"github",
			parts[0],
		)

	case "gitlab.com":
		return socialResult(
			"gitlab",
			parts[0],
		)

	case "instagram.com":
		switch strings.ToLower(
			parts[0],
		) {

		case "p",
			"reel",
			"reels",
			"stories",
			"explore",
			"accounts",
			"direct",
			"tv",
			"share",
			"about",
			"developer",
			"legal":

			return "", "", false
		}

		return socialResult(
			"instagram",
			parts[0],
		)

	case "twitter.com",
		"x.com":

		return socialResult(
			"twitter",
			parts[0],
		)

	case "facebook.com":
		if reservedFacebookRoute(
			parts[0],
		) {
			return "", "", false
		}

		return socialResult(
			"facebook",
			parts[0],
		)

	case "reddit.com":
		if len(parts) >= 2 &&
			(parts[0] == "user" ||
				parts[0] == "u") {

			return socialResult(
				"reddit",
				parts[1],
			)
		}

	case "youtube.com":
		if strings.HasPrefix(
			parts[0],
			"@",
		) {
			return socialResult(
				"youtube",
				strings.TrimPrefix(
					parts[0],
					"@",
				),
			)
		}

	case "tiktok.com":
		if strings.HasPrefix(
			parts[0],
			"@",
		) {
			return socialResult(
				"tiktok",
				strings.TrimPrefix(
					parts[0],
					"@",
				),
			)
		}

	case "linkedin.com":
		if len(parts) >= 2 &&
			strings.EqualFold(
				parts[0],
				"in",
			) {

			return socialResult(
				"linkedin",
				parts[1],
			)
		}
	}

	return "", "", false
}

func bareDomainNoise(
	value string,
) bool {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	if value == "" {
		return true
	}

	if strings.ContainsAny(
		value,
		"/?#",
	) {
		return false
	}

	labels :=
		strings.Split(
			value,
			".",
		)

	if len(labels) < 2 {
		return false
	}

	academicLabels :=
		map[string]struct{}{
			"dr":   {},
			"med":  {},
			"prof": {},
			"phd":  {},
			"md":   {},
			"msc":  {},
			"bsc":  {},
		}

	academicCount :=
		0

	for _, label := range labels {

		label =
			strings.TrimSpace(
				label,
			)

		if _, exists :=
			academicLabels[label]; exists {

			academicCount++
		}
	}

	return academicCount >= 2
}

func reservedFacebookRoute(
	value string,
) bool {
	switch strings.ToLower(
		strings.TrimSpace(
			value,
		),
	) {

	case "groups",
		"pages",
		"public",
		"events",
		"marketplace",
		"watch",
		"gaming",
		"reel",
		"reels",
		"stories",
		"photo.php",
		"photos",
		"story.php",
		"permalink.php",
		"profile.php",
		"people",
		"share",
		"sharer",
		"login",
		"help",
		"settings",
		"privacy",
		"about":

		return true
	}

	return false
}

func canonicalSocialURL(
	platform string,
	username string,
) string {
	username =
		normalizeUsername(
			username,
		)

	if !validUsername(
		username,
	) {
		return ""
	}

	escaped :=
		url.PathEscape(
			username,
		)

	switch strings.ToLower(
		strings.TrimSpace(
			platform,
		),
	) {

	case "github":
		return "https://github.com/" +
			escaped

	case "gitlab":
		return "https://gitlab.com/" +
			escaped

	case "instagram":
		return "https://www.instagram.com/" +
			escaped +
			"/"

	case "twitter":
		return "https://x.com/" +
			escaped

	case "facebook":
		return "https://www.facebook.com/" +
			escaped +
			"/"

	case "reddit":
		return "https://www.reddit.com/user/" +
			escaped +
			"/"

	case "youtube":
		return "https://www.youtube.com/@" +
			escaped

	case "tiktok":
		return "https://www.tiktok.com/@" +
			escaped

	case "linkedin":
		return "https://www.linkedin.com/in/" +
			escaped +
			"/"
	}

	return ""
}

func socialResult(
	platform string,
	username string,
) (
	string,
	string,
	bool,
) {
	username =
		normalizeUsername(
			username,
		)

	if !validUsername(
		username,
	) {
		return "", "", false
	}

	return platform,
		username,
		true
}

func uniqueLinks(
	values []model.ExternalLink,
) []model.ExternalLink {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.ExternalLink

	for _, value := range values {

		key :=
			value.URL

		if value.ResolvedURL != "" {
			key =
				value.ResolvedURL
		}

		key =
			strings.ToLower(
				key +
					":" +
					value.Source,
			)

		if index,
			exists :=
			indexes[key]; exists {

			result[index].Evidence =
				mergeEvidence(
					result[index].Evidence,
					value.Evidence,
				)

			if result[index].ResolvedURL == "" {

				result[index].ResolvedURL =
					value.ResolvedURL
			}

			if result[index].Category == "" {

				result[index].Category =
					value.Category
			}

			continue
		}

		indexes[key] =
			len(result)

		result =
			append(
				result,
				value,
			)
	}

	return result
}

func uniqueSocials(
	values []model.SocialReference,
) []model.SocialReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.SocialReference

	for _, value := range values {

		key :=
			strings.ToLower(
				value.Platform +
					":" +
					value.Username +
					":" +
					value.Source,
			)

		if index,
			exists :=
			indexes[key]; exists {

			result[index].Evidence =
				mergeEvidence(
					result[index].Evidence,
					value.Evidence,
				)

			if result[index].URL == "" {

				result[index].URL =
					value.URL
			}

			continue
		}

		indexes[key] =
			len(result)

		result =
			append(
				result,
				value,
			)
	}

	return result
}

func canonicalReferenceURL(
	value string,
) string {
	value =
		normalizeURL(
			value,
		)

	if value == "" {

		return ""
	}

	parsed,
		err :=
		url.Parse(
			value,
		)

	if err != nil {

		return value
	}

	host :=
		strings.ToLower(
			parsed.Hostname(),
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

	if host !=
		"instagram.com" {

		return value
	}

	parts :=
		strings.Split(
			strings.Trim(
				parsed.Path,
				"/",
			),
			"/",
		)

	if len(parts) < 2 {

		return value
	}

	switch strings.ToLower(
		parts[0],
	) {

	case "p",
		"reel",
		"reels",
		"tv":

		return "https://www.instagram.com/" +
			parts[0] +
			"/" +
			parts[1] +
			"/"
	}

	return value
}
