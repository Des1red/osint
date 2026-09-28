package username

import (
	"osint/internal/enrich"
	"osint/internal/logger"
	"osint/internal/platforms/variants"
)

func buildEnrichment(
	username string,
	result variants.PlatformResults,
) enrich.EnrichmentResult {
	input :=
		enrichmentInput(
			username,
			result,
		)

	enrichmentResult, err :=
		enrich.Enrich(
			input,
		)

	if err != nil {

		logger.LogError(
			"Username enrichment failed",
			err.Error(),
		)

		return enrich.EnrichmentResult{}
	}

	return enrichmentResult
}

func enrichmentInput(
	username string,
	result variants.PlatformResults,
) enrich.Input {
	input :=
		enrich.Input{
			RootUsername: username,

			Usernames: []enrich.UsernameInput{
				{
					Username: result.GitHub.Twitter,

					Source: "GitHub Twitter",
				},
			},

			Text: []enrich.TextInput{
				{
					Text: result.GitHub.Bio,

					Source: "GitHub bio",
				},

				{
					Text: result.Reddit.ProfileDescription,

					Source: "Reddit description",
				},

				{
					Text: result.Facebook.Description,

					Source: "Facebook description",
				},

				{
					Text: result.Instagram.Description,

					Source: "Instagram description",
				},

				{
					Text: result.Twitter.Description,

					Source: "Twitter / X description",
				},

				{
					Text: result.TikTok.Description,

					Source: "TikTok description",
				},
			},

			URLs: []enrich.URLInput{
				{
					URL: result.GitHub.Blog,

					Source: "GitHub blog",
				},
			},
		}

	if result.GitHub.Found {

		input.Accounts =
			append(
				input.Accounts,
				enrich.AccountInput{
					Platform: "GitHub",

					Username: result.GitHub.Username,

					ProfileURL: result.GitHub.ProfileURL,
				},
			)
	}

	if result.GitLab.Found {

		input.Accounts =
			append(
				input.Accounts,
				enrich.AccountInput{
					Platform: "GitLab",

					Username: result.GitLab.Username,

					ProfileURL: result.GitLab.ProfileURL,
				},
			)
	}

	if result.Reddit.Found {

		input.Accounts =
			append(
				input.Accounts,
				enrich.AccountInput{
					Platform: "Reddit",

					Username: result.Reddit.Username,

					ProfileURL: result.Reddit.ProfileURL,
				},
			)
	}

	if result.Facebook.Found {

		input.Accounts =
			append(
				input.Accounts,
				enrich.AccountInput{
					Platform: "Facebook",

					Username: result.Facebook.Username,

					ProfileURL: result.Facebook.ProfileURL,
				},
			)
	}

	if result.Instagram.Found {

		input.Accounts =
			append(
				input.Accounts,
				enrich.AccountInput{
					Platform: "Instagram",

					Username: result.Instagram.Username,

					ProfileURL: result.Instagram.ProfileURL,
				},
			)
	}

	if result.Twitter.Found {

		input.Accounts =
			append(
				input.Accounts,
				enrich.AccountInput{
					Platform: "Twitter",

					Username: result.Twitter.Username,

					ProfileURL: result.Twitter.ProfileURL,
				},
			)
	}

	if result.TikTok.Found {

		input.Accounts =
			append(
				input.Accounts,
				enrich.AccountInput{
					Platform: "TikTok",

					Username: result.TikTok.Username,

					ProfileURL: result.TikTok.ProfileURL,
				},
			)
	}

	return input
}
