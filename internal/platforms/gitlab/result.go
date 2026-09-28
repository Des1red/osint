package gitlab

type GitLabResult struct {
	Found bool

	ID       int64
	Username string
	Name     string
	State    string
	Locked   bool

	AvatarURL  string
	ProfileURL string

	Source string
}
