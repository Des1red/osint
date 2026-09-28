package knowledge

import (
	"osint/internal/platforms/facebook"
	"osint/internal/platforms/github"
	"osint/internal/platforms/gitlab"
	"osint/internal/platforms/instagram"
	"osint/internal/platforms/reddit"
	"osint/internal/platforms/tiktok"
	"osint/internal/platforms/twitter"
)

//
// Platform results.
//

type FacebookResult = facebook.FacebookResult

type GitHubResult = github.GitHubResult

type GitLabResult = gitlab.GitLabResult

type InstagramResult = instagram.InstagramResult

type RedditResult = reddit.RedditResult

type TikTokResult = tiktok.TikTokResult

type TwitterResult = twitter.TwitterResult
