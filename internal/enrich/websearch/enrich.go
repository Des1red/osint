package websearch

import (
	"osint/internal/enrich/dedupe"
	"osint/internal/enrich/model"
	"osint/internal/enrich/organize"
	"osint/internal/enrich/provenance"
	"osint/internal/enrich/relatives"
	"osint/internal/enrich/social"
	"osint/internal/enrich/websearch/discovery"
	"osint/internal/enrich/websearch/firsthop"
	"osint/internal/httpx/chrome"
	"osint/internal/logger"
)

func Enrich(
	result model.EnrichmentResult,
	input model.Input,
) (
	model.EnrichmentResult,
	error,
) {
	logger.HeaderStart(
		"WebSearch Enrichment",
	)

	defer logger.HeaderEnd(
		"WebSearch Enrichment",
	)

	//
	// One Chromium process for the entire
	// WebSearch enrichment operation.
	//
	logger.Info(
		"Starting Chromium...",
	)

	err :=
		chrome.Open()

	if err != nil {

		return result,
			err
	}

	defer chrome.Close()

	logger.Info(
		"Chromium ready.",
	)

	//
	// Discovery.
	//
	logger.Info(
		"Searching for public information...",
	)

	result,
		searchResults,
		evidence,
		err :=
		discovery.Enrich(
			result,
			input,
		)

	if err != nil {

		return result,
			err
	}

	logger.Info(
		"Search discovery complete.",
	)

	//
	// FirstHop.
	//
	logger.Info(
		"Inspecting discovered pages...",
	)

	result,
		firstHopEvidence :=
		firsthop.Enrich(
			result,
			input,
			searchResults,
		)

	evidence =
		provenance.Merge(
			evidence,
			firstHopEvidence,
		)

	logger.Info(
		"Page inspection complete.",
	)

	//
	// Associated people / relatives.
	//
	logger.Info(
		"Checking associated people...",
	)

	result,
		relativeEvidence :=
		relatives.Enrich(
			result,
			input,
		)

	evidence =
		provenance.Merge(
			evidence,
			relativeEvidence,
		)

	logger.Info(
		"Associated people complete.",
	)

	//
	// PRE-ORGANIZATION DEDUPE.
	//
	// The same semantic fact may have been
	// discovered through several evidence paths:
	//
	//     neutral Discovery evidence
	//     subject-aware FirstHop evidence
	//     relative FirstHop evidence
	//
	// Merge those copies before semantic
	// ownership is resolved.
	//
	logger.Info(
		"Organizing discovered information...",
	)

	result =
		dedupe.Result(
			result,
		)

	//
	// Determine ownership:
	//
	//     root target
	//     related person
	//     unattributed
	//
	result =
		organize.Organize(
			result,
			input,
			evidence,
		)

	//
	// POST-ORGANIZATION DEDUPE.
	//
	result =
		dedupe.Result(
			result,
		)

	logger.Info(
		"Information organization complete.",
	)

	//
	// Social Circle.
	//
	// This stage does NOT perform discovery.
	//
	// It inspects only social-post links which
	// have already been attributed to the root
	// target or an already-known person.
	//
	logger.Info(
		"Building social circle...",
	)

	result.SocialCircle =
		social.Enrich(
			result,
			input,
		)

	logger.Info(
		"Social circle complete.",
	)

	return result,
		nil
}
