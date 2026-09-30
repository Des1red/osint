package opencorp

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"osint/internal/httpx"
)

const searchURL = "https://opencorporates.com/officers"

const maxResults = 10

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

	response,
		err :=
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

	result :=
		parseResults(
			body,
		)

	return inspectOfficers(
		result,
	)
}
