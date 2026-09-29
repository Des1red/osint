package duckduckgo

import (
	"context"
	"fmt"
)

func Search(
	browserContext context.Context,
	queries []string,
) (
	[]Result,
	error,
) {
	responses :=
		runSearchMany(
			browserContext,
			queries,
			duckDuckGoSearchTab,
			duckDuckGoSearchDelay,
			DuckDuckGoSearchConcurrency,
		)

	var result []Result

	successfulQueries :=
		0

	var firstError error

	for _, response := range responses {

		if response.Err != nil {

			if firstError == nil {

				firstError =
					response.Err
			}

			continue
		}

		items,
			err :=
			parseDuckDuckGoResults(
				response.HTML,
			)

		if err != nil {

			if firstError == nil {

				firstError =
					err
			}

			continue
		}

		successfulQueries++

		for index := range items {

			items[index].Query =
				response.Query
		}

		result =
			append(
				result,
				items...,
			)
	}

	if successfulQueries == 0 &&
		firstError != nil {

		return nil,
			fmt.Errorf(
				"DuckDuckGo searches failed: %w",
				firstError,
			)
	}

	return result,
		nil
}
