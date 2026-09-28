package relatives

import (
	"strings"

	"osint/internal/enrich/extract"
	"osint/internal/enrich/model"
	"osint/internal/enrich/provenance"
	"osint/internal/enrich/websearch/discovery"
	"osint/internal/enrich/websearch/firsthop"
	"osint/internal/logger"
)

const maxRelativeSearches = 3

func Enrich(
	result model.EnrichmentResult,
	input model.Input,
) (
	model.EnrichmentResult,
	[]model.Evidence,
) {
	target :=
		strings.TrimSpace(
			input.FullName,
		)

	if target == "" ||
		len(result.People) == 0 {

		return result,
			nil
	}

	//
	// Freeze the currently known investigation
	// graph before relative expansion begins.
	//
	// Example:
	//
	//     Person1
	//     Person2
	//
	// A search performed for Person2 may therefore
	// surface evidence whose actual Subject is
	// Person1.
	//
	// Newly discovered people during this pass
	// are deliberately NOT added to this graph.
	// That keeps relative expansion one level
	// deep.
	//
	subjects :=
		provenance.Subjects(
			input,
			result.People,
		)

	logger.Debug(
		"Relative Investigation Graph",
		"target",
		target,
		"subjects",
		subjects,
	)

	searched :=
		0

	var evidence []model.Evidence

	for index := range result.People {

		if searched >=
			maxRelativeSearches {

			break
		}

		person :=
			result.People[index]

		name :=
			strings.TrimSpace(
				person.Name,
			)

		if name == "" {

			continue
		}

		searched++

		//
		// The discovered person becomes an
		// independent search anchor.
		//
		// Example:
		//
		// Root:
		//     Person1
		//
		// Search anchor:
		//     Person2
		//
		// Evidence found during this search keeps:
		//
		//     Anchor = Person2
		//
		// But FirstHop may determine:
		//
		//     Subject = Person1
		//
		personInput :=
			model.Input{
				FullName: name,
			}

		//
		// Run normal high-recall Discovery for
		// this relative.
		//
		webResult,
			rawResults,
			err :=
			discovery.Search(
				personInput,
			)

		if err != nil {

			logger.LogError(
				"Relative WebSearch failed for "+name,
				err.Error(),
			)

			continue
		}

		logger.Debug(
			"Relative WebSearch",
			"target",
			target,
			"person",
			name,
			"raw results",
			rawResults,
			"ranked results",
			webResult.Results,
		)

		//
		// Preserve the relative's own ranked
		// WebSearch results.
		//
		// These are still neutral search evidence.
		// Search context alone does not establish
		// a Subject.
		//
		searchResults,
			searchEvidence :=
			relativeSearchResults(
				personInput,
				webResult.Results,
			)

		person.SearchResults =
			searchResults

		evidence =
			provenance.Merge(
				evidence,
				searchEvidence,
			)

		//
		// Relative FirstHop.
		//
		// Important distinction:
		//
		// Search anchor:
		//
		//     Person2
		//
		// Known graph:
		//
		//     Person1
		//     Person2
		//
		// Therefore a result discovered from a
		// Person2 query may still be accepted when its
		// actual page independently matches Person1.
		//
		relativeResult,
			firstHopEvidence :=
			firsthop.EnrichSubjects(
				model.EnrichmentResult{},
				personInput,
				rawResults,
				subjects,
			)

		logger.Debug(
			"Relative FirstHop",
			"target",
			target,
			"anchor",
			name,
			"subjects",
			subjects,
			"usernames",
			len(
				relativeResult.Usernames,
			),
			"emails",
			len(
				relativeResult.Emails,
			),
			"phones",
			len(
				relativeResult.Phones,
			),
			"socials",
			len(
				relativeResult.Socials,
			),
			"links",
			len(
				relativeResult.Links,
			),
			"locations",
			len(
				relativeResult.Locations,
			),
			"organizations",
			len(
				relativeResult.Organizations,
			),
			"employment",
			len(
				relativeResult.Employment,
			),
			"education",
			len(
				relativeResult.Education,
			),
			"evidence count",
			len(
				firstHopEvidence,
			),
		)

		//
		// Preserve all subject-aware neutral
		// evidence generated while searching this
		// person.
		//
		evidence =
			provenance.Merge(
				evidence,
				firstHopEvidence,
			)

		//
		// Merge all raw facts into the shared
		// enrichment pool.
		//
		// Ownership is intentionally NOT decided
		// here.
		//
		result =
			mergeRelativeFirstHop(
				result,
				relativeResult,
			)

		//
		// Do not guess family relationships.
		//
		// Search the complete raw second-hop
		// result set for an explicit statement
		// connecting the root target and this
		// person.
		//
		if !person.RelationVerified {

			relation,
				relationEvidence,
				verified :=
				validateRelationship(
					target,
					personInput,
					rawResults,
				)

			if verified {

				person.Relation =
					relation

				person.RelationVerified =
					true

				person.RelationSource =
					relationEvidence.Source

				person.RelationEvidence =
					append(
						person.RelationEvidence,
						model.EvidenceReference{
							EvidenceID: relationEvidence.ID,
						},
					)

				//
				// Relationship validation may use
				// a raw search result which did not
				// survive normal ranking.
				//
				// Preserve that evidence separately
				// so RelationEvidence always resolves.
				//
				evidence =
					provenance.Merge(
						evidence,
						[]model.Evidence{
							relationEvidence,
						},
					)
			}
		}

		result.People[index] =
			person
	}

	return result,
		provenance.Merge(
			evidence,
		)
}

func relativeSearchResults(
	input model.Input,
	values []discovery.SearchResult,
) (
	[]model.RelativeSearchResult,
	[]model.Evidence,
) {
	results :=
		make(
			[]model.RelativeSearchResult,
			0,
			len(values),
		)

	evidence :=
		make(
			[]model.Evidence,
			0,
			len(values),
		)

	for _, value := range values {

		itemEvidence :=
			relativeEvidence(
				input,
				value,
			)

		evidence =
			append(
				evidence,
				itemEvidence,
			)

		results =
			append(
				results,
				model.RelativeSearchResult{
					Query: value.Query,

					Title: value.Title,

					URL: value.URL,

					Snippet: value.Snippet,

					Engine: value.Engine,

					EvidenceID: itemEvidence.ID,
				},
			)
	}

	return results,
		provenance.Merge(
			evidence,
		)
}

func relativeEvidence(
	input model.Input,
	item discovery.SearchResult,
) model.Evidence {
	//
	// Relative search evidence remains neutral
	// until it is independently matched against
	// the known graph.
	//
	// Anchor says only:
	//
	//     this result was discovered while
	//     searching Person2.
	//
	// It does NOT say that every fact belongs to
	// Person2.
	//
	evidence :=
		discovery.SearchEvidence(
			input,
			item,
		)

	evidence.Source =
		resultSource(
			item,
		)

	return evidence
}

func validateRelationship(
	target string,
	personInput model.Input,
	results []discovery.SearchResult,
) (
	string,
	model.Evidence,
	bool,
) {
	person :=
		strings.TrimSpace(
			personInput.FullName,
		)

	if target == "" ||
		person == "" {

		return "",
			model.Evidence{},
			false
	}

	for _, item := range results {

		itemEvidence :=
			relativeEvidence(
				personInput,
				item,
			)

		if itemEvidence.Text == "" {

			continue
		}

		relation,
			verified :=
			extract.FamilyRelation(
				target,
				person,
				itemEvidence.Text,
			)

		if !verified {

			continue
		}

		return relation,
			itemEvidence,
			true
	}

	return "",
		model.Evidence{},
		false
}

func resultSource(
	item discovery.SearchResult,
) string {
	source :=
		"Relative WebSearch"

	if strings.TrimSpace(
		item.Engine,
	) != "" {

		source +=
			" [" +
				strings.TrimSpace(
					item.Engine,
				) +
				"]"
	}

	if strings.TrimSpace(
		item.Query,
	) != "" {

		source +=
			" | Query: " +
				strings.TrimSpace(
					item.Query,
				)
	}

	if strings.TrimSpace(
		item.URL,
	) != "" {

		source +=
			" | URL: " +
				strings.TrimSpace(
					item.URL,
				)
	}

	return source
}
