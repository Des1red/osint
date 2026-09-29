package brave

import (
	"math/rand"
	"strings"
	"time"
)

const DefaultSearchConcurrency = 3

const braveVerificationTimeout = 2 * time.Minute

type searchTiming struct {
	MinGap time.Duration

	MaxGap time.Duration
}

var humanSearchDelay = searchTiming{
	MinGap: 2500 *
		time.Millisecond,

	MaxGap: 6 *
		time.Second,
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
