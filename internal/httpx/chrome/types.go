package chrome

import "time"

const DefaultTimeout = 20 * time.Second

type Request struct {
	URL string

	Timeout time.Duration
}

type Response struct {
	HTML string

	FinalURL string
}
