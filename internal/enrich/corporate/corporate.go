package corporate

import (
	"fmt"
	"strings"

	"osint/internal/enrich/corporate/opencorp"
	"osint/internal/enrich/dedupe"
	"osint/internal/enrich/model"
	"osint/internal/logger"
)

func Enrich(
	result model.EnrichmentResult,
	input model.Input,
) (
	model.EnrichmentResult,
	error,
) {
	fullName :=
		strings.TrimSpace(
			input.FullName,
		)

	if fullName == "" {

		return result,
			nil
	}

	logger.HeaderStart(
		"Corporate Enrichment",
	)

	defer logger.HeaderEnd(
		"Corporate Enrichment",
	)

	//
	// Stage 1:
	//
	// OpenCorporates.
	//
	logger.Info(
		"Searching OpenCorporates...",
	)

	openCorpResult,
		searchErr :=
		opencorp.Search(
			fullName,
		)

	if len(openCorpResult.Officers) == 0 {

		if searchErr != nil {

			return result,
				searchErr
		}

		logger.Info(
			"No OpenCorporates officer candidates found.",
		)

		return result,
			nil
	}

	logger.Info(
		fmt.Sprintf(
			"OpenCorporates returned %d officer candidate(s).",
			len(openCorpResult.Officers),
		),
	)

	result,
		retained :=
		mergeOpenCorp(
			result,
			fullName,
			openCorpResult,
		)

	logger.Info(
		fmt.Sprintf(
			"Retained %d matching corporate candidate(s).",
			retained,
		),
	)

	//
	// Search() may return a partial result if
	// individual officer pages failed.
	//
	// Keep the successful records and report
	// the failed inspections separately.
	//
	if searchErr != nil {

		logger.LogError(
			"Some OpenCorporates officer inspections failed",
			searchErr.Error(),
		)
	}

	//
	// Normalize anything contributed by the
	// corporate providers.
	//
	result =
		dedupe.Result(
			result,
		)

	logger.Info(
		"OpenCorporates processing complete.",
	)

	return result,
		nil
}
