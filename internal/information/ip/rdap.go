package ip

import (
	"fmt"
	"strings"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func rdapNodes(
	result knowledge.RDAPResult,
) []reportNode {
	nodes :=
		[]reportNode{
			rdapRegistrationNode(
				result,
			),
			rdapAddressRangeNode(
				result,
			),
			rdapRoutingNode(
				result,
			),
			rdapEntitiesNode(
				result.Entities,
			),
		}

	if output.Full() {
		nodes =
			append(
				nodes,
				rdapRegistryMetadataNode(
					result,
				),
				rdapEventsNode(
					"Events",
					result.Events,
				),
				rdapRemarksNode(
					"Remarks",
					result.Remarks,
				),
				rdapLinksNode(
					"Links",
					result.Links,
				),
				reportValue(
					"Source",
					result.Source,
				),
			)
	}

	return compactReportNodes(
		nodes,
	)
}

func rdapRegistrationNode(
	result knowledge.RDAPResult,
) reportNode {
	return reportSection(
		"Network Registration",
		reportValue(
			"IP",
			result.IP,
		),
		reportValue(
			"Object Class",
			result.ObjectClassName,
		),
		reportValue(
			"Name",
			result.Name,
		),
		reportValue(
			"Handle",
			result.Handle,
		),
		reportValue(
			"Parent Handle",
			result.ParentHandle,
		),
		reportValue(
			"Type",
			result.Type,
		),
		reportSlice(
			"Status",
			result.Status,
		),
	)
}

func rdapAddressRangeNode(
	result knowledge.RDAPResult,
) reportNode {
	return reportSection(
		"Address Range",
		reportValue(
			"Start IP",
			result.StartIP,
		),
		reportValue(
			"End IP",
			result.EndIP,
		),
		reportValue(
			"IP Version",
			result.IPVersion,
		),
		reportValue(
			"Country",
			result.Country,
		),
		rdapCIDRsNode(
			result.CIDRs,
		),
	)
}

func rdapCIDRsNode(
	cidrs []knowledge.RDAPCIDR,
) reportNode {
	var nodes []reportNode

	for _, cidr := range cidrs {
		switch {
		case strings.TrimSpace(
			cidr.IPv4Prefix,
		) != "":

			nodes =
				append(
					nodes,
					reportItem(
						fmt.Sprintf(
							"%s/%d",
							cidr.IPv4Prefix,
							cidr.Length,
						),
					),
				)

		case strings.TrimSpace(
			cidr.IPv6Prefix,
		) != "":

			nodes =
				append(
					nodes,
					reportItem(
						fmt.Sprintf(
							"%s/%d",
							cidr.IPv6Prefix,
							cidr.Length,
						),
					),
				)
		}
	}

	return reportSection(
		fmt.Sprintf(
			"CIDRs (%d)",
			len(nodes),
		),
		nodes...,
	)
}

func rdapRoutingNode(
	result knowledge.RDAPResult,
) reportNode {
	var nodes []reportNode

	for _, asn := range result.OriginASNs {
		nodes =
			append(
				nodes,
				reportItem(
					fmt.Sprintf(
						"AS%d",
						asn,
					),
				),
			)
	}

	return reportSection(
		"Routing",
		reportSection(
			fmt.Sprintf(
				"Origin ASNs (%d)",
				len(nodes),
			),
			nodes...,
		),
	)
}

func rdapRegistryMetadataNode(
	result knowledge.RDAPResult,
) reportNode {
	return reportSection(
		"Registry Metadata",
		reportValue(
			"Port 43",
			result.Port43,
		),
		reportSlice(
			"Conformance",
			result.Conformance,
		),
	)
}

func rdapEventsNode(
	label string,
	events []knowledge.RDAPEvent,
) reportNode {
	var nodes []reportNode

	for _, event := range events {
		eventLabel :=
			strings.TrimSpace(
				event.Action,
			)

		date :=
			strings.TrimSpace(
				event.Date,
			)

		if eventLabel == "" {
			eventLabel =
				"Event"
		}

		if date != "" {
			eventLabel +=
				": " +
					date
		}

		var children []reportNode

		if strings.TrimSpace(
			event.Actor,
		) != "" {
			children =
				append(
					children,
					reportValue(
						"Actor",
						event.Actor,
					),
				)
		}

		nodes =
			append(
				nodes,
				reportNode{
					Label: eventLabel,

					Children: children,
				},
			)
	}

	return reportSection(
		label,
		nodes...,
	)
}

func rdapEntitiesNode(
	entities []knowledge.RDAPEntity,
) reportNode {
	var nodes []reportNode

	for _, entity := range entities {
		nodes =
			append(
				nodes,
				rdapEntityNode(
					entity,
				),
			)
	}

	return reportSection(
		"Entities",
		nodes...,
	)
}

func rdapEntityNode(
	entity knowledge.RDAPEntity,
) reportNode {
	handle :=
		strings.TrimSpace(
			entity.Handle,
		)

	name :=
		strings.TrimSpace(
			entity.Name,
		)

	label :=
		handle

	if label == "" {
		label =
			name
	} else if name != "" {
		label +=
			" (" +
				name +
				")"
	}

	if label == "" {
		label =
			"Entity"
	}

	children :=
		[]reportNode{
			reportSlice(
				"Roles",
				entity.Roles,
			),
			reportSlice(
				"Status",
				entity.Status,
			),
			reportValue(
				"Kind",
				entity.Kind,
			),
			reportValue(
				"Organization",
				entity.Organization,
			),
			reportValue(
				"Title",
				entity.Title,
			),
			reportValue(
				"Role",
				entity.Role,
			),
			rdapEntityContactsNode(
				entity,
			),
		}

	if output.Full() {
		children =
			append(
				children,
				reportValue(
					"Port 43",
					entity.Port43,
				),
				rdapPublicIDsNode(
					entity.PublicIDs,
				),
				rdapEventsNode(
					"Events",
					entity.Events,
				),
				rdapRemarksNode(
					"Remarks",
					entity.Remarks,
				),
				rdapLinksNode(
					"Links",
					entity.Links,
				),
			)
	}

	if len(entity.Entities) > 0 {
		children =
			append(
				children,
				rdapEntitiesNode(
					entity.Entities,
				),
			)
	}

	return reportNode{
		Label: label,

		Children: compactReportNodes(
			children,
		),
	}
}

func rdapEntityContactsNode(
	entity knowledge.RDAPEntity,
) reportNode {
	var nodes []reportNode

	for _, email := range entity.Emails {
		nodes =
			append(
				nodes,
				reportValue(
					"Email",
					email,
				),
			)
	}

	for _, phone := range entity.Phones {
		nodes =
			append(
				nodes,
				reportValue(
					"Phone",
					phone,
				),
			)
	}

	for _, address := range entity.Addresses {
		nodes =
			append(
				nodes,
				reportValue(
					"Address",
					address,
				),
			)
	}

	for _, url := range entity.URLs {
		nodes =
			append(
				nodes,
				reportValue(
					"URL",
					url,
				),
			)
	}

	return reportSection(
		"Contacts",
		nodes...,
	)
}

func rdapPublicIDsNode(
	publicIDs []knowledge.RDAPPublicID,
) reportNode {
	var nodes []reportNode

	for _, publicID := range publicIDs {
		value :=
			strings.TrimSpace(
				publicID.Identifier,
			)

		if value == "" {
			continue
		}

		if strings.TrimSpace(
			publicID.Type,
		) != "" {
			value =
				strings.TrimSpace(
					publicID.Type,
				) +
					": " +
					value
		}

		nodes =
			append(
				nodes,
				reportItem(
					value,
				),
			)
	}

	return reportSection(
		"Public IDs",
		nodes...,
	)
}

func rdapRemarksNode(
	label string,
	remarks []knowledge.RDAPRemark,
) reportNode {
	var nodes []reportNode

	for _, remark := range remarks {
		remarkLabel :=
			strings.TrimSpace(
				remark.Title,
			)

		if remarkLabel == "" {
			remarkLabel =
				"Remark"
		}

		var descriptions []reportNode

		for _, description := range remark.Description {
			descriptions =
				append(
					descriptions,
					reportItem(
						description,
					),
				)
		}

		children :=
			[]reportNode{
				reportValue(
					"Type",
					remark.Type,
				),
				reportSection(
					"Description",
					descriptions...,
				),
			}

		if output.Full() {
			children =
				append(
					children,
					rdapLinksNode(
						"Links",
						remark.Links,
					),
				)
		}

		nodes =
			append(
				nodes,
				reportNode{
					Label: remarkLabel,

					Children: compactReportNodes(
						children,
					),
				},
			)
	}

	return reportSection(
		label,
		nodes...,
	)
}

func rdapLinksNode(
	label string,
	links []knowledge.RDAPLink,
) reportNode {
	var nodes []reportNode

	for _, link := range links {
		href :=
			strings.TrimSpace(
				link.Href,
			)

		if href == "" {
			continue
		}

		if strings.TrimSpace(
			link.Rel,
		) != "" {
			nodes =
				append(
					nodes,
					reportValue(
						link.Rel,
						href,
					),
				)

			continue
		}

		nodes =
			append(
				nodes,
				reportItem(
					href,
				),
			)
	}

	return reportSection(
		label,
		nodes...,
	)
}
