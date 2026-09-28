package dedupe

import (
	dedupeinternal "osint/internal/enrich/dedupe/internal"
	"osint/internal/enrich/model"
)

func Result(
	result model.EnrichmentResult,
) model.EnrichmentResult {
	//
	// Root target.
	//
	result.Usernames =
		dedupeinternal.UsernameResult(
			result.Usernames,
		)

	result.Emails =
		dedupeinternal.EmailResult(
			result.Emails,
		)

	result.Phones =
		dedupeinternal.PhoneResult(
			result.Phones,
		)

	result.Socials =
		dedupeinternal.SocialResult(
			result.Socials,
		)

	result.Links =
		dedupeinternal.LinkResult(
			result.Links,
		)

	result.Locations =
		dedupeinternal.LocationResult(
			result.Locations,
		)

	result.Organizations =
		dedupeinternal.OrganizationResult(
			result.Organizations,
		)

	result.Employment =
		dedupeinternal.EmploymentResult(
			result.Employment,
		)

	result.Education =
		dedupeinternal.EducationResult(
			result.Education,
		)

	result.RelatedAccounts =
		dedupeinternal.AccountResult(
			result.RelatedAccounts,
		)

	//
	// Unattributed facts.
	//
	result.UnattributedUsernames =
		dedupeinternal.UsernameResult(
			result.UnattributedUsernames,
		)

	result.UnattributedEmails =
		dedupeinternal.EmailResult(
			result.UnattributedEmails,
		)

	result.UnattributedPhones =
		dedupeinternal.PhoneResult(
			result.UnattributedPhones,
		)

	result.UnattributedSocials =
		dedupeinternal.SocialResult(
			result.UnattributedSocials,
		)

	result.UnattributedLinks =
		dedupeinternal.LinkResult(
			result.UnattributedLinks,
		)

	result.UnattributedLocations =
		dedupeinternal.LocationResult(
			result.UnattributedLocations,
		)

	result.UnattributedOrganizations =
		dedupeinternal.OrganizationResult(
			result.UnattributedOrganizations,
		)

	result.UnattributedEmployment =
		dedupeinternal.EmploymentResult(
			result.UnattributedEmployment,
		)

	result.UnattributedEducation =
		dedupeinternal.EducationResult(
			result.UnattributedEducation,
		)

	result.UnattributedRelatedAccounts =
		dedupeinternal.AccountResult(
			result.UnattributedRelatedAccounts,
		)

	//
	// Related people.
	//
	// Every person now supports the same fact
	// families as the root target.
	//
	for index := range result.People {

		result.People[index].Usernames =
			dedupeinternal.UsernameResult(
				result.People[index].Usernames,
			)

		result.People[index].Emails =
			dedupeinternal.EmailResult(
				result.People[index].Emails,
			)

		result.People[index].Phones =
			dedupeinternal.PhoneResult(
				result.People[index].Phones,
			)

		result.People[index].Socials =
			dedupeinternal.SocialResult(
				result.People[index].Socials,
			)

		result.People[index].Links =
			dedupeinternal.LinkResult(
				result.People[index].Links,
			)

		result.People[index].Locations =
			dedupeinternal.LocationResult(
				result.People[index].Locations,
			)

		result.People[index].Organizations =
			dedupeinternal.OrganizationResult(
				result.People[index].Organizations,
			)

		result.People[index].Employment =
			dedupeinternal.EmploymentResult(
				result.People[index].Employment,
			)

		result.People[index].Education =
			dedupeinternal.EducationResult(
				result.People[index].Education,
			)

		result.People[index].RelatedAccounts =
			dedupeinternal.AccountResult(
				result.People[index].RelatedAccounts,
			)
	}

	return result
}
