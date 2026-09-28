package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	sharedapi "osint/internal/platforms/shared/api"
)

type gitLabResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	State    string `json:"state"`
	Locked   bool   `json:"locked"`

	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

func UsernameLookup(
	username string,
) (GitLabResult, error) {
	requestURL :=
		"https://gitlab.com/api/v4/users?username=" +
			url.QueryEscape(username)

	headers :=
		make(
			http.Header,
		)

	headers.Set(
		"User-Agent",
		"osint-master",
	)

	headers.Set(
		"Accept",
		"application/json",
	)

	response, err :=
		sharedapi.Get(
			requestURL,
			headers,
		)

	if err != nil {
		return GitLabResult{},
			err
	}

	if response.StatusCode !=
		http.StatusOK {

		return GitLabResult{},
			fmt.Errorf(
				"HTTP %d",
				response.StatusCode,
			)
	}

	var data []gitLabResponse

	err =
		sharedapi.Decode(
			response.Body,
			&data,
		)

	if err != nil {
		return GitLabResult{},
			err
	}

	for _, user := range data {

		if !strings.EqualFold(
			user.Username,
			username,
		) {
			continue
		}

		return GitLabResult{
				Found: true,

				ID: user.ID,

				Username: user.Username,

				Name: user.Name,

				State: user.State,

				Locked: user.Locked,

				AvatarURL: user.AvatarURL,

				ProfileURL: user.WebURL,

				Source: requestURL,
			},
			nil
	}

	return GitLabResult{
			Found: false,

			Source: requestURL,
		},
		nil
}
