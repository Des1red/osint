package reputation

import (
	"fmt"
	"sync"
)

func Reputation(
	ip string,
) (ReputationResult, error) {
	result := ReputationResult{
		IP: ip,
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	configuredProviders := 0
	successfulProviders := 0

	wg.Add(3)

	go func() {
		defer wg.Done()

		data, configured, err :=
			abuseIPDB(
				ip,
			)

		if !configured {
			return
		}

		mu.Lock()
		configuredProviders++
		mu.Unlock()

		if err != nil {
			mu.Lock()

			result.Errors =
				append(
					result.Errors,
					ProviderError{
						Provider: "abuseipdb",
						Error:    err.Error(),
					},
				)

			mu.Unlock()

			return
		}

		mu.Lock()

		result.AbuseIPDB =
			data

		successfulProviders++

		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		data, configured, err :=
			virusTotal(
				ip,
			)

		if !configured {
			return
		}

		mu.Lock()
		configuredProviders++
		mu.Unlock()

		if err != nil {
			mu.Lock()

			result.Errors =
				append(
					result.Errors,
					ProviderError{
						Provider: "virustotal",
						Error:    err.Error(),
					},
				)

			mu.Unlock()

			return
		}

		mu.Lock()

		result.VirusTotal =
			data

		successfulProviders++

		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		data, configured, err :=
			ipQS(
				ip,
			)

		if !configured {
			return
		}

		mu.Lock()
		configuredProviders++
		mu.Unlock()

		if err != nil {
			mu.Lock()

			result.Errors =
				append(
					result.Errors,
					ProviderError{
						Provider: "ipqualityscore",
						Error:    err.Error(),
					},
				)

			mu.Unlock()

			return
		}

		mu.Lock()

		result.IPQS =
			data

		successfulProviders++

		mu.Unlock()
	}()

	wg.Wait()

	if configuredProviders == 0 {

		return ReputationResult{},
			nil
	}

	if successfulProviders == 0 {
		return result,
			fmt.Errorf(
				"all reputation providers failed",
			)
	}

	return result, nil
}
