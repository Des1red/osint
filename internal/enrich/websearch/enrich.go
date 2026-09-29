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
)

func Enrich(
	result model.EnrichmentResult,
	input model.Input,
) (
	model.EnrichmentResult,
	error,
) {
	//
	// One Chromium process for the entire
	// WebSearch enrichment operation.
	//
	err :=
		chrome.Open()

	if err != nil {

		return result,
			err
	}

	defer chrome.Close()

	//
	// Discovery.
	//
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

	//
	// FirstHop.
	//
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

	//
	// Associated people / relatives.
	//
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

	//
	// Social Circle.
	//
	// This stage does NOT perform discovery.
	//
	// It inspects only social-post links which
	// have already been attributed to the root
	// target or an already-known person.
	//
	result.SocialCircle =
		social.Enrich(
			result,
			input,
		)

	return result,
		nil
}
