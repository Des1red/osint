package twitter

type TwitterResult struct {
	Found         bool
	Accessible    bool
	LoginRequired bool

	Username string
	Name     string

	Description string

	ProfileURL string
	Source     string
}
