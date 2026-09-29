package model

// Search evidence collected during the
// second-hop WebSearch for an associated
// person.
type RelativeSearchResult struct {
	Query string

	Title string

	URL string

	Snippet string

	Engine string

	EvidenceID string
}

type PersonReference struct {
	Name string

	//
	// Evidence which originally established
	// the existence of this person.
	//
	Source string

	Evidence []EvidenceReference

	//
	// Facts associated with this person after
	// semantic organization.
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
	// Family relationship relative to the
	// original search target.
	//
	Relation string

	RelationVerified bool

	//
	// Separate evidence that explicitly
	// established the family relationship.
	//
	RelationSource string

	RelationEvidence []EvidenceReference

	//
	// Second-hop WebSearch results for this
	// person.
	//
	SearchResults []RelativeSearchResult
}
