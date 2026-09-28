package platforms

import (
	"osint/internal/platforms/facebook"
	"osint/internal/platforms/github"
	"osint/internal/platforms/gitlab"
	"osint/internal/platforms/instagram"
	"osint/internal/platforms/reddit"
	"osint/internal/platforms/tiktok"
	"osint/internal/platforms/twitter"
	"osint/internal/platforms/variants"

	"osint/internal/logger"
)

func UsernameLookup(
	username string,
) variants.PlatformResults {
	return usernameLookup(
		username,
	)
}

func usernameLookup(
	username string,
) variants.PlatformResults {
	var result variants.PlatformResults

	githubResult, err :=
		github.UsernameLookup(
			username,
		)

	if err != nil {
		logger.LogError(
			"GitHub failed for "+username,
			err.Error(),
		)
	} else {
		result.GitHub =
			githubResult
	}

	gitlabResult, err :=
		gitlab.UsernameLookup(
			username,
		)

	if err != nil {
		logger.LogError(
			"GitLab failed for "+username,
			err.Error(),
		)
	} else {
		result.GitLab =
			gitlabResult
	}

	redditResult, err :=
		reddit.UsernameLookup(
			username,
		)

	if err != nil {
		logger.LogError(
			"Reddit failed for "+username,
			err.Error(),
		)
	} else {
		result.Reddit =
			redditResult
	}

	facebookResult, err :=
		facebook.UsernameLookup(
			username,
		)

	if err != nil {
		logger.LogError(
			"Facebook failed for "+username,
			err.Error(),
		)
	} else {
		result.Facebook =
			facebookResult
	}

	instagramResult, err :=
		instagram.UsernameLookup(
			username,
		)

	if err != nil {
		logger.LogError(
			"Instagram failed for "+username,
			err.Error(),
		)
	} else {
		result.Instagram =
			instagramResult
	}

	twitterResult, err :=
		twitter.UsernameLookup(
			username,
		)

	if err != nil {
		logger.LogError(
			"Twitter failed for "+username,
			err.Error(),
		)
	} else {
		result.Twitter =
			twitterResult
	}

	tiktokResult, err :=
		tiktok.UsernameLookup(
			username,
		)

	if err != nil {
		logger.LogError(
			"TikTok failed for "+username,
			err.Error(),
		)
	} else {
		result.TikTok =
			tiktokResult
	}

	return result
}

func VariantLookup(
	username string,
	root variants.PlatformResults,
) []variants.Result {
	candidates :=
		variants.Generate(
			username,
		)

	return variants.Search(
		candidates,
		root,
		usernameLookup,
	)
}
