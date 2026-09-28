package ip

import (
	"fmt"
	"strings"

	"osint/internal/information/output"
)

type reportNode struct {
	Label string

	Children []reportNode
}

func reportSection(
	label string,
	children ...reportNode,
) reportNode {
	label =
		strings.TrimSpace(
			label,
		)

	children =
		compactReportNodes(
			children,
		)

	if label == "" ||
		len(children) == 0 {
		return reportNode{}
	}

	return reportNode{
		Label: label,

		Children: children,
	}
}

func reportItem(
	value string,
) reportNode {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return reportNode{}
	}

	return reportNode{
		Label: value,
	}
}

func reportValue(
	name string,
	value string,
) reportNode {
	name =
		strings.TrimSpace(
			name,
		)

	value =
		strings.TrimSpace(
			value,
		)

	if name == "" ||
		value == "" {
		return reportNode{}
	}

	return reportNode{
		Label: fmt.Sprintf(
			"%s: %s",
			name,
			value,
		),
	}
}

func reportBool(
	name string,
	value bool,
) reportNode {
	return reportValue(
		name,
		fmt.Sprintf(
			"%t",
			value,
		),
	)
}

func reportSlice(
	name string,
	values []string,
) reportNode {
	var filtered []string

	for _, value := range values {
		value =
			strings.TrimSpace(
				value,
			)

		if value == "" {
			continue
		}

		filtered =
			append(
				filtered,
				value,
			)
	}

	if len(filtered) == 0 {
		return reportNode{}
	}

	return reportValue(
		name,
		strings.Join(
			filtered,
			", ",
		),
	)
}

func compactReportNodes(
	nodes []reportNode,
) []reportNode {
	result :=
		make(
			[]reportNode,
			0,
			len(nodes),
		)

	for _, node := range nodes {
		node.Label =
			strings.TrimSpace(
				node.Label,
			)

		if node.Label == "" {
			continue
		}

		node.Children =
			compactReportNodes(
				node.Children,
			)

		result =
			append(
				result,
				node,
			)
	}

	return result
}

func printReportNodes(
	prefix string,
	nodes []reportNode,
) {
	nodes =
		compactReportNodes(
			nodes,
		)

	for index, node := range nodes {
		last :=
			index ==
				len(nodes)-1

		output.TreeItem(
			prefix,
			last,
			node.Label,
		)

		if len(node.Children) == 0 {
			continue
		}

		printReportNodes(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			node.Children,
		)
	}
}
