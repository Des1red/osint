package discovery

import "osint/internal/enrich/model"

func Enrich(
	result model.EnrichmentResult,
	input model.Input,
) (
	model.EnrichmentResult,
	[]SearchResult,
	[]model.Evidence,
	error,
) {
	collection, err :=
		collect(
			input,
		)

	if err != nil {

		return result,
			nil,
			nil,
			err
	}

	result =
		mergeCollection(
			result,
			collection.Result,
		)

	return result,
		collection.RawResults,
		collection.Evidence,
		nil
}

func mergeCollection(
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

	result.People =
		append(
			result.People,
			discovered.People...,
		)

	return result
}
