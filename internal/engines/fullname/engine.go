package fullname

import (
	"strings"

	"osint/internal/logger"
	"osint/internal/models"
)

func FullNameEngine() (
	FullNameResult,
	error,
) {
	logger.HeaderStart(
		"Full Name",
	)

	defer logger.HeaderEnd(
		"Full Name",
	)

	fullName :=
		strings.TrimSpace(
			models.ScopeInput.FullName,
		)

	var result FullNameResult

	logger.Info(
		"Validating full name...",
	)

	err :=
		validate(
			fullName,
		)

	if err != nil {

		return result,
			err
	}

	logger.Info(
		"Full name valid.",
	)

	//
	// Stage 1:
	//
	// Generate full-name lookup candidates.
	//
	logger.Info(
		"Preparing search candidates...",
	)

	searchCandidates :=
		prepareCandidates(
			fullName,
		)

	logger.Info(
		"Search candidates prepared.",
	)

	//
	// Stage 2:
	//
	// Discover possible accounts.
	//
	logger.Info(
		"Looking up possible accounts...",
	)

	result.Candidates =
		lookupCandidates(
			fullName,
			searchCandidates,
		)

	logger.Info(
		"Account lookup complete.",
	)

	//
	// Stage 3:
	//
	// Match discovered accounts against
	// the requested full name.
	//
	logger.Info(
		"Matching discovered accounts...",
	)

	result.Matches =
		matchCandidates(
			fullName,
			result.Candidates,
		)

	logger.Info(
		"Account matching complete.",
	)

	//
	// Stage 4:
	//
	// Enrichment always runs.
	//
	// Matched accounts are optional context,
	// not a requirement.
	//
	logger.Info(
		"Starting full-name enrichment...",
	)

	result.Enrichment =
		enrichFullName(
			fullName,
			result.Matches,
		)

	logger.Info(
		"Full-name enrichment complete.",
	)

	return result,
		nil
}
