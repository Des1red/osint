package ip

import (
	"osint/internal/engines/ip/geo"
	"osint/internal/engines/ip/history"
	"osint/internal/engines/ip/network"
	"osint/internal/engines/ip/rdap"
	"osint/internal/engines/ip/reputation"
)

type IPResult struct {
	RDAP       rdap.RDAPResult
	Geo        geo.GeoResult
	Network    network.NetworkResult
	History    history.HistoryResult
	Reputation reputation.ReputationResult
}
