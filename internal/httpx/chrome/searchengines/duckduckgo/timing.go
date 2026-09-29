package duckduckgo

import (
	"context"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const DuckDuckGoSearchConcurrency = 5

type searchTiming struct {
	MinGap time.Duration

	MaxGap time.Duration
}

var duckDuckGoSearchDelay = searchTiming{
	MinGap: 700 *
		time.Millisecond,

	MaxGap: 700 *
		time.Millisecond,
}

type searchTabFunc func(
	browserContext context.Context,
	query string,
) pageResponse

func runSearchMany(
	browserContext context.Context,
	queries []string,
	searchTab searchTabFunc,
	timing searchTiming,
	concurrency int,
) []pageResponse {
	queries =
		uniqueQueries(
			queries,
		)

	if len(queries) == 0 {

		return nil
	}

	if concurrency <= 0 {

		concurrency =
			DuckDuckGoSearchConcurrency
	}

	responses :=
		make(
			[]pageResponse,
			len(queries),
		)

	semaphore :=
		make(
			chan struct{},
			concurrency,
		)

	launchDelays :=
		searchLaunchDelays(
			len(queries),
			timing,
		)

	var waitGroup sync.WaitGroup

	for index, query := range queries {

		waitGroup.Add(
			1,
		)

		go func(
			index int,
			query string,
		) {
			defer waitGroup.Done()

			delay :=
				launchDelays[index]

			if delay > 0 {

				timer :=
					time.NewTimer(
						delay,
					)

				defer timer.Stop()

				select {

				case <-browserContext.Done():

					return

				case <-timer.C:
				}
			}

			select {

			case semaphore <- struct{}{}:

			case <-browserContext.Done():

				return
			}

			defer func() {

				<-semaphore
			}()

			responses[index] =
				searchTab(
					browserContext,
					query,
				)
		}(
			index,
			query,
		)
	}

	waitGroup.Wait()

	return responses
}

func searchLaunchDelays(
	count int,
	timing searchTiming,
) []time.Duration {
	delays :=
		make(
			[]time.Duration,
			count,
		)

	if count <= 1 ||
		timing.MaxGap <= 0 {

		return delays
	}

	minGap :=
		timing.MinGap

	maxGap :=
		timing.MaxGap

	if minGap < 0 {

		minGap =
			0
	}

	if maxGap < minGap {

		maxGap =
			minGap
	}

	var cumulativeDelay time.Duration

	for index := 1; index < count; index++ {

		gap :=
			minGap

		gapRange :=
			maxGap -
				minGap

		if gapRange > 0 {

			gap +=
				time.Duration(
					rand.Int63n(
						int64(
							gapRange,
						) + 1,
					),
				)
		}

		cumulativeDelay +=
			gap

		delays[index] =
			cumulativeDelay
	}

	return delays
}

func uniqueQueries(
	queries []string,
) []string {
	seen :=
		make(
			map[string]struct{},
		)

	var result []string

	for _, query := range queries {

		query =
			strings.TrimSpace(
				query,
			)

		if query == "" {

			continue
		}

		key :=
			strings.ToLower(
				query,
			)

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		result =
			append(
				result,
				query,
			)
	}

	return result
}
