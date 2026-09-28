package extract

import (
	"strings"

	"osint/internal/enrich/model"
	"osint/internal/matcher"
)

func Extract(
	input model.Input,
) model.EnrichmentResult {
	var result model.EnrichmentResult

	result.Emails =
		extractEmails(
			input,
		)

	result.Phones =
		extractPhones(
			input,
		)

	result.Usernames =
		extractUsernames(
			input,
		)

	result.Locations =
		extractLocations(
			input,
		)

	result.Organizations =
		extractOrganizations(
			input,
		)

	result.Employment =
		extractEmployment(
			input,
		)

	result.Education =
		extractEducation(
			input,
		)

	result.People =
		extractPeople(
			input,
		)

	relationshipOrganizations,
		relationshipEmployment :=
		extractRelationships(
			input,
		)

	result.Organizations =
		append(
			result.Organizations,
			relationshipOrganizations...,
		)

	result.Employment =
		append(
			result.Employment,
			relationshipEmployment...,
		)

	links,
		socials :=
		extractLinks(
			input,
		)

	result.Links =
		links

	result.Socials =
		socials

	classifySocials(
		&result,
		input.RootUsername,
	)

	result.Emails =
		uniqueEmails(
			result.Emails,
		)

	result.Phones =
		uniquePhones(
			result.Phones,
		)

	result.Usernames =
		uniqueUsernames(
			result.Usernames,
		)

	result.Locations =
		uniqueLocations(
			result.Locations,
		)

	result.Organizations =
		uniqueOrganizations(
			result.Organizations,
		)

	result.Employment =
		uniqueEmployment(
			result.Employment,
		)

	result.Education =
		uniqueEducation(
			result.Education,
		)

	result.People =
		uniquePeople(
			result.People,
		)

	result.Links =
		uniqueLinks(
			result.Links,
		)

	result.Socials =
		uniqueSocials(
			result.Socials,
		)

	return result
}

func classifySocials(
	result *model.EnrichmentResult,
	rootUsername string,
) {
	rootUsername =
		strings.TrimSpace(
			rootUsername,
		)

	for index := range result.Socials {

		if rootUsername == "" {

			result.Socials[index].Match =
				matcher.None

			continue
		}

		result.Socials[index].Match =
			matcher.Classify(
				rootUsername,
				result.Socials[index].Username,
			)
	}

	for _, social := range result.Socials {

		if strings.TrimSpace(
			social.Username,
		) == "" {

			continue
		}

		//
		// A numeric platform profile ID is useful
		// as part of the SocialReference, but it
		// should not be presented as a generic
		// username.
		//
		if !genericUsernameCandidate(
			social.Username,
		) {

			continue
		}

		if rootUsername != "" &&
			social.Match ==
				matcher.Exact {

			continue
		}

		result.Usernames =
			append(
				result.Usernames,
				model.UsernameReference{
					Username: social.Username,

					Source: social.Source,

					Match: social.Match,

					Evidence: social.Evidence,
				},
			)
	}
}
