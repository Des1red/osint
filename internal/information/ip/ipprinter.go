package ip

import (
	"fmt"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func IpPrinter() bool {
	result :=
		knowledge.Data.IP

	var sections []reportNode

	if result.RDAP.IP != "" {

		sections =
			append(
				sections,
				reportNode{
					Label: "RDAP",

					Children: rdapNodes(
						result.RDAP,
					),
				},
			)
	}

	if result.Geo.IP != "" {

		sections =
			append(
				sections,
				reportNode{
					Label: "Geolocation",

					Children: geoNodes(
						result.Geo,
					),
				},
			)
	}

	if result.Network.IP != "" {

		sections =
			append(
				sections,
				reportNode{
					Label: "Network",

					Children: networkNodes(
						result.Network,
					),
				},
			)
	}

	if result.History.IP != "" {

		sections =
			append(
				sections,
				reportNode{
					Label: "History",

					Children: historyNodes(
						result.History,
					),
				},
			)
	}

	if result.Reputation.IP != "" {

		sections =
			append(
				sections,
				reportNode{
					Label: "Reputation",

					Children: reputationNodes(
						result.Reputation,
					),
				},
			)
	}

	sections =
		compactReportNodes(
			sections,
		)

	if len(sections) == 0 {
		return false
	}

	fmt.Fprintln(
		output.Writer(),
		"IP Intelligence",
	)

	fmt.Fprintln(
		output.Writer(),
		"===============",
	)

	printReportNodes(
		"",
		sections,
	)

	return true
}
