package gr11888

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"osint/internal/httpx"
)

const searchURL = "https://www.11888.gr/white-pages/"

func Search(
	surname string,
) (
	Result,
	error,
) {
	surname =
		strings.TrimSpace(
			surname,
		)

	if surname == "" {

		return Result{},
			nil
	}

	params :=
		url.Values{}

	params.Set(
		"query",
		surname,
	)

	response,
		err :=
		httpx.Get(
			searchURL +
				"?" +
				params.Encode(),
		)

	if err != nil {

		return Result{},
			err
	}

	if response.StatusCode !=
		http.StatusOK {

		return Result{},
			fmt.Errorf(
				"11888 HTTP %d",
				response.StatusCode,
			)
	}

	recordURLs,
		err :=
		parseRecordURLs(
			string(
				response.Body,
			),
		)

	if err != nil {

		return Result{},
			err
	}

	result :=
		Result{
			Query: surname,
		}

	for _, recordURL := range recordURLs {

		entryResponse,
			err :=
			httpx.Get(
				recordURL,
			)

		if err != nil {

			continue
		}

		if entryResponse.StatusCode !=
			http.StatusOK {

			continue
		}

		entry,
			err :=
			parseEntry(
				string(
					entryResponse.Body,
				),
				entryResponse.FinalURL,
			)

		if err != nil {

			continue
		}

		result.Entries =
			append(
				result.Entries,
				entry,
			)
	}

	return result,
		nil
}
