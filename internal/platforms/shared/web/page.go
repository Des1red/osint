package web

import (
	"strings"

	"osint/internal/httpx"
)

type ProfilePage struct {
	StatusCode int

	FinalURL string

	Body string

	Content string
}

func Fetch(
	profileURL string,
) (ProfilePage, error) {
	response, err :=
		httpx.Get(
			profileURL,
		)

	if err != nil {
		return ProfilePage{},
			err
	}

	body :=
		string(
			response.Body,
		)

	return ProfilePage{
			StatusCode: response.StatusCode,

			FinalURL: response.FinalURL,

			Body: body,

			Content: strings.ToLower(
				body,
			),
		},
		nil
}
