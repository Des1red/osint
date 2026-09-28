package accounts

import (
	"strings"

	"osint/internal/enrich/model"
)

func Organize(
	result model.EnrichmentResult,
	input model.Input,
	evidence []model.Evidence,
) model.EnrichmentResult {
	fullName :=
		strings.TrimSpace(
			input.FullName,
		)

	if fullName == "" {

		return result
	}

	//
	// Root target.
	//
	targetEvidence :=
		evidenceForIdentity(
			fullName,
			evidence,
		)

	result.Usernames,
		result.Socials,
		result.RelatedAccounts =
		organizeIdentity(
			result.Usernames,
			result.Socials,
			result.RelatedAccounts,
			input,
			targetEvidence,
		)

	//
	// Related people.
	//
	for index := range result.People {

		personName :=
			strings.TrimSpace(
				result.People[index].Name,
			)

		if personName == "" {

			continue
		}

		personInput :=
			model.Input{
				FullName: personName,
			}

		personEvidence :=
			evidenceForIdentity(
				personName,
				evidence,
			)

		result.People[index].Usernames,
			result.People[index].Socials,
			result.People[index].RelatedAccounts =
			organizeIdentity(
				result.People[index].Usernames,
				result.People[index].Socials,
				result.People[index].RelatedAccounts,
				personInput,
				personEvidence,
			)
	}

	return result
}

func organizeIdentity(
	usernames []model.UsernameReference,
	socials []model.SocialReference,
	related []model.RelatedAccountReference,
	input model.Input,
	evidence []model.Evidence,
) (
	[]model.UsernameReference,
	[]model.SocialReference,
	[]model.RelatedAccountReference,
) {
	owned :=
		newOwnedAccounts(
			input,
		)

	for _, item := range evidence {

		platform,
			username,
			ok :=
			ownedProfileFromEvidence(
				input,
				item,
			)

		if !ok {

			continue
		}

		owned.add(
			platform,
			username,
		)
	}

	related =
		append(
			related,
			relatedAccountsFromEvidence(
				evidence,
				owned,
			)...,
		)

	var ownedUsernames []model.UsernameReference

	for _, username := range usernames {

		key :=
			normalizeUsername(
				username.Username,
			)

		if key == "" {

			continue
		}

		if owned.hasUsername(
			key,
		) {

			ownedUsernames =
				append(
					ownedUsernames,
					username,
				)

			continue
		}

		//
		// Being present somewhere inside
		// subject evidence is not sufficient to
		// establish a related account.
		//
		// Example:
		//
		//     Person2 appears on a Piknu following
		//     page containing hundreds of handles.
		//
		// Those handles are page contents, not
		// meaningful Person2-related accounts.
		//
		if !usernameAssociationAllowed(
			username,
			evidence,
		) {

			continue
		}

		related =
			append(
				related,
				model.RelatedAccountReference{
					Username: username.Username,

					Association: "Mentioned",

					Source: username.Source,

					Evidence: username.Evidence,
				},
			)
	}

	var ownedSocials []model.SocialReference

	for _, social := range socials {

		platform :=
			normalizePlatform(
				social.Platform,
			)

		username :=
			normalizeUsername(
				social.Username,
			)

		if username == "" {

			continue
		}

		if owned.hasSocial(
			platform,
			username,
		) {

			ownedSocials =
				append(
					ownedSocials,
					social,
				)

			continue
		}

		if !socialAssociationAllowed(
			social,
			evidence,
		) {

			continue
		}

		related =
			append(
				related,
				model.RelatedAccountReference{
					Platform: social.Platform,

					Username: social.Username,

					ProfileURL: social.URL,

					Association: "Referenced",

					Source: social.Source,

					Evidence: social.Evidence,
				},
			)
	}

	return ownedUsernames,
		ownedSocials,
		related
}
