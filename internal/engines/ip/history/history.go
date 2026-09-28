package history

import (
	"fmt"
	"sync"
)

func History(
	ip string,
) (
	HistoryResult,
	error,
) {
	result :=
		HistoryResult{
			IP: ip,
		}

	var mutex sync.Mutex

	var waitGroup sync.WaitGroup

	successfulCollectors :=
		0

	waitGroup.Add(
		2,
	)

	go func() {
		defer waitGroup.Done()

		status,
			err :=
			routingStatus(
				ip,
			)

		mutex.Lock()

		defer mutex.Unlock()

		if err != nil {

			result.Errors =
				append(
					result.Errors,
					ProviderError{
						Provider: "RIPEstat routing-status",

						Error: err.Error(),
					},
				)

			return
		}

		result.Status =
			status

		result.Sources =
			appendUnique(
				result.Sources,
				"RIPEstat routing-status",
			)

		successfulCollectors++
	}()

	go func() {
		defer waitGroup.Done()

		routing,
			err :=
			routingHistory(
				ip,
			)

		mutex.Lock()

		defer mutex.Unlock()

		if err != nil {

			result.Errors =
				append(
					result.Errors,
					ProviderError{
						Provider: "RIPEstat routing-history",

						Error: err.Error(),
					},
				)

			return
		}

		result.Routing =
			routing

		result.Sources =
			appendUnique(
				result.Sources,
				"RIPEstat routing-history",
			)

		successfulCollectors++
	}()

	waitGroup.Wait()

	if successfulCollectors == 0 {

		return result,
			fmt.Errorf(
				"all history collectors failed",
			)
	}

	return result,
		nil
}

func appendUnique(
	values []string,
	value string,
) []string {
	if value == "" {
		return values
	}

	for _, existing := range values {

		if existing == value {
			return values
		}
	}

	return append(
		values,
		value,
	)
}
