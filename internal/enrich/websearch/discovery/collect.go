package discovery

import (
	"osint/internal/enrich/extract"
	"osint/internal/enrich/model"
	"osint/internal/logger"
)

type Collection struct {
	Result model.EnrichmentResult

	RankedResults []SearchResult

	RawResults []SearchResult

	Evidence []model.Evidence
}

func Collect(
	input model.Input,
) (
	Collection,
	error,
) {
	webResult,
		rawResults,
		err :=
		Search(
			input,
		)

	if err != nil {

		return Collection{},
			err
	}

	if len(rawResults) == 0 {

		return Collection{},
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

	return Collection{
		Result: discovered,

		RankedResults: webResult.Results,

		RawResults: rawResults,

		Evidence: evidence,
	}, nil
}
