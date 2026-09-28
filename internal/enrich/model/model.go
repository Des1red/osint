package model

import (
	"osint/internal/matcher"
)

// Evidence anchor.
//
// The anchor describes the identity or entity
// whose investigation produced the evidence.
//
// It does NOT mean that every fact inside the
// evidence belongs to that entity.
type EvidenceAnchorKind string

const (
	EvidenceAnchorUnknown EvidenceAnchorKind = ""

	EvidenceAnchorPerson EvidenceAnchorKind = "person"

	EvidenceAnchorAccount EvidenceAnchorKind = "account"

	EvidenceAnchorOrganization EvidenceAnchorKind = "organization"

	EvidenceAnchorDomain EvidenceAnchorKind = "domain"
)

type EvidenceAnchor struct {
	Kind EvidenceAnchorKind

	Value string
}

// Evidence subject.
//
// A subject describes an entity that the
// evidence itself was independently matched
// against.
//
// Example:
//
// Search:
//
//	Person2
//
// Result page:
//
//	Contact: Person1
//	City: MILOS
//
// Anchor:
//
//	Person2
//
// Subject:
//
//	Person1
//
// Subject still does not automatically mean
// that every individual fact is owned by that
// entity.
//
// Type-specific organization can use stronger
// evidence when available.
//
// Example:
//
//	Person1 +30...
//	Person2 +44...
//
// The phone organizer can still resolve each
// phone from its exact occurrence.
type EvidenceSubject struct {
	Kind EvidenceAnchorKind

	Value string
}

// Reference back to the evidence from which a
// fact was extracted.
//
// Start and End are byte offsets inside
// Evidence.Text.
//
// When Start == End, the reference identifies
// the evidence as a whole rather than an exact
// text occurrence.
type EvidenceReference struct {
	EvidenceID string

	Start int

	End int
}

type AccountInput struct {
	Platform string

	Username string

	ProfileURL string

	EvidenceID string
}

type UsernameInput struct {
	Username string

	Source string

	EvidenceID string
}

type TextInput struct {
	Text string

	Source string

	EvidenceID string
}

type URLInput struct {
	URL string

	Source string

	EvidenceID string
}

type OrganizationInput struct {
	Organization string

	Source string

	EvidenceID string
}

type LocationInput struct {
	Location string

	Source string

	EvidenceID string
}

type EmploymentInput struct {
	Title string

	Organization string

	StartDate string

	EndDate string

	Current bool

	Summary string

	Source string

	EvidenceID string
}

type EducationInput struct {
	School string

	Degrees []string

	Majors []string

	StartDate string

	EndDate string

	Summary string

	Source string

	EvidenceID string
}

type Input struct {
	FullName string

	RootUsername string

	Accounts []AccountInput

	Usernames []UsernameInput

	Text []TextInput

	URLs []URLInput

	Organizations []OrganizationInput

	Locations []LocationInput

	Employment []EmploymentInput

	Education []EducationInput
}

// Neutral evidence passed between enrichment
// stages.
//
// Anchor:
//
//	why / through whom we discovered it.
//
// Subjects:
//
//	which known entities the evidence itself
//	independently matches.
//
// Neither field by itself means that every fact
// inside the evidence is owned by that entity.
type Evidence struct {
	ID string

	Anchor EvidenceAnchor

	Subjects []EvidenceSubject

	Source string

	Query string

	Engine string

	Title string

	//
	// Text contains the exact text supplied to
	// text-based extractors.
	//
	// EvidenceReference offsets point into this
	// value.
	//
	Text string

	URL string
}

type UsernameReference struct {
	Username string

	Source string

	Match matcher.Level

	Evidence []EvidenceReference
}

type EmailReference struct {
	Email string

	Type string

	Source string

	Evidence []EvidenceReference
}

type PhoneReference struct {
	Phone string

	Source string

	Evidence []EvidenceReference
}

type SocialReference struct {
	Platform string

	Username string

	URL string

	Source string

	Match matcher.Level

	Evidence []EvidenceReference
}

type ExternalLink struct {
	URL string

	ResolvedURL string

	Category string

	Source string

	Evidence []EvidenceReference
}

type LocationReference struct {
	Location string

	Source string

	Evidence []EvidenceReference
}

type OrganizationReference struct {
	Organization string

	Source string

	Evidence []EvidenceReference
}

type EmploymentReference struct {
	Title string

	Organization string

	StartDate string

	EndDate string

	Current bool

	Summary string

	Source string

	Evidence []EvidenceReference
}

type EducationReference struct {
	School string

	Degrees []string

	Majors []string

	StartDate string

	EndDate string

	Summary string

	Source string

	Evidence []EvidenceReference
}

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

type RelatedAccountReference struct {
	Platform string

	Username string

	ProfileURL string

	EvidenceURL string

	Association string

	Source string

	Evidence []EvidenceReference
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
