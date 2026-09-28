package opencorp

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"osint/internal/httpx"
)

const searchURL = "https://opencorporates.com/officers"

const maxResults = 10

var officerLinkPattern = regexp.MustCompile(
	`(?i)href=["'](/officers/\d+)["']`,
)

func Search(
	fullName string,
) (
	Result,
	error,
) {
	fullName =
		strings.TrimSpace(
			fullName,
		)

	if fullName == "" {
		return Result{},
			nil
	}

	params :=
		url.Values{}

	params.Set(
		"q",
		fullName,
	)

	requestURL :=
		searchURL +
			"?" +
			params.Encode()

	response, err :=
		httpx.Get(
			requestURL,
		)

	if err != nil {
		return Result{},
			err
	}

	if response.StatusCode !=
		http.StatusOK {

		return Result{},
			fmt.Errorf(
				"OpenCorporates HTTP %d",
				response.StatusCode,
			)
	}

	body :=
		string(
			response.Body,
		)

	return parseResults(
			body,
		),
		nil
}

func parseResults(
	body string,
) Result {
	var result Result

	matches :=
		officerLinkPattern.FindAllStringSubmatch(
			body,
			-1,
		)

	seen :=
		make(
			map[string]struct{},
		)

	for _, match := range matches {

		if len(match) < 2 {
			continue
		}

		path :=
			strings.TrimSpace(
				match[1],
			)

		if path == "" {
			continue
		}

		profileURL :=
			"https://opencorporates.com" +
				path

		key :=
			strings.ToLower(
				profileURL,
			)

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		result.Officers =
			append(
				result.Officers,
				Officer{
					ProfileURL: profileURL,
				},
			)

		if len(result.Officers) >=
			maxResults {

			break
		}
	}

	return result
}
