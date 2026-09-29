package model

type EnrichmentResult struct {
	//
	// Facts attributed to the requested root
	// identity.
	//
	Usernames []UsernameReference

	Emails []EmailReference

	Phones []PhoneReference

	Socials []SocialReference

	Links []ExternalLink

	Locations []LocationReference

	Organizations []OrganizationReference

	Employment []EmploymentReference

	Education []EducationReference

	RelatedAccounts []RelatedAccountReference

	//
	// Facts for which semantic attribution could
	// not be established strongly enough.
	//
	UnattributedUsernames []UsernameReference

	UnattributedEmails []EmailReference

	UnattributedPhones []PhoneReference

	UnattributedSocials []SocialReference

	UnattributedLinks []ExternalLink

	UnattributedLocations []LocationReference

	UnattributedOrganizations []OrganizationReference

	UnattributedEmployment []EmploymentReference

	UnattributedEducation []EducationReference

	UnattributedRelatedAccounts []RelatedAccountReference

	//
	// Other known people in the current
	// investigation graph.
	//
	People []PersonReference
}
