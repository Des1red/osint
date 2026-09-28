package domain

import (
	"sort"
	"sync"

	"osint/internal/engines/domain/takeover"
)

const subdomainLookupWorkers = 12

type subdomainJob struct {
	Index int

	Host string
}

func resolveSubdomains(
	hosts []string,
) []SubdomainResult {
	if len(hosts) == 0 {

		return nil
	}

	results :=
		make(
			[]SubdomainResult,
			len(hosts),
		)

	keep :=
		make(
			[]bool,
			len(hosts),
		)

	jobs :=
		make(
			chan subdomainJob,
		)

	var waitGroup sync.WaitGroup

	for worker := 0; worker < subdomainLookupWorkers; worker++ {

		waitGroup.Add(
			1,
		)

		go func() {

			defer waitGroup.Done()

			for job := range jobs {

				result,
					found :=
					resolveSubdomain(
						job.Host,
					)

				results[job.Index] =
					result

				keep[job.Index] =
					found
			}
		}()
	}

	for index, host := range hosts {

		jobs <- subdomainJob{
			Index: index,

			Host: host,
		}
	}

	close(
		jobs,
	)

	waitGroup.Wait()

	filtered :=
		make(
			[]SubdomainResult,
			0,
			len(results),
		)

	for index, result := range results {

		if !keep[index] {

			continue
		}

		filtered =
			append(
				filtered,
				result,
			)
	}

	sort.Slice(
		filtered,
		func(
			left int,
			right int,
		) bool {
			return filtered[left].Host <
				filtered[right].Host
		},
	)

	return filtered
}

func resolveSubdomain(
	host string,
) (
	SubdomainResult,
	bool,
) {
	result :=
		SubdomainResult{
			Host: host,
		}

	//
	// A / AAAA.
	//
	ipv4,
		ipv6,
		err :=
		lookupAddresses(
			host,
		)

	if err == nil {

		result.A =
			ipv4

		result.AAAA =
			ipv6
	}

	//
	// CNAME.
	//
	cname,
		err :=
		lookupCNAME(
			host,
		)

	if err == nil {

		result.CNAME =
			cname
	}

	//
	// Certificate Transparency contains
	// historical names as well.
	//
	// Keep only names that still expose
	// current DNS information.
	//
	if len(result.A) == 0 &&
		len(result.AAAA) == 0 &&
		result.CNAME == "" {

		return result,
			false
	}

	if result.CNAME != "" {

		result.Takeover =
			takeover.Check(
				result.CNAME,
			)
	}

	return result,
		true
}
