package tiktok

type TikTokResult struct {
	Found         bool
	Accessible    bool
	LoginRequired bool

	Username string
	Name     string

	Description string

	ProfileURL string
	Source     string
}
