package fullname

import (
	"osint/internal/enrich"
)

type Candidate struct {
	SearchCandidate string

	Platform string

	Username string

	Name string

	Description string

	Email string

	Location string

	Organization string

	Website string

	ProfileURL string
}

type Match struct {
	Platform string

	Username string

	Name string

	Description string

	Email string

	Location string

	Organization string

	Website string

	ProfileURL string
}

type FullNameResult struct {
	Candidates []Candidate

	Matches []Match

	Enrichment enrich.EnrichmentResult
}
