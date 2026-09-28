package username

import (
	"osint/internal/enrich"
	"osint/internal/platforms/variants"
)

type UsernameResult struct {
	variants.PlatformResults

	Enrichment enrich.EnrichmentResult

	Variants []variants.Result
}
