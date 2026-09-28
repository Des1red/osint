package knowledge

import (
	"osint/internal/engines/ip/reputation"
)

//
// IP — reputation.
//

type ReputationResult = reputation.ReputationResult

type AbuseIPDBResult = reputation.AbuseIPDBResult

type VirusTotalResult = reputation.VirusTotalResult

type IPQSResult = reputation.IPQSResult

type ReputationProviderError = reputation.ProviderError
