package geo

type ProviderValue struct {
	Provider string
	Value    string
}

type FieldResult struct {
	Value      string
	Confidence string
	Agreement  int
	Total      int
	Providers  []ProviderValue
}

type ProviderCoordinates struct {
	Provider  string
	Latitude  float64
	Longitude float64
}

type CoordinateResult struct {
	Latitude   float64
	Longitude  float64
	Available  bool
	Confidence string
	SpreadKM   float64

	Providers []ProviderCoordinates
}

type ProviderError struct {
	Provider string
	Error    string
}

type ProviderResult struct {
	Provider string

	IP string

	Country     string
	CountryCode string

	Region     string
	RegionCode string

	City       string
	PostalCode string

	Latitude  float64
	Longitude float64

	Timezone string

	Continent     string
	ContinentCode string

	Source string
}

type GeoResult struct {
	IP string

	Country     FieldResult
	CountryCode FieldResult

	Region     FieldResult
	RegionCode FieldResult

	City       FieldResult
	PostalCode FieldResult

	Coordinates CoordinateResult

	Timezone FieldResult

	Continent     FieldResult
	ContinentCode FieldResult

	Reliability string
	Warning     string

	Providers []ProviderResult
	Errors    []ProviderError
}
