package searchengines

import (
	"context"
	"fmt"

	"osint/internal/httpx/chrome"
)

type engineFunc func(
	browserContext context.Context,
	queries []string,
) (
	[]Result,
	error,
)

func Search(
	queries []string,
) (
	[]Result,
	error,
) {
	queries =
		uniqueQueries(
			queries,
		)

	if len(queries) == 0 {

		return nil,
			nil
	}

	browserContext,
		err :=
		chrome.Context()

	if err != nil {

		return nil,
			err
	}

	engines :=
		[]engineFunc{
			duckDuckGo,
			brave,
		}

	var result []Result

	successfulEngines :=
		0

	var firstError error

	for _, engine := range engines {

		items,
			err :=
			engine(
				browserContext,
				queries,
			)

		if err != nil {

			if firstError == nil {

				firstError =
					err
			}

			continue
		}

		successfulEngines++

		result =
			append(
				result,
				items...,
			)
	}

	if successfulEngines == 0 &&
		firstError != nil {

		return nil,
			fmt.Errorf(
				"all search engines failed: %w",
				firstError,
			)
	}

	return result,
		nil
}
