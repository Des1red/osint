package username

import (
	"osint/internal/knowledge"
)

func matchLabel(
	level knowledge.MatchLevel,
) string {
	switch level {

	case knowledge.MatchExact:
		return "Exact"

	case knowledge.MatchClose:
		return "Close"

	case knowledge.MatchBroad:
		return "Broad"

	default:
		return "None"
	}
}
