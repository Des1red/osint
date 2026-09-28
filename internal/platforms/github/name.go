package github

import (
	"fmt"
	"net/http"
	"net/url"

	sharedapi "osint/internal/platforms/shared/api"
)

const maxNameResults = 10

type githubSearchResponse struct {
	Items []githubSearchUser `json:"items"`
}

type githubSearchUser struct {
	Login string `json:"login"`
}

func NameLookup(
	name string,
) ([]GitHubResult, error) {
	requestURL :=
		"https://api.github.com/search/users?q=" +
			url.QueryEscape(name) +
			"&per_page=" +
			fmt.Sprint(maxNameResults)

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
		return nil,
			err
	}

	if response.StatusCode !=
		http.StatusOK {

		return nil,
			fmt.Errorf(
				"HTTP %d",
				response.StatusCode,
			)
	}

	var search githubSearchResponse

	err =
		sharedapi.Decode(
			response.Body,
			&search,
		)

	if err != nil {
		return nil,
			err
	}

	results :=
		make(
			[]GitHubResult,
			0,
			len(search.Items),
		)

	for _, candidate := range search.Items {

		if candidate.Login == "" {
			continue
		}

		result, err :=
			UsernameLookup(
				candidate.Login,
			)

		if err != nil {
			continue
		}

		if !result.Found {
			continue
		}

		results =
			append(
				results,
				result,
			)
	}

	return results,
		nil
}
