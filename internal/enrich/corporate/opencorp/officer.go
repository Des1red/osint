package opencorp

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"osint/internal/httpx"
)

const baseURL = "https://opencorporates.com"

func officerFromPath(
	path string,
) Officer {
	path =
		strings.TrimSpace(
			path,
		)

	if path == "" {

		return Officer{}
	}

	return Officer{
		ProfileURL: baseURL +
			path,
	}
}

func inspectOfficers(
	result Result,
) (
	Result,
	error,
) {
	var inspectErrors []error

	for index := range result.Officers {

		officer :=
			result.Officers[index]

		inspected,
			err :=
			inspectOfficer(
				officer,
			)

		if err != nil {

			inspectErrors =
				append(
					inspectErrors,
					err,
				)

			continue
		}

		result.Officers[index] =
			inspected
	}

	return result,
		errors.Join(
			inspectErrors...,
		)
}

func inspectOfficer(
	officer Officer,
) (
	Officer,
	error,
) {
	if strings.TrimSpace(
		officer.ProfileURL,
	) == "" {

		return officer,
			fmt.Errorf(
				"officer profile URL is empty",
			)
	}

	response,
		err :=
		httpx.Get(
			officer.ProfileURL,
		)

	if err != nil {

		return officer,
			fmt.Errorf(
				"failed to inspect %s: %w",
				officer.ProfileURL,
				err,
			)
	}

	if response.StatusCode !=
		http.StatusOK {

		return officer,
			fmt.Errorf(
				"officer %s returned HTTP %d",
				officer.ProfileURL,
				response.StatusCode,
			)
	}

	body :=
		string(
			response.Body,
		)

	officer =
		parseOfficer(
			body,
			officer,
		)

	return officer,
		nil
}
