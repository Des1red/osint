package enrich

import (
	"strings"

	"osint/internal/enrich/corporate"
	"osint/internal/enrich/directory"
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

	// //
	// // Stage 2:
	// //
	// // General web discovery.
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
	// Corporate intelligence.
	//
	// The current corporate providers require
	// a real-world name anchor.
	//
	if strings.TrimSpace(
		input.FullName,
	) != "" {

		logger.Info(
			"Starting corporate enrichment...",
		)

		result,
			err =
			corporate.Enrich(
				result,
				input,
			)

		if err != nil {

			logger.LogError(
				"Corporate enrichment failed",
				err.Error(),
			)

		} else {

			logger.Info(
				"Corporate enrichment complete.",
			)
		}
	}

	//
	// Stage 4:
	//
	// Public directory discovery.
	//
	if strings.TrimSpace(
		input.Surname,
	) != "" {

		logger.Info(
			"Starting directory enrichment...",
		)

		directoryResults,
			directoryErr :=
			directory.Enrich(
				input,
			)

		if directoryErr != nil {

			logger.LogError(
				"Directory enrichment failed",
				directoryErr.Error(),
			)

		} else {

			result.Directories =
				append(
					result.Directories,
					directoryResults...,
				)

			logger.Info(
				"Directory enrichment complete.",
			)
		}
	} else {
		logger.Info(
			"Skipped directory enrichment, no surname set",
		)
	}

	return result,
		nil
}
