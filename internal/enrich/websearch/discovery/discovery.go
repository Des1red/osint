package discovery

import (
	"osint/internal/enrich/extract"
	"osint/internal/enrich/model"
	"osint/internal/logger"
)

func Enrich(
	result model.EnrichmentResult,
	input model.Input,
) (
	model.EnrichmentResult,
	[]SearchResult,
	[]model.Evidence,
	error,
) {
	webResult,
		rawResults,
		err :=
		Search(
			input,
		)

	if err != nil {

		return result,
			nil,
			nil,
			err
	}

	if len(rawResults) == 0 {

		return result,
			nil,
			nil,
			nil
	}

	logger.Debug(
		"Raw Search Results",
		"results",
		rawResults,
	)

	logger.Debug(
		"Ranked Search Results",
		"results",
		webResult.Results,
	)

	//
	// Preserve both primary identity fields.
	//
	extractionInput :=
		model.Input{
			FullName: input.FullName,

			RootUsername: input.RootUsername,
		}

	//
	// Neutral evidence preserved for the
	// organization stage.
	//
	evidence :=
		make(
			[]model.Evidence,
			0,
			len(webResult.Results),
		)

	for _, item := range webResult.Results {

		//
		// Convert this ranked search result into
		// one stable neutral evidence record.
		//
		itemEvidence :=
			SearchEvidence(
				input,
				item,
			)

		evidence =
			append(
				evidence,
				itemEvidence,
			)

		//
		// Result URL itself is extraction
		// evidence.
		//
		// It points back to the same evidence
		// record as the search-result text.
		//
		if itemEvidence.URL != "" {

			extractionInput.URLs =
				append(
					extractionInput.URLs,
					model.URLInput{
						URL: itemEvidence.URL,

						Source: itemEvidence.Source,

						EvidenceID: itemEvidence.ID,
					},
				)
		}

		//
		// Evidence.Text is deliberately the exact
		// string passed into the text extractor.
		//
		// This guarantees that future occurrence
		// offsets point into the correct evidence
		// text.
		//
		if itemEvidence.Text != "" {

			extractionInput.Text =
				append(
					extractionInput.Text,
					model.TextInput{
						Text: itemEvidence.Text,

						Source: itemEvidence.Source,

						EvidenceID: itemEvidence.ID,
					},
				)
		}
	}

	logger.Debug(
		"Extraction Input",
		"input",
		extractionInput,
	)

	//
	// Discovery only extracts candidate facts.
	//
	// It does not determine whether those facts
	// belong to the requested target or merely
	// relate to it.
	//
	discovered :=
		extract.Extract(
			extractionInput,
		)

	logger.Debug(
		"Extraction Result",
		"result",
		discovered,
	)

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

	result.People =
		append(
			result.People,
			discovered.People...,
		)

	return result,
		rawResults,
		evidence,
		nil
}
