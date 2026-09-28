package history

type ProviderError struct {
	Provider string

	Error string
}

type SeenResult struct {
	Available bool

	Time string

	Origin string

	Prefix string
}

type RoutingTimeline struct {
	StartTime string

	EndTime string

	FullPeersSeeing int

	Visibility float64

	VisibilityAvailable bool
}

type RoutingPrefix struct {
	Prefix string

	Timelines []RoutingTimeline
}

type RoutingOrigin struct {
	Origin string

	Prefixes []RoutingPrefix
}

type RoutingHistoryResult struct {
	Available bool

	Resource string

	QueryStartTime string

	QueryEndTime string

	Origins []RoutingOrigin

	Source string
}

type RoutingStatusResult struct {
	Available bool

	Resource string

	QueryTime string

	FirstSeen SeenResult

	LastSeen SeenResult

	Source string
}

type HistoryResult struct {
	IP string

	Status RoutingStatusResult

	Routing RoutingHistoryResult

	Sources []string

	Errors []ProviderError
}
