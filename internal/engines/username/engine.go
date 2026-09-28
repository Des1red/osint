package username

import (
	"osint/internal/models"
	"osint/internal/platforms"
)

func UsernameEngine() (
	UsernameResult,
	error,
) {
	username :=
		models.ScopeInput.Username

	var result UsernameResult

	err :=
		validate(
			username,
		)

	if err != nil {

		return result,
			err
	}

	//
	// Stage 1:
	//
	// Exact username lookup.
	//
	result.PlatformResults =
		platforms.UsernameLookup(
			username,
		)

	//
	// Stage 2:
	//
	// Enrichment from exact results.
	//
	result.Enrichment =
		buildEnrichment(
			username,
			result.PlatformResults,
		)

	//
	// Stage 3:
	//
	// Generated username variants.
	//
	result.Variants =
		platforms.VariantLookup(
			username,
			result.PlatformResults,
		)

	return result,
		nil
}
