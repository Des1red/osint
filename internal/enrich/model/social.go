package model

type SocialCircle struct {
	Connections []SocialCircleConnection

	Evidence []Evidence
}

type SocialCircleConnection struct {
	From string

	//
	// To is either:
	//
	//     an already-known person name
	//
	// or:
	//
	//     @username
	//
	To string

	Platform string

	Username string

	ProfileURL string

	Posts []string

	Evidence []EvidenceReference
}
