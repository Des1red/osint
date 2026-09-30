package username

import (
	"osint/internal/logger"
	"osint/internal/models"
	"osint/internal/platforms"
)

func UsernameEngine() (
	UsernameResult,
	error,
) {
	logger.HeaderStart(
		"Username",
	)

	defer logger.HeaderEnd(
		"Username",
	)

	username :=
		models.ScopeInput.Username

	var result UsernameResult

	logger.Info(
		"Validating username...",
	)

	err :=
		validate(
			username,
		)

	if err != nil {

		return result,
			err
	}

	logger.Info(
		"Username valid.",
	)

	//
	// Stage 1:
	//
	// Exact username lookup.
	//
	logger.Info(
		"Looking up exact username matches...",
	)

	result.PlatformResults =
		platforms.UsernameLookup(
			username,
		)

	logger.Info(
		"Exact username lookup complete.",
	)

	//
	// Stage 2:
	//
	// Enrichment from exact results.
	//
	logger.Info(
		"Enriching username results...",
	)

	result.Enrichment =
		buildEnrichment(
			username,
			result.PlatformResults,
		)

	logger.Info(
		"Username enrichment complete.",
	)

	//
	// Stage 3:
	//
	// Generated username variants.
	//
	logger.Info(
		"Looking up username variants...",
	)

	result.Variants =
		platforms.VariantLookup(
			username,
			result.PlatformResults,
		)

	logger.Info(
		"Username variant lookup complete.",
	)

	return result,
		nil
}
