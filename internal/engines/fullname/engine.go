package fullname

import (
	"strings"

	"osint/internal/models"
)

func FullNameEngine() (
	FullNameResult,
	error,
) {
	fullName :=
		strings.TrimSpace(
			models.ScopeInput.FullName,
		)

	var result FullNameResult

	err :=
		validate(
			fullName,
		)

	if err != nil {

		return result,
			err
	}

	//
	// Stage 1:
	//
	// Generate full-name lookup candidates.
	//
	searchCandidates :=
		prepareCandidates(
			fullName,
		)

	//
	// Stage 2:
	//
	// Discover possible accounts.
	//
	result.Candidates =
		lookupCandidates(
			fullName,
			searchCandidates,
		)

	//
	// Stage 3:
	//
	// Match discovered accounts against
	// the requested full name.
	//
	result.Matches =
		matchCandidates(
			fullName,
			result.Candidates,
		)

	//
	// Stage 4:
	//
	// Enrichment always runs.
	//
	// Matched accounts are optional context,
	// not a requirement.
	//
	result.Enrichment =
		enrichFullName(
			fullName,
			result.Matches,
		)

	return result,
		nil
}
