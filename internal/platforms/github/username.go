package github

import (
	"fmt"
	"net/http"
	"net/url"

	sharedapi "osint/internal/platforms/shared/api"
)

type githubResponse struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	NodeID    string `json:"node_id"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`

	Name     string `json:"name"`
	Company  string `json:"company"`
	Blog     string `json:"blog"`
	Location string `json:"location"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`

	TwitterUsername string `json:"twitter_username"`

	PublicRepos int `json:"public_repos"`
	PublicGists int `json:"public_gists"`

	Followers int `json:"followers"`
	Following int `json:"following"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	URL       string `json:"url"`
}

func UsernameLookup(
	username string,
) (GitHubResult, error) {
	requestURL :=
		"https://api.github.com/users/" +
			url.PathEscape(username)

	headers :=
		make(
			http.Header,
		)

	headers.Set(
		"Accept",
		"application/vnd.github+json",
	)

	headers.Set(
		"X-GitHub-Api-Version",
		"2026-03-10",
	)

	headers.Set(
		"User-Agent",
		"osint-master",
	)

	response, err :=
		sharedapi.Get(
			requestURL,
			headers,
		)

	if err != nil {
		return GitHubResult{},
			err
	}

	switch response.StatusCode {

	case http.StatusOK:

	case http.StatusNotFound:
		return GitHubResult{
				Found: false,

				Source: requestURL,
			},
			nil

	default:
		return GitHubResult{},
			fmt.Errorf(
				"HTTP %d",
				response.StatusCode,
			)
	}

	var data githubResponse

	err =
		sharedapi.Decode(
			response.Body,
			&data,
		)

	if err != nil {
		return GitHubResult{},
			err
	}

	return GitHubResult{
			Found: true,

			Username: data.Login,

			ID: data.ID,

			NodeID: data.NodeID,

			Type: data.Type,

			SiteAdmin: data.SiteAdmin,

			Name: data.Name,

			Company: data.Company,

			Blog: data.Blog,

			Location: data.Location,

			Email: data.Email,

			Bio: data.Bio,

			Twitter: data.TwitterUsername,

			PublicRepos: data.PublicRepos,

			PublicGists: data.PublicGists,

			Followers: data.Followers,

			Following: data.Following,

			CreatedAt: data.CreatedAt,

			UpdatedAt: data.UpdatedAt,

			AvatarURL: data.AvatarURL,

			ProfileURL: data.HTMLURL,

			APIURL: data.URL,

			Source: requestURL,
		},
		nil
}
