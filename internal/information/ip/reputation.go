package ip

import (
	"fmt"
	"time"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func reputationNodes(
	result knowledge.ReputationResult,
) []reportNode {
	nodes :=
		[]reportNode{
			abuseIPDBNode(
				result.AbuseIPDB,
			),
			virusTotalNode(
				result.VirusTotal,
			),
			ipQSNode(
				result.IPQS,
			),
		}

	if output.Full() {
		nodes =
			append(
				nodes,
				reputationErrorsNode(
					result.Errors,
				),
			)
	}

	return compactReportNodes(
		nodes,
	)
}

func abuseIPDBNode(
	result knowledge.AbuseIPDBResult,
) reportNode {
	if !result.Available {
		return reportNode{}
	}

	children :=
		[]reportNode{
			reportValue(
				"Abuse Confidence Score",
				fmt.Sprintf(
					"%d",
					result.AbuseConfidenceScore,
				),
			),
			reportValue(
				"Total Reports",
				fmt.Sprintf(
					"%d",
					result.TotalReports,
				),
			),
			reportValue(
				"Distinct Users",
				fmt.Sprintf(
					"%d",
					result.DistinctUsers,
				),
			),
			reportValue(
				"Last Reported",
				result.LastReportedAt,
			),
			reportValue(
				"Usage Type",
				result.UsageType,
			),
			reportValue(
				"Domain",
				result.Domain,
			),
			reportBool(
				"Public",
				result.IsPublic,
			),
		}

	if result.IsWhitelisted != nil {
		children =
			append(
				children,
				reportBool(
					"Whitelisted",
					*result.IsWhitelisted,
				),
			)
	}

	if output.Full() {
		children =
			append(
				children,
				reportValue(
					"Source",
					result.Source,
				),
			)
	}

	return reportSection(
		"AbuseIPDB",
		children...,
	)
}

func virusTotalNode(
	result knowledge.VirusTotalResult,
) reportNode {
	if !result.Available {
		return reportNode{}
	}

	children :=
		[]reportNode{
			reportValue(
				"Reputation",
				fmt.Sprintf(
					"%d",
					result.Reputation,
				),
			),
			reportSection(
				"Analysis",
				reportValue(
					"Harmless",
					fmt.Sprintf(
						"%d",
						result.Stats.Harmless,
					),
				),
				reportValue(
					"Malicious",
					fmt.Sprintf(
						"%d",
						result.Stats.Malicious,
					),
				),
				reportValue(
					"Suspicious",
					fmt.Sprintf(
						"%d",
						result.Stats.Suspicious,
					),
				),
				reportValue(
					"Undetected",
					fmt.Sprintf(
						"%d",
						result.Stats.Undetected,
					),
				),
				reportValue(
					"Timeout",
					fmt.Sprintf(
						"%d",
						result.Stats.Timeout,
					),
				),
			),
			reportSection(
				"Community Votes",
				reportValue(
					"Harmless",
					fmt.Sprintf(
						"%d",
						result.Votes.Harmless,
					),
				),
				reportValue(
					"Malicious",
					fmt.Sprintf(
						"%d",
						result.Votes.Malicious,
					),
				),
			),
			virusTotalTagsNode(
				result.Tags,
			),
		}

	if result.LastAnalysisDate > 0 {
		children =
			append(
				children,
				reportValue(
					"Last Analysis",
					time.Unix(
						result.LastAnalysisDate,
						0,
					).Format(
						time.RFC3339,
					),
				),
			)
	}

	if output.Full() {
		children =
			append(
				children,
				reportValue(
					"Source",
					result.Source,
				),
			)
	}

	return reportSection(
		"VirusTotal",
		children...,
	)
}

func virusTotalTagsNode(
	tags []string,
) reportNode {
	var nodes []reportNode

	for _, tag := range tags {
		nodes =
			append(
				nodes,
				reportItem(
					tag,
				),
			)
	}

	return reportSection(
		"Tags",
		nodes...,
	)
}

func ipQSNode(
	result knowledge.IPQSResult,
) reportNode {
	if !result.Available {
		return reportNode{}
	}

	children :=
		[]reportNode{
			reportValue(
				"Fraud Score",
				fmt.Sprintf(
					"%d",
					result.FraudScore,
				),
			),
			reportBool(
				"Proxy",
				result.Proxy,
			),
			reportBool(
				"VPN",
				result.VPN,
			),
			reportBool(
				"Tor",
				result.Tor,
			),
			reportBool(
				"Bot Status",
				result.BotStatus,
			),
			reportBool(
				"Recent Abuse",
				result.RecentAbuse,
			),
			reportBool(
				"Frequent Abuser",
				result.FrequentAbuser,
			),
			reportValue(
				"Abuse Velocity",
				result.AbuseVelocity,
			),
		}

	if output.Full() {
		children =
			append(
				children,
				reportValue(
					"Source",
					result.Source,
				),
			)
	}

	return reportSection(
		"IPQualityScore",
		children...,
	)
}

func reputationErrorsNode(
	errors []knowledge.ReputationProviderError,
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
