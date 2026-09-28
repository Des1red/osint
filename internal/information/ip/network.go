package ip

import (
	"fmt"
	"strings"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

const maxCompactNetworkItems = 10

func networkNodes(
	result knowledge.NetworkResult,
) []reportNode {
	nodes :=
		[]reportNode{
			networkSummaryNode(
				result,
			),
			networkProviderNode(
				result,
			),
			networkRoutingNode(
				result.ASNs,
			),
			reverseDNSNode(
				result.ReverseDNS,
			),
			networkClassificationNode(
				result,
			),
			anycastNode(
				result.Anycast,
			),
		}

	if output.Full() {
		nodes =
			append(
				nodes,
				networkSourcesNode(
					result.Sources,
				),
				networkErrorsNode(
					result.Errors,
				),
			)
	}

	return compactReportNodes(
		nodes,
	)
}

func networkSummaryNode(
	result knowledge.NetworkResult,
) reportNode {
	return reportSection(
		"Summary",
		reportValue(
			"IP",
			result.IP,
		),
		reportValue(
			"Routed Prefix",
			result.Prefix,
		),
	)
}

func networkProviderNode(
	result knowledge.NetworkResult,
) reportNode {
	return reportSection(
		"Provider",
		reportValue(
			"ISP",
			result.ISP,
		),
		reportValue(
			"Organization",
			result.Organization,
		),
		reportValue(
			"Domain",
			result.Domain,
		),
	)
}

func networkRoutingNode(
	asns []knowledge.NetworkASNResult,
) reportNode {
	if len(asns) == 0 {
		return reportNode{}
	}

	var nodes []reportNode

	for _, asn := range asns {
		label :=
			fmt.Sprintf(
				"AS%d",
				asn.ASN,
			)

		if strings.TrimSpace(
			asn.Holder,
		) != "" {
			label +=
				" — " +
					strings.TrimSpace(
						asn.Holder,
					)
		}

		nodes =
			append(
				nodes,
				reportNode{
					Label: label,

					Children: compactReportNodes(
						[]reportNode{
							reportBool(
								"Announced",
								asn.Announced,
							),
							announcedPrefixesNode(
								asn.AnnouncedPrefixes,
							),
							asnNeighboursNode(
								asn.NeighbourCounts,
								asn.Neighbours,
							),
						},
					),
				},
			)
	}

	return reportSection(
		"Routing",
		reportSection(
			"Autonomous Systems",
			nodes...,
		),
	)
}

func announcedPrefixesNode(
	prefixes []string,
) reportNode {
	if len(prefixes) == 0 {
		return reportNode{}
	}

	limit :=
		len(prefixes)

	if !output.Full() &&
		limit >
			maxCompactNetworkItems {
		limit =
			maxCompactNetworkItems
	}

	var nodes []reportNode

	for index := 0; index < limit; index++ {
		nodes =
			append(
				nodes,
				reportItem(
					prefixes[index],
				),
			)
	}

	if limit < len(prefixes) {
		nodes =
			append(
				nodes,
				reportItem(
					fmt.Sprintf(
						"... %d more",
						len(prefixes)-limit,
					),
				),
			)
	}

	return reportSection(
		fmt.Sprintf(
			"Announced Prefixes (%d)",
			len(prefixes),
		),
		nodes...,
	)
}

func asnNeighboursNode(
	counts knowledge.NetworkASNNeighbourCounts,
	neighbours []knowledge.NetworkASNNeighbour,
) reportNode {
	if counts.Unique == 0 &&
		len(neighbours) == 0 {
		return reportNode{}
	}

	children :=
		[]reportNode{
			reportValue(
				"Unique",
				fmt.Sprintf(
					"%d",
					counts.Unique,
				),
			),
			reportValue(
				"Left",
				fmt.Sprintf(
					"%d",
					counts.Left,
				),
			),
			reportValue(
				"Right",
				fmt.Sprintf(
					"%d",
					counts.Right,
				),
			),
			reportValue(
				"Uncertain",
				fmt.Sprintf(
					"%d",
					counts.Uncertain,
				),
			),
		}

	if len(neighbours) > 0 {
		limit :=
			len(neighbours)

		if !output.Full() &&
			limit >
				maxCompactNetworkItems {
			limit =
				maxCompactNetworkItems
		}

		var neighbourNodes []reportNode

		for index := 0; index < limit; index++ {
			neighbour :=
				neighbours[index]

			label :=
				fmt.Sprintf(
					"AS%d",
					neighbour.ASN,
				)

			if strings.TrimSpace(
				neighbour.Position,
			) != "" {
				label +=
					" [" +
						strings.TrimSpace(
							neighbour.Position,
						) +
						"]"
			}

			var details []reportNode

			if neighbour.PathCount > 0 {
				details =
					append(
						details,
						reportValue(
							"Path Count",
							fmt.Sprintf(
								"%d",
								neighbour.PathCount,
							),
						),
					)
			}

			if neighbour.PeerCountV4 > 0 {
				details =
					append(
						details,
						reportValue(
							"IPv4 Peers",
							fmt.Sprintf(
								"%d",
								neighbour.PeerCountV4,
							),
						),
					)
			}

			if neighbour.PeerCountV6 > 0 {
				details =
					append(
						details,
						reportValue(
							"IPv6 Peers",
							fmt.Sprintf(
								"%d",
								neighbour.PeerCountV6,
							),
						),
					)
			}

			neighbourNodes =
				append(
					neighbourNodes,
					reportNode{
						Label: label,

						Children: details,
					},
				)
		}

		if limit < len(neighbours) {
			neighbourNodes =
				append(
					neighbourNodes,
					reportItem(
						fmt.Sprintf(
							"... %d more",
							len(neighbours)-limit,
						),
					),
				)
		}

		children =
			append(
				children,
				reportSection(
					"Neighbours",
					neighbourNodes...,
				),
			)
	}

	return reportSection(
		"ASN Neighbours",
		children...,
	)
}

func reverseDNSNode(
	names []string,
) reportNode {
	var nodes []reportNode

	for _, name := range names {
		nodes =
			append(
				nodes,
				reportItem(
					name,
				),
			)
	}

	return reportSection(
		"Reverse DNS",
		nodes...,
	)
}

func networkClassificationNode(
	result knowledge.NetworkResult,
) reportNode {
	children :=
		[]reportNode{
			reportValue(
				"Connection Type",
				result.ConnectionType,
			),
			networkFlagNode(
				"Mobile",
				result.Mobile,
			),
			networkFlagNode(
				"Hosting",
				result.Hosting,
			),
			networkFlagNode(
				"Proxy",
				result.Proxy,
			),
			networkFlagNode(
				"VPN",
				result.VPN,
			),
			networkFlagNode(
				"Tor",
				result.Tor,
			),
			reportValue(
				"Threat",
				result.Threat,
			),
		}

	return reportSection(
		"Connection Classification",
		children...,
	)
}

func networkFlagNode(
	name string,
	flag knowledge.NetworkFlagResult,
) reportNode {
	if !flag.Available {
		return reportNode{}
	}

	node :=
		reportBool(
			name,
			flag.Value,
		)

	if output.Full() &&
		strings.TrimSpace(
			flag.Source,
		) != "" {
		node.Children =
			[]reportNode{
				reportValue(
					"Source",
					flag.Source,
				),
			}
	}

	return node
}

func anycastNode(
	result knowledge.NetworkAnycastResult,
) reportNode {
	if !result.Available {
		return reportNode{}
	}

	children :=
		[]reportNode{
			reportBool(
				"Detected",
				result.Anycast,
			),
			reportValue(
				"Confidence",
				result.Confidence,
			),
			reportValue(
				"Prefix",
				result.Prefix,
			),
			reportValue(
				"Mapped Prefix",
				result.MappedPrefix,
			),
			reportValue(
				"Backing Prefix",
				result.BackingPrefix,
			),
			anycastASNsNode(
				result,
			),
			reportValue(
				"Observation Date",
				result.Date,
			),
		}

	if output.Full() {
		children =
			append(
				children,
				anycastMeasurementsNode(
					result,
				),
				anycastLocationsNode(
					result,
				),
			)
	}

	return reportSection(
		"Anycast",
		children...,
	)
}

func anycastASNsNode(
	result knowledge.NetworkAnycastResult,
) reportNode {
	if len(result.ASNs) == 0 {
		return reportNode{}
	}

	values :=
		make(
			[]string,
			0,
			len(result.ASNs),
		)

	for _, asn := range result.ASNs {
		values =
			append(
				values,
				fmt.Sprintf(
					"AS%d",
					asn,
				),
			)
	}

	return reportValue(
		"ASNs",
		strings.Join(
			values,
			", ",
		),
	)
}

func anycastMeasurementsNode(
	result knowledge.NetworkAnycastResult,
) reportNode {
	return reportSection(
		"Measurements",
		reportValue(
			"AB ICMP",
			fmt.Sprintf(
				"%d",
				result.ABICMP,
			),
		),
		reportValue(
			"AB TCP",
			fmt.Sprintf(
				"%d",
				result.ABTCP,
			),
		),
		reportValue(
			"AB DNS",
			fmt.Sprintf(
				"%d",
				result.ABDNS,
			),
		),
		reportValue(
			"GCD ICMP",
			fmt.Sprintf(
				"%d",
				result.GCDICMP,
			),
		),
		reportValue(
			"GCD TCP",
			fmt.Sprintf(
				"%d",
				result.GCDTCP,
			),
		),
	)
}

func anycastLocationsNode(
	result knowledge.NetworkAnycastResult,
) reportNode {
	var nodes []reportNode

	for _, location := range result.Locations {
		label :=
			strings.TrimSpace(
				location.ID,
			)

		if label == "" {
			label =
				"Location"
		}

		var details []reportNode

		details =
			append(
				details,
				reportValue(
					"City",
					location.City,
				),
				reportValue(
					"Country",
					location.Country,
				),
				reportValue(
					"Coordinates",
					fmt.Sprintf(
						"%.6f, %.6f",
						location.Latitude,
						location.Longitude,
					),
				),
			)

		nodes =
			append(
				nodes,
				reportNode{
					Label: label,

					Children: details,
				},
			)
	}

	return reportSection(
		fmt.Sprintf(
			"Observed Locations (%d)",
			len(result.Locations),
		),
		nodes...,
	)
}

func networkSourcesNode(
	sources []string,
) reportNode {
	var nodes []reportNode

	for _, source := range sources {
		nodes =
			append(
				nodes,
				reportItem(
					source,
				),
			)
	}

	return reportSection(
		"Sources",
		nodes...,
	)
}

func networkErrorsNode(
	errors []knowledge.NetworkProviderError,
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
