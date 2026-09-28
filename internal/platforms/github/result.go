package github

type GitHubResult struct {
	Found bool

	Username  string
	ID        int64
	NodeID    string
	Type      string
	SiteAdmin bool

	Name     string
	Company  string
	Blog     string
	Location string
	Email    string
	Bio      string
	Twitter  string

	PublicRepos int
	PublicGists int

	Followers int
	Following int

	CreatedAt string
	UpdatedAt string

	AvatarURL  string
	ProfileURL string
	APIURL     string

	Source string
}
