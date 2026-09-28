package knowledge

import (
	"osint/internal/engines/ip/network"
)

//
// IP — network.
//

type NetworkResult = network.NetworkResult

type NetworkASNResult = network.ASNResult

type NetworkASNNeighbourCounts = network.ASNNeighbourCounts

type NetworkASNNeighbour = network.ASNNeighbour

type NetworkFlagResult = network.FlagResult

type NetworkAnycastResult = network.AnycastResult

type NetworkProviderError = network.ProviderError
