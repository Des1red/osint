package api

import (
	"encoding/json"
	"net/http"

	"osint/internal/httpx"
)

type Response struct {
	StatusCode int
	Body       []byte
	FinalURL   string
}

func Get(
	requestURL string,
	headers http.Header,
) (Response, error) {
	response, err :=
		httpx.Do(
			httpx.Request{
				Method: http.MethodGet,

				URL: requestURL,

				Headers: headers,
			},
		)

	if err != nil {
		return Response{},
			err
	}

	return Response{
			StatusCode: response.StatusCode,

			Body: response.Body,

			FinalURL: response.FinalURL,
		},
		nil
}

func Decode(
	body []byte,
	target any,
) error {
	return json.Unmarshal(
		body,
		target,
	)
}
