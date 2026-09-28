package gitlab

import (
	"fmt"
	"net/http"
	"net/url"

	sharedapi "osint/internal/platforms/shared/api"
)

const maxNameResults = 10

func NameLookup(
	name string,
) ([]GitLabResult, error) {
	requestURL :=
		"https://gitlab.com/api/v4/users?search=" +
			url.QueryEscape(name) +
			"&per_page=" +
			fmt.Sprint(maxNameResults)

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

	var data []gitLabResponse

	err =
		sharedapi.Decode(
			response.Body,
			&data,
		)

	if err != nil {
		return nil,
			err
	}

	results :=
		make(
			[]GitLabResult,
			0,
			len(data),
		)

	for _, user := range data {
		results =
			append(
				results,
				GitLabResult{
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
			)
	}

	return results,
		nil
}
