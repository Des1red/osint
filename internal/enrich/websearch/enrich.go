package websearch

import (
	"osint/internal/enrich/dedupe"
	"osint/internal/enrich/model"
	"osint/internal/enrich/organize"
	"osint/internal/enrich/provenance"
	"osint/internal/enrich/relatives"
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
	// This lets one fact carry the complete
	// evidence set into organize.Organize().
	//
	// Example:
	//
	//     Person1.
	//
	// may have:
	//
	//     Discovery evidence with no Subject
	//
	// and:
	//
	//     FirstHop evidence
	//     Subject = Person1
	//
	// Those should become one fact before
	// ownership resolution.
	//
	result =
		dedupe.Result(
			result,
		)

	//
	// All collectors have now contributed their
	// raw facts and those facts have been
	// coalesced by semantic value.
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
	// Organization may move several copies into
	// the same semantic destination.
	//
	// Run dedupe again to normalize the final
	// result buckets.
	//
	result =
		dedupe.Result(
			result,
		)

	return result,
		nil
}
