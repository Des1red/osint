package network

type ASNNeighbourCounts struct {
	Left      int
	Right     int
	Uncertain int
	Unique    int
}

type ASNNeighbour struct {
	ASN      uint32
	Position string

	PathCount int

	PeerCountV4 int
	PeerCountV6 int
}

type ASNResult struct {
	ASN       uint32
	Holder    string
	Announced bool

	AnnouncedPrefixes []string

	NeighbourCounts ASNNeighbourCounts
	Neighbours      []ASNNeighbour
}

type FlagResult struct {
	Available bool
	Value     bool
	Source    string
}

type AnycastLocation struct {
	ID string

	City    string
	Country string

	Latitude  float64
	Longitude float64
}

type AnycastResult struct {
	Available  bool
	Anycast    bool
	Confidence string

	Prefix        string
	BackingPrefix string
	MappedPrefix  string

	ASNs []uint32

	ABICMP int
	ABTCP  int
	ABDNS  int

	GCDICMP int
	GCDTCP  int

	Locations []AnycastLocation

	Date   string
	Source string
}

type ProviderError struct {
	Provider string
	Error    string
}

type NetworkResult struct {
	IP string

	Prefix string

	ASNs []ASNResult

	ReverseDNS []string

	ISP            string
	Organization   string
	Domain         string
	ConnectionType string

	Mobile  FlagResult
	Hosting FlagResult
	Proxy   FlagResult
	VPN     FlagResult
	Tor     FlagResult

	Anycast AnycastResult

	Threat string

	Sources []string
	Errors  []ProviderError
}
