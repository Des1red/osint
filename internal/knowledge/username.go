package knowledge

import (
	"osint/internal/matcher"
	"osint/internal/platforms/variants"
)

//
// Username variants.
//

type VariantResult = variants.Result

//
// Username matching.
//

type MatchLevel = matcher.Level

const (
	MatchExact = matcher.Exact

	MatchClose = matcher.Close

	MatchBroad = matcher.Broad

	MatchNone = matcher.None
)
