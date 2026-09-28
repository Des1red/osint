package geo

import "strings"

func normalize(
	value string,
) string {
	return strings.ToLower(
		strings.TrimSpace(value),
	)
}

func normalizeContinent(
	value string,
) string {
	switch normalize(value) {
	case "north america":
		return "North America"

	case "south america":
		return "South America"

	case "europe":
		return "Europe"

	case "asia":
		return "Asia"

	case "africa":
		return "Africa"

	case "oceania":
		return "Oceania"

	case "antarctica":
		return "Antarctica"

	default:
		return strings.TrimSpace(value)
	}
}

func normalizeContinentCode(
	value string,
) string {
	switch strings.ToUpper(
		strings.TrimSpace(value),
	) {
	case "NA":
		return "NA"

	case "SA":
		return "SA"

	case "EU":
		return "EU"

	case "AS":
		return "AS"

	case "AF":
		return "AF"

	case "OC":
		return "OC"

	case "AN":
		return "AN"

	default:
		return ""
	}
}
