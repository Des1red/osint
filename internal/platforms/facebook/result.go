package facebook

type FacebookResult struct {
	Found         bool
	Accessible    bool
	LoginRequired bool

	Username string
	Name     string

	Description string
	Category    string

	ProfileURL string

	Source string
}
