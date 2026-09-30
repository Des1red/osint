package opencorp

type Result struct {
	Officers []Officer
}

type Officer struct {
	Name string

	Position string

	Company string

	Jurisdiction string

	StartDate string

	EndDate string

	ProfileURL string

	CompanyURL string
}
