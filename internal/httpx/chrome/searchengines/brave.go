package searchengines

import (
	"context"

	braveengine "osint/internal/httpx/chrome/searchengines/brave"
)

func brave(
	browserContext context.Context,
	queries []string,
) (
	[]Result,
	error,
) {
	items,
		err :=
		braveengine.Search(
			browserContext,
			queries,
		)

	if err != nil {

		return nil,
			err
	}

	result :=
		make(
			[]Result,
			0,
			len(items),
		)

	for _, item := range items {

		result =
			append(
				result,
				Result{
					Query: item.Query,

					Title: item.Title,

					URL: item.URL,

					Snippet: item.Snippet,

					Engine: item.Engine,
				},
			)
	}

	return result,
		nil
}
