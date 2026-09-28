package ip

import (
	"fmt"
	"strings"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

const maxCompactHistoryTimelines = 5

func historyNodes(
	result knowledge.HistoryResult,
) []reportNode {
	nodes :=
		[]reportNode{
			historyPresenceNode(
				result.Status,
			),

			routingHistoryNode(
				result.Routing,
			),
		}

	if output.Full() {

		nodes =
			append(
				nodes,
				historySourcesNode(
					result.Sources,
				),

				historyErrorsNode(
					result.Errors,
				),
			)
	}

	return compactReportNodes(
		nodes,
	)
}

func historyPresenceNode(
	result knowledge.HistoryRoutingStatusResult,
) reportNode {
	if !result.Available {
		return reportNode{}
	}

	children :=
		[]reportNode{
			historySeenNode(
				"First Seen",
				result.FirstSeen,
			),

			historySeenNode(
				"Last Seen",
				result.LastSeen,
			),
		}

	if output.Full() {

		children =
			append(
				children,
				reportValue(
					"Query Time",
					result.QueryTime,
				),

				reportValue(
					"Resource",
					result.Resource,
				),

				reportValue(
					"Source",
					result.Source,
				),
			)
	}

	return reportSection(
		"BGP Presence",
		children...,
	)
}

func historySeenNode(
	label string,
	result knowledge.HistorySeenResult,
) reportNode {
	if !result.Available {
		return reportNode{}
	}

	return reportSection(
		label,

		reportValue(
			"Time",
			result.Time,
		),

		reportValue(
			"Origin",
			historyOrigin(
				result.Origin,
			),
		),

		reportValue(
			"Prefix",
			result.Prefix,
		),
	)
}

func routingHistoryNode(
	result knowledge.HistoryRoutingResult,
) reportNode {
	if !result.Available {
		return reportNode{}
	}

	var origins []reportNode

	for _, origin := range result.Origins {

		var prefixes []reportNode

		for _, prefix := range origin.Prefixes {

			var timelines []reportNode

			limit :=
				len(
					prefix.Timelines,
				)

			if !output.Full() &&
				limit >
					maxCompactHistoryTimelines {

				limit =
					maxCompactHistoryTimelines
			}

			for index := 0; index < limit; index++ {

				timeline :=
					prefix.Timelines[index]

				label :=
					historyPeriodLabel(
						timeline,
					)

				var children []reportNode

				if timeline.FullPeersSeeing >
					0 {

					children =
						append(
							children,
							reportValue(
								"Full Peers Seeing",
								fmt.Sprintf(
									"%d",
									timeline.
										FullPeersSeeing,
								),
							),
						)
				}

				if timeline.VisibilityAvailable {

					visibility :=
						"unavailable"

					if timeline.Visibility >= 0 {

						visibility =
							fmt.Sprintf(
								"%.2f%%",
								timeline.
									Visibility*
									100,
							)
					}

					children =
						append(
							children,
							reportValue(
								"Visibility",
								visibility,
							),
						)
				}

				timelines =
					append(
						timelines,
						reportNode{
							Label: label,

							Children: children,
						},
					)
			}

			if limit <
				len(
					prefix.Timelines,
				) {

				timelines =
					append(
						timelines,
						reportItem(
							fmt.Sprintf(
								"... %d older periods",
								len(
									prefix.
										Timelines,
								)-limit,
							),
						),
					)
			}

			prefixes =
				append(
					prefixes,
					reportSection(
						prefix.Prefix,
						timelines...,
					),
				)
		}

		origins =
			append(
				origins,
				reportSection(
					historyOrigin(
						origin.Origin,
					),
					prefixes...,
				),
			)
	}

	children :=
		origins

	if output.Full() {

		children =
			append(
				children,
				reportValue(
					"Query Start",
					result.QueryStartTime,
				),

				reportValue(
					"Query End",
					result.QueryEndTime,
				),

				reportValue(
					"Resource",
					result.Resource,
				),

				reportValue(
					"Source",
					result.Source,
				),
			)
	}

	return reportSection(
		"Routing History",
		children...,
	)
}

func historyPeriodLabel(
	timeline knowledge.HistoryRoutingTimeline,
) string {
	start :=
		strings.TrimSpace(
			timeline.StartTime,
		)

	end :=
		strings.TrimSpace(
			timeline.EndTime,
		)

	switch {

	case start != "" &&
		end != "":

		return start +
			" → " +
			end

	case start != "":

		return start +
			" → present"

	case end != "":

		return "until " +
			end
	}

	return "Routing period"
}

func historyOrigin(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	if strings.HasPrefix(
		strings.ToUpper(
			value,
		),
		"AS",
	) {
		return value
	}

	return "AS" +
		value
}

func historySourcesNode(
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

func historyErrorsNode(
	errors []knowledge.HistoryProviderError,
) reportNode {
	var nodes []reportNode

	for _, providerErr := range errors {

		nodes =
			append(
				nodes,
				reportNode{
					Label: providerErr.
						Provider,

					Children: []reportNode{
						reportValue(
							"Error",
							providerErr.
								Error,
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
