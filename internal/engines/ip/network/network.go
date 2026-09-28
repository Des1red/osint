package network

import (
	"sync"
)

func Network(
	ip string,
) (NetworkResult, error) {
	result := NetworkResult{
		IP: ip,
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(5)

	go func() {
		defer wg.Done()

		prefix, asns, err :=
			ripeNetworkInfo(ip)

		if err != nil {
			mu.Lock()

			result.Errors = append(
				result.Errors,
				ProviderError{
					Provider: "RIPEstat",
					Error:    err.Error(),
				},
			)

			mu.Unlock()
			return
		}

		asnResults,
			asnErrors :=
			collectASNs(asns)

		mu.Lock()

		result.Prefix = prefix
		result.ASNs = asnResults

		result.Errors = append(
			result.Errors,
			asnErrors...,
		)

		result.Sources =
			appendUnique(
				result.Sources,
				"RIPEstat",
			)

		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		names, err :=
			reverseDNS(ip)

		if err != nil {
			return
		}

		mu.Lock()

		result.ReverseDNS = names

		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		data, err :=
			ipAPI(ip)

		if err != nil {
			mu.Lock()

			result.Errors = append(
				result.Errors,
				ProviderError{
					Provider: "ip-api.com",
					Error:    err.Error(),
				},
			)

			mu.Unlock()
			return
		}

		mu.Lock()

		applyIPAPI(
			&result,
			data,
		)

		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		data, err :=
			ipWho(ip)

		if err != nil {
			//
			// Missing API key is optional,
			// so don't treat it like a failed
			// network lookup.
			//
			return
		}

		mu.Lock()

		applyIPWho(
			&result,
			data,
		)

		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		data, err := manycast(ip)
		if err != nil {
			mu.Lock()

			result.Errors = append(
				result.Errors,
				ProviderError{
					Provider: "manycast.net",
					Error:    err.Error(),
				},
			)

			mu.Unlock()

			return
		}

		mu.Lock()

		result.Anycast = data

		result.Sources = appendUnique(
			result.Sources,
			"manycast.net",
		)

		mu.Unlock()
	}()

	wg.Wait()

	return result, nil
}

func collectASNs(
	asns []uint32,
) ([]ASNResult, []ProviderError) {
	results := make(
		[]ASNResult,
		len(asns),
	)

	var errors []ProviderError

	var mu sync.Mutex
	var wg sync.WaitGroup

	for i, asn := range asns {
		wg.Add(1)

		go func(
			index int,
			asn uint32,
		) {
			defer wg.Done()

			result,
				asnErrors :=
				collectASN(asn)

			results[index] =
				result

			if len(asnErrors) > 0 {
				mu.Lock()

				errors = append(
					errors,
					asnErrors...,
				)

				mu.Unlock()
			}
		}(
			i,
			asn,
		)
	}

	wg.Wait()

	return results, errors
}

func collectASN(
	asn uint32,
) (ASNResult, []ProviderError) {
	result := ASNResult{
		ASN: asn,
	}

	var errors []ProviderError

	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()

		data, err :=
			ripeASOverview(asn)

		if err != nil {
			mu.Lock()

			errors = append(
				errors,
				ProviderError{
					Provider: "RIPEstat AS overview",
					Error:    err.Error(),
				},
			)

			mu.Unlock()
			return
		}

		mu.Lock()

		result.Holder =
			data.Holder

		result.Announced =
			data.Announced

		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		prefixes, err :=
			ripeAnnouncedPrefixes(asn)

		if err != nil {
			mu.Lock()

			errors = append(
				errors,
				ProviderError{
					Provider: "RIPEstat announced prefixes",
					Error:    err.Error(),
				},
			)

			mu.Unlock()
			return
		}

		mu.Lock()

		result.AnnouncedPrefixes =
			prefixes

		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		counts,
			neighbours,
			err :=
			ripeASNNeighbours(asn)

		if err != nil {
			mu.Lock()

			errors = append(
				errors,
				ProviderError{
					Provider: "RIPEstat ASN neighbours",
					Error:    err.Error(),
				},
			)

			mu.Unlock()
			return
		}

		mu.Lock()

		result.NeighbourCounts =
			counts

		result.Neighbours =
			neighbours

		mu.Unlock()
	}()

	wg.Wait()

	return result, errors
}
