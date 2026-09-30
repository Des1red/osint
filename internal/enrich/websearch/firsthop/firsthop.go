package firsthop

import (
	"strings"

	"osint/internal/enrich/extract"
	"osint/internal/enrich/model"
	"osint/internal/enrich/provenance"
	"osint/internal/enrich/websearch/discovery"
	"osint/internal/logger"
)

func Enrich(
	result model.EnrichmentResult,
	input model.Input,
	searchResults []discovery.SearchResult,
) (
	model.EnrichmentResult,
	[]model.Evidence,
) {
	return EnrichSubjects(
		result,
		input,
		searchResults,
		provenance.Subjects(
			input,
			result.People,
		),
	)
}

// EnrichSubjects runs FirstHop with an explicit
// investigation graph.
//
// This is used by relative expansion where the
// search anchor changes to the related person,
// while the graph still needs to contain the
// original root target and the other already
// known people.
func EnrichSubjects(
	result model.EnrichmentResult,
	input model.Input,
	searchResults []discovery.SearchResult,
	subjects []model.EvidenceSubject,
) (
	model.EnrichmentResult,
	[]model.Evidence,
) {
	searchAnchor :=
		strings.TrimSpace(
			input.FullName,
		)

	if searchAnchor == "" {

		logger.Debug(
			"FirstHop Skipped",
			"reason",
			"full name is empty",
		)

		return result,
			nil
	}

	subjects =
		firstHopSubjects(
			input,
			subjects,
		)

	if len(subjects) == 0 {

		logger.Debug(
			"FirstHop Skipped",
			"anchor",
			searchAnchor,
			"reason",
			"no known subjects",
		)

		return result,
			nil
	}

	candidates :=
		matchingSubjectResults(
			subjects,
			searchResults,
		)

	logger.Debug(
		"FirstHop Candidates",
		"anchor",
		searchAnchor,
		"known subjects",
		subjects,
		"search results",
		len(searchResults),
		"candidate count",
		len(candidates),
		"candidates",
		candidates,
	)

	if len(candidates) == 0 {

		logger.Debug(
			"FirstHop Skipped",
			"anchor",
			searchAnchor,
			"reason",
			"no known-subject candidates",
		)

		return result,
			nil
	}

	var evidence []model.Evidence

	//
	// FirstHop may use raw search results that
	// were not selected by Discovery's ranking.
	//
	// Re-evaluate each candidate against the
	// current graph and preserve those matched
	// subjects on the search evidence itself.
	//
	for _, candidate := range candidates {

		matchedSubjects :=
			matchingSearchResultSubjects(
				candidate,
				subjects,
			)

		if len(matchedSubjects) == 0 {

			continue
		}

		for _, subject := range matchedSubjects {

			candidateEvidence :=
				modelEvidenceFromSearchSubject(
					input,
					candidate,
					subject,
				)

			evidence =
				append(
					evidence,
					candidateEvidence,
				)

			extractionInput :=
				subjectExtractionInput(
					input,
					[]model.EvidenceSubject{
						subject,
					},
				)

			if candidateEvidence.Text != "" {

				extractionInput.Text =
					append(
						extractionInput.Text,
						model.TextInput{
							Text: candidateEvidence.Text,

							Source: candidateEvidence.Source,

							EvidenceID: candidateEvidence.ID,
						},
					)
			}

			if candidateEvidence.URL != "" {

				extractionInput.URLs =
					append(
						extractionInput.URLs,
						model.URLInput{
							URL: candidateEvidence.URL,

							Source: candidateEvidence.Source,

							EvidenceID: candidateEvidence.ID,
						},
					)
			}

			discovered :=
				extract.Extract(
					extractionInput,
				)

			result =
				appendExtractedResult(
					result,
					discovered,
				)
		}
	}
	//
	// Fetch each candidate once and test the
	// resulting page against every known person
	// in the investigation graph.
	//
	pages :=
		inspect(
			searchAnchor,
			subjects,
			candidates,
		)

	logger.Debug(
		"FirstHop Inspection Complete",
		"anchor",
		searchAnchor,
		"candidate count",
		len(candidates),
		"subject-local pages",
		len(pages),
	)

	for _, page := range pages {

		pageEvidence :=
			modelEvidenceFromPage(
				input,
				page,
			)

		evidence =
			append(
				evidence,
				pageEvidence,
			)

		//
		// Each accepted page is now subject-local.
		//
		// Extraction therefore uses the matched
		// subject as FullName rather than the
		// search anchor.
		//
		// Example:
		//
		// Search anchor:
		//     Person2
		//
		// Page subject:
		//     Person1
		//
		// Relationship / organization extraction
		// must run in Person1 context.
		//
		extractionInput :=
			subjectExtractionInput(
				input,
				[]model.EvidenceSubject{
					page.Subject,
				},
			)

		if pageEvidence.Text != "" {

			extractionInput.Text =
				append(
					extractionInput.Text,
					model.TextInput{
						Text: pageEvidence.Text,

						Source: pageEvidence.Source,

						EvidenceID: pageEvidence.ID,
					},
				)
		}

		for _, organization := range page.Organizations {

			extractionInput.Organizations =
				append(
					extractionInput.Organizations,
					model.OrganizationInput{
						Organization: organization,

						Source: pageEvidence.Source,

						EvidenceID: pageEvidence.ID,
					},
				)
		}

		for _, location := range page.Locations {

			extractionInput.Locations =
				append(
					extractionInput.Locations,
					model.LocationInput{
						Location: location,

						Source: pageEvidence.Source,

						EvidenceID: pageEvidence.ID,
					},
				)
		}

		for _, employment := range page.Employment {

			extractionInput.Employment =
				append(
					extractionInput.Employment,
					model.EmploymentInput{
						Title: employment.Title,

						Organization: employment.Organization,

						StartDate: employment.StartDate,

						EndDate: employment.EndDate,

						Current: employment.Current,

						Summary: employment.Summary,

						Source: pageEvidence.Source,

						EvidenceID: pageEvidence.ID,
					},
				)
		}

		for _, education := range page.Education {

			extractionInput.Education =
				append(
					extractionInput.Education,
					model.EducationInput{
						School: education.School,

						Degrees: education.Degrees,

						Majors: education.Majors,

						StartDate: education.StartDate,

						EndDate: education.EndDate,

						Summary: education.Summary,

						Source: pageEvidence.Source,

						EvidenceID: pageEvidence.ID,
					},
				)
		}

		logger.Debug(
			"FirstHop Subject Extraction Input",
			"anchor",
			searchAnchor,
			"subject",
			page.Subject.Value,
			"text items",
			len(extractionInput.Text),
			"urls",
			len(extractionInput.URLs),
			"locations",
			len(extractionInput.Locations),
			"organizations",
			len(extractionInput.Organizations),
			"employment",
			len(extractionInput.Employment),
			"education",
			len(extractionInput.Education),
		)

		discovered :=
			extract.Extract(
				extractionInput,
			)

		logger.Debug(
			"FirstHop Subject Extraction Result",
			"anchor",
			searchAnchor,
			"subject",
			page.Subject.Value,
			"result",
			discovered,
		)

		result =
			appendExtractedResult(
				result,
				discovered,
			)
	}

	return result,
		provenance.Merge(
			evidence,
		)
}

func firstHopSubjects(
	input model.Input,
	values []model.EvidenceSubject,
) []model.EvidenceSubject {
	values =
		append(
			[]model.EvidenceSubject(nil),
			values...,
		)

	anchor :=
		provenance.Anchor(
			input,
		)

	if anchor.Kind !=
		model.EvidenceAnchorUnknown &&
		strings.TrimSpace(
			anchor.Value,
		) != "" {

		values =
			append(
				values,
				model.EvidenceSubject{
					Kind: anchor.Kind,

					Value: anchor.Value,
				},
			)
	}

	return normalizeFirstHopSubjects(
		values,
	)
}

func subjectExtractionInput(
	input model.Input,
	subjects []model.EvidenceSubject,
) model.Input {
	var result model.Input

	subjects =
		normalizeFirstHopSubjects(
			subjects,
		)

	if len(subjects) != 1 {

		return result
	}

	subject :=
		subjects[0]

	if subject.Kind ==
		model.EvidenceAnchorPerson {

		result.FullName =
			strings.TrimSpace(
				subject.Value,
			)
	}

	//
	// RootUsername only belongs to the search
	// anchor.
	//
	// Do not use Person2's username matcher while
	// extracting a Person1 subject page.
	//
	anchor :=
		provenance.Anchor(
			input,
		)

	if subject.Kind ==
		anchor.Kind &&
		strings.EqualFold(
			strings.TrimSpace(
				subject.Value,
			),
			strings.TrimSpace(
				anchor.Value,
			),
		) {

		result.RootUsername =
			input.RootUsername
	}

	return result
}

func appendExtractedResult(
	result model.EnrichmentResult,
	discovered model.EnrichmentResult,
) model.EnrichmentResult {
	result.Usernames =
		append(
			result.Usernames,
			discovered.Usernames...,
		)

	result.Emails =
		append(
			result.Emails,
			discovered.Emails...,
		)

	result.Phones =
		append(
			result.Phones,
			discovered.Phones...,
		)

	result.Socials =
		append(
			result.Socials,
			discovered.Socials...,
		)

	result.Links =
		append(
			result.Links,
			discovered.Links...,
		)

	result.Locations =
		append(
			result.Locations,
			discovered.Locations...,
		)

	result.Organizations =
		append(
			result.Organizations,
			discovered.Organizations...,
		)

	result.Employment =
		append(
			result.Employment,
			discovered.Employment...,
		)

	result.Education =
		append(
			result.Education,
			discovered.Education...,
		)

	result.RelatedAccounts =
		append(
			result.RelatedAccounts,
			discovered.RelatedAccounts...,
		)

	//
	// IMPORTANT:
	//
	// FirstHop runs against an already-known
	// investigation graph.
	//
	// Subject-local extraction must not expand
	// that graph by re-discovering people which
	// appear in the evidence.
	//
	// Example:
	//
	//     Known subjects:
	//         Person1
	//         Person2
	//
	//     Evidence:
	//         Company was founded by
	//         Person2 and Person1
	//
	// When extracting in Person1 context,
	// extractPeople() can see Person2.
	//
	// When extracting in Person2 context,
	// extractPeople() can see Person1.
	//
	// Appending discovered.People here would
	// therefore create:
	//
	//     Person2
	//     Person2
	//     Person1
	//
	// and corrupt later ownership resolution.
	//
	// People discovery belongs to the earlier
	// graph-building stage. FirstHop only
	// enriches facts for subjects already in
	// that graph.
	//

	return result
}
