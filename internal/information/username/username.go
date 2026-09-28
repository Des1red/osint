package username

import (
	enrichmentinfo "osint/internal/information/enrichment"
	"osint/internal/information/platforms"
	"osint/internal/knowledge"
)

func UsernamePrinter() bool {
	hasData :=
		false

	result :=
		knowledge.Data.Username

	if result.GitHub.Found {
		platforms.PrintGitHub(
			result.GitHub,
		)

		hasData =
			true
	} else {
		platforms.PrintNothingFound(
			"GitHub",
		)
	}

	if result.GitLab.Found {
		platforms.PrintGitLab(
			result.GitLab,
		)

		hasData =
			true
	} else {
		platforms.PrintNothingFound(
			"GitLab",
		)
	}

	if result.Reddit.Found {
		platforms.PrintReddit(
			result.Reddit,
		)

		hasData =
			true
	} else {
		platforms.PrintNothingFound(
			"Reddit",
		)
	}

	if result.Facebook.Found {
		platforms.PrintFacebook(
			result.Facebook,
		)

		hasData =
			true
	} else {
		platforms.PrintNothingFound(
			"Facebook",
		)
	}

	if result.Instagram.Found {
		platforms.PrintInstagram(
			result.Instagram,
		)

		hasData =
			true
	} else {
		platforms.PrintNothingFound(
			"Instagram",
		)
	}

	if result.Twitter.Found {
		platforms.PrintTwitter(
			result.Twitter,
		)

		hasData =
			true
	} else {
		platforms.PrintNothingFound(
			"Twitter / X",
		)
	}

	if result.TikTok.Found {
		platforms.PrintTikTok(
			result.TikTok,
		)

		hasData =
			true
	} else {
		platforms.PrintNothingFound(
			"TikTok",
		)
	}

	if enrichmentinfo.Print(
		result.Enrichment,
	) {
		hasData =
			true
	}

	if len(
		result.Variants,
	) > 0 {

		printVariants(
			result.Variants,
		)

		hasData =
			true
	}

	return hasData
}
