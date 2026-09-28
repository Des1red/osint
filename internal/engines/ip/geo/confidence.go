package geo

func fieldConfidence(
	bestCount int,
	total int,
	ratio float64,
) string {
	if total < 2 {
		return "low"
	}

	if bestCount == total {
		return "high"
	}

	if ratio > 0.5 {
		return "medium"
	}

	return "low"
}

func coordinateConfidence(
	spreadKM float64,
) string {
	switch {
	case spreadKM <= 10:
		return "high"

	case spreadKM <= 50:
		return "medium"

	default:
		return "low"
	}
}

func geoReliability(
	result GeoResult,
) (string, string) {
	if len(result.Providers) == 1 {
		return "low",
			"only one geolocation provider returned data"
	}

	if result.Coordinates.SpreadKM > 1000 {
		return "very low",
			"providers report geographically distant locations; IP may belong to distributed or anycast infrastructure"
	}

	if result.Coordinates.SpreadKM > 100 {
		return "low",
			"providers significantly disagree on location"
	}

	if result.City.Confidence == "high" &&
		result.Region.Confidence == "high" &&
		result.Country.Confidence == "high" {
		return "high", ""
	}

	if result.Country.Confidence == "high" {
		return "medium", ""
	}

	return "low",
		"geolocation providers disagree"
}
