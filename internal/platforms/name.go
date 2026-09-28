package platforms

import (
	"osint/internal/logger"
	"osint/internal/platforms/facebook"
	"osint/internal/platforms/github"
	"osint/internal/platforms/gitlab"
)

type NameCandidate struct {
	SearchCandidate string

	Platform string

	Username string

	Name string

	Description string

	Email string

	Location string

	Organization string

	Website string

	ProfileURL string
}

func NameLookup(
	fullName string,
	searchCandidates []string,
) []NameCandidate {
	var candidates []NameCandidate

	//
	// GitHub.
	//
	githubResults, err :=
		github.NameLookup(
			fullName,
		)

	if err != nil {
		logger.LogError(
			"GitHub name lookup failed for "+fullName,
			err.Error(),
		)
	} else {
		for _, result := range githubResults {

			if !result.Found {
				continue
			}

			candidates =
				append(
					candidates,
					NameCandidate{
						SearchCandidate: fullName,

						Platform: "GitHub",

						Username: result.Username,

						Name: result.Name,

						Description: result.Bio,

						Email: result.Email,

						Location: result.Location,

						Organization: result.Company,

						Website: result.Blog,

						ProfileURL: result.ProfileURL,
					},
				)
		}
	}

	//
	// GitLab.
	//
	gitlabResults, err :=
		gitlab.NameLookup(
			fullName,
		)

	if err != nil {
		logger.LogError(
			"GitLab name lookup failed for "+fullName,
			err.Error(),
		)
	} else {
		for _, result := range gitlabResults {

			if !result.Found {
				continue
			}

			candidates =
				append(
					candidates,
					NameCandidate{
						SearchCandidate: fullName,

						Platform: "GitLab",

						Username: result.Username,

						Name: result.Name,

						ProfileURL: result.ProfileURL,
					},
				)
		}
	}

	//
	// Facebook.
	//
	for _, searchCandidate := range searchCandidates {

		facebookResults, err :=
			facebook.NameLookup(
				searchCandidate,
			)

		if err != nil {
			logger.LogError(
				"Facebook name lookup failed for "+searchCandidate,
				err.Error(),
			)

			continue
		}

		for _, result := range facebookResults {

			if !result.Found {
				continue
			}

			candidates =
				append(
					candidates,
					NameCandidate{
						SearchCandidate: searchCandidate,

						Platform: "Facebook",

						Username: result.Username,

						Name: result.Name,

						Description: result.Description,

						ProfileURL: result.ProfileURL,
					},
				)
		}
	}

	return candidates
}
