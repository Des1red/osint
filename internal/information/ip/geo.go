package ip

import (
	"fmt"
	"strings"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func geoNodes(
	result knowledge.GeoResult,
) []reportNode {
	nodes :=
		[]reportNode{
			reportSection(
				"Summary",
				reportValue(
					"IP",
					result.IP,
				),
				reportValue(
					"Location Reliability",
					result.Reliability,
				),
				reportValue(
					"Warning",
					result.Warning,
				),
			),

			reportSection(
				"Location",
				geoFieldNode(
					"Country",
					result.Country,
				),
				geoFieldNode(
					"Country Code",
					result.CountryCode,
				),
				geoFieldNode(
					"Region",
					result.Region,
				),
				geoFieldNode(
					"Region Code",
					result.RegionCode,
				),
				geoFieldNode(
					"City",
					result.City,
				),
				geoFieldNode(
					"Postal Code",
					result.PostalCode,
				),
			),

			geoCoordinatesNode(
				result.Coordinates,
			),

			reportSection(
				"Regional Metadata",
				geoFieldNode(
					"Timezone",
					result.Timezone,
				),
				geoFieldNode(
					"Continent",
					result.Continent,
				),
				geoFieldNode(
					"Continent Code",
					result.ContinentCode,
				),
			),
		}

	if output.Full() {
		nodes =
			append(
				nodes,
				geoProvidersNode(
					result.Providers,
				),
				geoErrorsNode(
					result.Errors,
				),
			)
	}

	return compactReportNodes(
		nodes,
	)
}

func geoFieldNode(
	name string,
	field knowledge.GeoFieldResult,
) reportNode {
	if field.Total == 0 {
		return reportNode{}
	}

	value :=
		strings.TrimSpace(
			field.Value,
		)

	if value == "" {
		value =
			"uncertain"
	}

	label :=
		fmt.Sprintf(
			"%s: %s",
			name,
			value,
		)

	if strings.TrimSpace(
		field.Confidence,
	) != "" {
		label +=
			" [" +
				field.Confidence +
				"]"
	}

	node :=
		reportNode{
			Label: label,
		}

	if !output.Full() ||
		len(field.Providers) == 0 {
		return node
	}

	var providers []reportNode

	for _, provider := range field.Providers {
		providerValue :=
			strings.TrimSpace(
				provider.Value,
			)

		if providerValue == "" {
			providerValue =
				"uncertain"
		}

		providers =
			append(
				providers,
				reportValue(
					provider.Provider,
					providerValue,
				),
			)
	}

	node.Children =
		[]reportNode{
			reportSection(
				"Provider Values",
				providers...,
			),
		}

	return node
}

func geoCoordinatesNode(
	result knowledge.GeoCoordinateResult,
) reportNode {
	if len(result.Providers) == 0 {
		return reportNode{}
	}

	var children []reportNode

	if result.Available {
		children =
			append(
				children,
				reportValue(
					"Consensus",
					fmt.Sprintf(
						"%.6f, %.6f [%s]",
						result.Latitude,
						result.Longitude,
						result.Confidence,
					),
				),
			)
	} else {
		children =
			append(
				children,
				reportValue(
					"Consensus",
					"uncertain ["+
						result.Confidence+
						"]",
				),
			)
	}

	if len(result.Providers) > 1 {
		children =
			append(
				children,
				reportValue(
					"Maximum Spread",
					fmt.Sprintf(
						"%.2f km",
						result.SpreadKM,
					),
				),
			)
	}

	if output.Full() {
		var providers []reportNode

		for _, provider := range result.Providers {
			providers =
				append(
					providers,
					reportNode{
						Label: provider.Provider,

						Children: []reportNode{
							reportValue(
								"Latitude",
								fmt.Sprintf(
									"%.6f",
									provider.Latitude,
								),
							),
							reportValue(
								"Longitude",
								fmt.Sprintf(
									"%.6f",
									provider.Longitude,
								),
							),
						},
					},
				)
		}

		children =
			append(
				children,
				reportSection(
					"Provider Coordinates",
					providers...,
				),
			)
	}

	return reportSection(
		"Coordinates",
		children...,
	)
}

func geoProvidersNode(
	providers []knowledge.GeoProviderResult,
) reportNode {
	var nodes []reportNode

	for _, provider := range providers {
		nodes =
			append(
				nodes,
				reportNode{
					Label: provider.Provider,

					Children: []reportNode{
						reportValue(
							"Source",
							provider.Source,
						),
					},
				},
			)
	}

	return reportSection(
		"Providers",
		nodes...,
	)
}

func geoErrorsNode(
	errors []knowledge.GeoProviderError,
) reportNode {
	var nodes []reportNode

	for _, providerErr := range errors {
		nodes =
			append(
				nodes,
				reportNode{
					Label: providerErr.Provider,

					Children: []reportNode{
						reportValue(
							"Error",
							providerErr.Error,
						),
					},
				},
			)
	}

	return reportSection(
		"Provider Errors",
		nodes...,
	)
}
