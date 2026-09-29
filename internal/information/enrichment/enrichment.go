package enrichment

import (
	"fmt"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func Print(
	result knowledge.EnrichmentResult,
) bool {
	if !hasData(
		result,
	) {

		return false
	}

	fmt.Fprintln(
		output.Writer(),
	)

	fmt.Fprintln(
		output.Writer(),
		"Discovered Information",
	)

	fmt.Fprintln(
		output.Writer(),
		"======================",
	)

	printIdentity(
		result,
	)

	printContactAndLocation(
		result,
	)

	printProfessional(
		result,
	)

	printPeople(
		result.People,
	)

	printSocialCircle(
		result,
	)

	printReferences(
		result,
	)

	return true
}

func hasData(
	result knowledge.EnrichmentResult,
) bool {
	return len(result.Usernames) > 0 ||
		len(result.UnattributedUsernames) > 0 ||
		len(result.RelatedAccounts) > 0 ||
		len(result.UnattributedRelatedAccounts) > 0 ||
		len(result.Emails) > 0 ||
		len(result.UnattributedEmails) > 0 ||
		len(result.Phones) > 0 ||
		len(result.UnattributedPhones) > 0 ||
		len(result.Socials) > 0 ||
		len(result.UnattributedSocials) > 0 ||
		len(result.Links) > 0 ||
		len(result.UnattributedLinks) > 0 ||
		len(result.Locations) > 0 ||
		len(result.UnattributedLocations) > 0 ||
		len(result.Organizations) > 0 ||
		len(result.UnattributedOrganizations) > 0 ||
		len(result.Employment) > 0 ||
		len(result.UnattributedEmployment) > 0 ||
		len(result.Education) > 0 ||
		len(result.UnattributedEducation) > 0 ||
		len(result.People) > 0 ||
		len(result.SocialCircle.Connections) > 0
}
