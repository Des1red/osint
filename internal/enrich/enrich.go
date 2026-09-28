package enrich

import (
	"osint/internal/enrich/extract"
	"osint/internal/enrich/websearch"
	"osint/internal/logger"
)

func Enrich(
	input Input,
) (
	EnrichmentResult,
	error,
) {
	//
	// Stage 1:
	//
	// Extract information from evidence already
	// supplied by the calling engine.
	//
	result :=
		extract.Extract(
			input,
		)

	//
	// Stage 2:
	//
	// General web discovery.
	//
	result, err :=
		websearch.Enrich(
			result,
			input,
		)

	if err != nil {

		logger.LogError(
			"Web search enrichment failed",
			err.Error(),
		)
	}

	//
	// Stage 3:
	//
	// OpenCorporates.
	//

	return result,
		nil
}
