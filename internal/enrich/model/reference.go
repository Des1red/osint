package model

import (
	"osint/internal/matcher"
)

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

type RelatedAccountReference struct {
	Platform string

	Username string

	ProfileURL string

	EvidenceURL string

	Association string

	Source string

	Evidence []EvidenceReference
}
