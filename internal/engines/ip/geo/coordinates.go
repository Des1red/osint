package geo

import "math"

func consensusCoordinates(
	providers []ProviderResult,
) CoordinateResult {
	var values []ProviderCoordinates

	for _, provider := range providers {
		if !hasCoordinates(provider) {
			continue
		}

		values = append(
			values,
			ProviderCoordinates{
				Provider:  provider.Provider,
				Latitude:  provider.Latitude,
				Longitude: provider.Longitude,
			},
		)
	}

	if len(values) == 0 {
		return CoordinateResult{
			Confidence: "unknown",
		}
	}

	if len(values) == 1 {
		return CoordinateResult{
			Latitude:   values[0].Latitude,
			Longitude:  values[0].Longitude,
			Available:  true,
			Confidence: "low",
			Providers:  values,
		}
	}

	maxSpread := maxCoordinateSpread(
		values,
	)

	confidence := coordinateConfidence(
		maxSpread,
	)

	result := CoordinateResult{
		Confidence: confidence,
		SpreadKM:   maxSpread,
		Providers:  values,
	}

	if confidence == "low" {
		return result
	}

	latitude, longitude := averageCoordinates(
		values,
	)

	result.Latitude = latitude
	result.Longitude = longitude
	result.Available = true

	return result
}

func hasCoordinates(
	provider ProviderResult,
) bool {
	return provider.Latitude != 0 ||
		provider.Longitude != 0
}

func maxCoordinateSpread(
	values []ProviderCoordinates,
) float64 {
	var maximum float64

	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			distance := haversine(
				values[i].Latitude,
				values[i].Longitude,
				values[j].Latitude,
				values[j].Longitude,
			)

			if distance > maximum {
				maximum = distance
			}
		}
	}

	return maximum
}

func averageCoordinates(
	values []ProviderCoordinates,
) (float64, float64) {
	var latitude float64
	var longitude float64

	for _, value := range values {
		latitude += value.Latitude
		longitude += value.Longitude
	}

	count := float64(
		len(values),
	)

	return latitude / count,
		longitude / count
}

func haversine(
	lat1 float64,
	lon1 float64,
	lat2 float64,
	lon2 float64,
) float64 {
	const earthRadiusKM = 6371.0

	lat1Rad := degreesToRadians(lat1)
	lat2Rad := degreesToRadians(lat2)

	deltaLat := degreesToRadians(
		lat2 - lat1,
	)

	deltaLon := degreesToRadians(
		lon2 - lon1,
	)

	a := math.Sin(deltaLat/2)*
		math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*
			math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*
			math.Sin(deltaLon/2)

	c := 2 * math.Atan2(
		math.Sqrt(a),
		math.Sqrt(1-a),
	)

	return earthRadiusKM * c
}

func degreesToRadians(
	degrees float64,
) float64 {
	return degrees * math.Pi / 180
}
