package geo

import "strings"

func buildConsensus(
	ip string,
	providers []ProviderResult,
) GeoResult {
	result := GeoResult{
		IP: ip,

		Country: consensusField(
			providers,
			func(p ProviderResult) string {
				return p.Country
			},
		),

		CountryCode: consensusField(
			providers,
			func(p ProviderResult) string {
				return p.CountryCode
			},
		),

		Region: consensusField(
			providers,
			func(p ProviderResult) string {
				return p.Region
			},
		),

		RegionCode: consensusField(
			providers,
			func(p ProviderResult) string {
				return p.RegionCode
			},
		),

		City: consensusField(
			providers,
			func(p ProviderResult) string {
				return p.City
			},
		),

		PostalCode: consensusField(
			providers,
			func(p ProviderResult) string {
				return p.PostalCode
			},
		),

		Coordinates: consensusCoordinates(
			providers,
		),

		Timezone: consensusField(
			providers,
			func(p ProviderResult) string {
				return p.Timezone
			},
		),

		Continent: consensusField(
			providers,
			func(p ProviderResult) string {
				return normalizeContinent(
					p.Continent,
				)
			},
		),

		ContinentCode: consensusField(
			providers,
			func(p ProviderResult) string {
				return normalizeContinentCode(
					p.ContinentCode,
				)
			},
		),

		Providers: providers,
	}

	result.Reliability, result.Warning =
		geoReliability(result)

	return result
}

func consensusField(
	providers []ProviderResult,
	getValue func(ProviderResult) string,
) FieldResult {
	var values []ProviderValue

	counts := make(map[string]int)
	originalValues := make(map[string]string)

	for _, provider := range providers {
		value := strings.TrimSpace(
			getValue(provider),
		)

		if value == "" {
			continue
		}

		values = append(
			values,
			ProviderValue{
				Provider: provider.Provider,
				Value:    value,
			},
		)

		normalized := normalize(value)

		counts[normalized]++

		if _, exists := originalValues[normalized]; !exists {
			originalValues[normalized] = value
		}
	}

	total := len(values)

	if total == 0 {
		return FieldResult{
			Confidence: "unknown",
		}
	}

	if total == 1 {
		return FieldResult{
			Value:      values[0].Value,
			Confidence: "low",
			Agreement:  1,
			Total:      1,
			Providers:  values,
		}
	}

	var bestValue string
	bestCount := 0
	tied := false

	for value, count := range counts {
		if count > bestCount {
			bestValue = value
			bestCount = count
			tied = false
			continue
		}

		if count == bestCount {
			tied = true
		}
	}

	if tied {
		return FieldResult{
			Confidence: "low",
			Agreement:  bestCount,
			Total:      total,
			Providers:  values,
		}
	}

	ratio := float64(bestCount) /
		float64(total)

	return FieldResult{
		Value: originalValues[bestValue],

		Confidence: fieldConfidence(
			bestCount,
			total,
			ratio,
		),

		Agreement: bestCount,
		Total:     total,
		Providers: values,
	}
}
