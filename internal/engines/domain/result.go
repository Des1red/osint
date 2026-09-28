package domain

import (
	"osint/internal/engines/domain/takeover"
)

type MXRecord struct {
	Preference uint16

	Host string
}

type TakeoverResult = takeover.Result

type SubdomainResult struct {
	Host string

	A []string

	AAAA []string

	CNAME string

	Takeover TakeoverResult
}

type DomainResult struct {
	Domain string

	Found bool

	A []string

	AAAA []string

	CNAME string

	MX []MXRecord

	NS []string

	TXT []string

	SubdomainSource string

	//
	// Total unique names discovered from
	// certificate transparency before the
	// bounded DNS-resolution limit.
	//
	SubdomainCandidates int

	//
	// Number of discovered names actually
	// checked through DNS.
	//
	SubdomainChecked int

	//
	// Currently DNS-relevant subdomains.
	//
	// Stale CT names with no current A,
	// AAAA or CNAME record are discarded.
	//
	Subdomains []SubdomainResult
}
