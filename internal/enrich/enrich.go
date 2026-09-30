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
	logger.HeaderStart(
		"Enrichment",
	)

	defer logger.HeaderEnd(
		"Enrichment",
	)

	//
	// Stage 1:
	//
	// Extract information from evidence already
	// supplied by the calling engine.
	//
	logger.Info(
		"Extracting supplied evidence...",
	)

	result :=
		extract.Extract(
			input,
		)

	logger.Info(
		"Evidence extraction complete.",
	)

	//
	// Stage 2:
	//
	// General web discovery.
	//
	logger.Info(
		"Starting web discovery...",
	)

	result,
		err :=
		websearch.Enrich(
			result,
			input,
		)

	if err != nil {

		logger.LogError(
			"Web search enrichment failed",
			err.Error(),
		)

	} else {

		logger.Info(
			"Web discovery complete.",
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
