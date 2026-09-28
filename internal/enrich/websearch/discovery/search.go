package discovery

import (
	"fmt"
	"strings"
	"sync"

	"osint/internal/httpx/chrome/searchengines"
)

const (
	//
	// Final results passed into enrichment
	// after ranking and diversity selection.
	//
	maxResults = 15
)

var searchCacheMu sync.RWMutex

var searchCache = make(
	map[string][]SearchResult,
)

func search(
	queries []string,
) (
	Result,
	error,
) {
	queries =
		uniqueQueries(
			queries,
		)

	if len(queries) == 0 {
		return Result{},
			nil
	}

	resultsByQuery :=
		make(
			map[string][]SearchResult,
		)

	var pending []string

	cachedQueries :=
		0

	//
	// Check the in-memory query cache first.
	//
	for _, query := range queries {

		cached,
			ok :=
			cachedSearch(
				query,
			)

		if ok {

			resultsByQuery[searchCacheKey(
				query,
			)] =
				cached

			cachedQueries++

			continue
		}

		pending =
			append(
				pending,
				query,
			)
	}

	//
	// Discovery knows only that it is asking
	// the search-engine orchestrator for web
	// search results.
	//
	// It does not know which engines are being
	// used or how their pages are parsed.
	//
	if len(pending) > 0 {

		items, err :=
			searchengines.Search(
				pending,
			)

		if err != nil {

			//
			// Preserve already-cached data if a
			// temporary search-engine failure
			// occurs.
			//
			if cachedQueries == 0 {

				return Result{},
					fmt.Errorf(
						"all web searches failed: %w",
						err,
					)
			}

		} else {

			for _, item := range items {

				key :=
					searchCacheKey(
						item.Query,
					)

				if key == "" {
					continue
				}

				resultsByQuery[key] =
					append(
						resultsByQuery[key],
						item,
					)
			}

			//
			// Cache only queries for which at
			// least one result was returned.
			//
			// This avoids turning a temporary
			// failed/empty provider response into
			// a permanent empty cache entry for
			// this process.
			//
			for _, query := range pending {

				key :=
					searchCacheKey(
						query,
					)

				values,
					exists :=
					resultsByQuery[key]

				if !exists ||
					len(values) == 0 {

					continue
				}

				storeSearchCache(
					query,
					values,
				)
			}
		}
	}

	var result Result

	//
	// Restore original query order.
	//
	// No final result cap is applied here.
	// Ranking happens afterwards.
	//
	for _, query := range queries {

		items :=
			resultsByQuery[searchCacheKey(
				query,
			)]

		result.Results =
			append(
				result.Results,
				items...,
			)
	}

	return result,
		nil
}

func searchCacheKey(
	query string,
) string {
	return strings.ToLower(
		strings.Join(
			strings.Fields(
				strings.TrimSpace(
					query,
				),
			),
			" ",
		),
	)
}

func cachedSearch(
	query string,
) (
	[]SearchResult,
	bool,
) {
	key :=
		searchCacheKey(
			query,
		)

	searchCacheMu.RLock()

	values,
		exists :=
		searchCache[key]

	searchCacheMu.RUnlock()

	if !exists {
		return nil,
			false
	}

	result :=
		append(
			[]SearchResult(nil),
			values...,
		)

	return result,
		true
}

func storeSearchCache(
	query string,
	results []SearchResult,
) {
	key :=
		searchCacheKey(
			query,
		)

	values :=
		append(
			[]SearchResult(nil),
			results...,
		)

	searchCacheMu.Lock()

	searchCache[key] =
		values

	searchCacheMu.Unlock()
}
