package variants

import (
	"osint/internal/matcher"
	"osint/internal/platforms/facebook"
	"osint/internal/platforms/github"
	"osint/internal/platforms/gitlab"
	"osint/internal/platforms/instagram"
	"osint/internal/platforms/reddit"
	"osint/internal/platforms/tiktok"
	"osint/internal/platforms/twitter"
)

type PlatformResults struct {
	GitHub github.GitHubResult

	GitLab gitlab.GitLabResult

	Reddit reddit.RedditResult

	Facebook facebook.FacebookResult

	Instagram instagram.InstagramResult

	Twitter twitter.TwitterResult

	TikTok tiktok.TikTokResult
}

func (
	result PlatformResults,
) FoundAny() bool {
	return result.GitHub.Found ||
		result.GitLab.Found ||
		result.Reddit.Found ||
		result.Facebook.Found ||
		result.Instagram.Found ||
		result.Twitter.Found ||
		result.TikTok.Found
}

type Result struct {
	Username string

	Match matcher.Level

	Source string

	Platforms PlatformResults
}
