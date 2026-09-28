package geo

import (
	"fmt"
	"sort"
	"sync"
)

type providerResponse struct {
	Result ProviderResult

	Error ProviderError

	Success bool
}

func Geo(ip string) (GeoResult, error) {
	if cached, ok := getCached(ip); ok {
		return cached, nil
	}

	responses := collectProviders(ip)

	var providers []ProviderResult
	var providerErrors []ProviderError

	for _, response := range responses {
		if response.Success {
			providers = append(
				providers,
				response.Result,
			)

			continue
		}

		providerErrors = append(
			providerErrors,
			response.Error,
		)
	}

	if len(providers) == 0 {
		return GeoResult{}, fmt.Errorf(
			"all geolocation providers failed",
		)
	}

	result := buildConsensus(
		ip,
		providers,
	)

	result.Errors = providerErrors

	setCached(
		ip,
		result,
	)

	return result, nil
}

func collectProviders(
	ip string,
) []providerResponse {
	responses := make(
		chan providerResponse,
		3,
	)

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()

		result, err := ipAPI(ip)

		if err != nil {
			responses <- providerResponse{
				Error: ProviderError{
					Provider: "ipapi.co",
					Error:    err.Error(),
				},
			}

			return
		}

		responses <- providerResponse{
			Result:  result,
			Success: true,
		}
	}()

	go func() {
		defer wg.Done()

		result, err := ipWhoIs(ip)

		if err != nil {
			responses <- providerResponse{
				Error: ProviderError{
					Provider: "ipwho.is",
					Error:    err.Error(),
				},
			}

			return
		}

		responses <- providerResponse{
			Result:  result,
			Success: true,
		}
	}()

	go func() {
		defer wg.Done()

		result, err := freeIPAPI(ip)

		if err != nil {
			responses <- providerResponse{
				Error: ProviderError{
					Provider: "freeipapi.com",
					Error:    err.Error(),
				},
			}

			return
		}

		responses <- providerResponse{
			Result:  result,
			Success: true,
		}
	}()

	go func() {
		wg.Wait()

		close(responses)
	}()

	var results []providerResponse

	for response := range responses {
		results = append(
			results,
			response,
		)
	}

	sort.Slice(
		results,
		func(i, j int) bool {
			return providerName(results[i]) <
				providerName(results[j])
		},
	)

	return results
}

func providerName(
	response providerResponse,
) string {
	if response.Success {
		return response.Result.Provider
	}

	return response.Error.Provider
}
