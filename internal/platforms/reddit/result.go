package reddit

type RedditResult struct {
	Found      bool
	Accessible bool

	Username string

	PostKarma    int
	CommentKarma int
	TotalKarma   int

	RedditAge string
	CakeDay   string

	ProfileDescription string

	ProfileURL string
	Source     string
}
