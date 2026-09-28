package discovery

import (
	"osint/internal/enrich/model"
)

func Search(
	input model.Input,
) (
	Result,
	[]SearchResult,
	error,
) {
	queries :=
		searchQueries(
			input,
		)

	if len(queries) == 0 {
		return Result{},
			nil,
			nil
	}

	webResult, err :=
		search(
			queries,
		)

	if err != nil {
		return Result{},
			nil,
			err
	}

	//
	// Preserve the complete raw result set.
	//
	// Relationship validation may find useful
	// explicit evidence outside the final
	// ranked result set.
	//
	rawResults :=
		append(
			[]SearchResult(nil),
			webResult.Results...,
		)

	//
	// Normal Discovery ranking.
	//
	webResult.Results =
		prioritizeResults(
			input,
			webResult.Results,
		)

	return webResult,
		rawResults,
		nil
}
