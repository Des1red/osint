package knowledge

import (
	"osint/internal/engines/ip/geo"
)

//
// IP — geolocation.
//

type GeoResult = geo.GeoResult

type GeoFieldResult = geo.FieldResult

type GeoCoordinateResult = geo.CoordinateResult

type GeoProviderResult = geo.ProviderResult

type GeoProviderError = geo.ProviderError
