package model

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
