package brave

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
	queries =
		uniqueQueries(
			queries,
		)

	if len(queries) == 0 {

		return nil,
			nil
	}

	responses :=
		runBraveSearchMany(
			browserContext,
			queries,
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
			parseBraveResults(
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
				"Brave searches failed: %w",
				firstError,
			)
	}

	return result,
		nil
}
