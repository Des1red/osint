package discovery

import (
	"osint/internal/httpx/chrome/searchengines"
)

type Result struct {
	Results []SearchResult
}

type SearchResult = searchengines.Result
