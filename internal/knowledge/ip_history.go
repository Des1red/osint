package knowledge

import (
	"osint/internal/engines/ip/history"
)

//
// IP — history.
//

type HistoryResult = history.HistoryResult

type HistorySeenResult = history.SeenResult

type HistoryRoutingResult = history.RoutingHistoryResult

type HistoryRoutingOrigin = history.RoutingOrigin

type HistoryRoutingPrefix = history.RoutingPrefix

type HistoryRoutingTimeline = history.RoutingTimeline

type HistoryProviderError = history.ProviderError

type HistoryRoutingStatusResult = history.RoutingStatusResult
