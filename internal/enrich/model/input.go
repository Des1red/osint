package model

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
	Surname  string

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
