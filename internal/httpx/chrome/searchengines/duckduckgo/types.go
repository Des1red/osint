package duckduckgo

type Result struct {
	Query string

	Title string

	URL string

	Snippet string

	Engine string
}

type pageResponse struct {
	Query string

	HTML string

	FinalURL string

	Err error
}
