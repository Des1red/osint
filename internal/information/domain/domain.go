package domain

import (
	"fmt"
	"strconv"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

const terminalSubdomainLimit = 30

type treeSection struct {
	Name string

	Values []string
}

func DomainPrinter() bool {
	result :=
		knowledge.Data.Domain

	fmt.Fprintln(
		output.Writer(),
		"DNS Lookup",
	)

	fmt.Fprintln(
		output.Writer(),
		"----------",
	)

	if !result.Found {

		output.TreeItem(
			"",
			true,
			result.Domain,
		)

		prefix :=
			output.TreeChildPrefix(
				"",
				true,
			)

		output.TreeItem(
			prefix,
			true,
			"Nothing found",
		)

		return false
	}

	printDNS(
		result,
	)

	printSubdomains(
		result,
	)

	return true
}

func printDNS(
	result knowledge.DomainResult,
) {
	output.TreeItem(
		"",
		true,
		result.Domain,
	)

	prefix :=
		output.TreeChildPrefix(
			"",
			true,
		)

	sections :=
		dnsSections(
			result,
		)

	for index, section := range sections {

		last :=
			index ==
				len(sections)-1

		output.TreeItem(
			prefix,
			last,
			section.Name,
		)

		childPrefix :=
			output.TreeChildPrefix(
				prefix,
				last,
			)

		for valueIndex, value := range section.Values {

			output.TreeItem(
				childPrefix,
				valueIndex ==
					len(section.Values)-1,
				value,
			)
		}
	}
}

func dnsSections(
	result knowledge.DomainResult,
) []treeSection {
	var sections []treeSection

	if len(result.A) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Name: "A",

					Values: result.A,
				},
			)
	}

	if len(result.AAAA) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Name: "AAAA",

					Values: result.AAAA,
				},
			)
	}

	if result.CNAME != "" {

		sections =
			append(
				sections,
				treeSection{
					Name: "CNAME",

					Values: []string{
						result.CNAME,
					},
				},
			)
	}

	if len(result.MX) > 0 {

		values :=
			make(
				[]string,
				0,
				len(result.MX),
			)

		for _, record := range result.MX {

			values =
				append(
					values,
					fmt.Sprintf(
						"%d %s",
						record.Preference,
						record.Host,
					),
				)
		}

		sections =
			append(
				sections,
				treeSection{
					Name: "MX",

					Values: values,
				},
			)
	}

	if len(result.NS) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Name: "NS",

					Values: result.NS,
				},
			)
	}

	if len(result.TXT) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Name: "TXT",

					Values: result.TXT,
				},
			)
	}

	return sections
}

func printSubdomains(
	result knowledge.DomainResult,
) {
	if result.SubdomainCandidates == 0 &&
		len(result.Subdomains) == 0 {

		return
	}

	fmt.Fprintln(
		output.Writer(),
	)

	fmt.Fprintln(
		output.Writer(),
		"Subdomains",
	)

	fmt.Fprintln(
		output.Writer(),
		"----------",
	)

	potential :=
		0

	for _, subdomain := range result.Subdomains {

		if subdomain.Takeover.Potential {

			potential++
		}
	}

	fields :=
		[]output.TreeField{
			{
				Name: "Source",

				Value: result.SubdomainSource,
			},
			{
				Name: "Candidates Discovered",

				Value: strconv.Itoa(
					result.SubdomainCandidates,
				),
			},
			{
				Name: "Candidates Checked",

				Value: strconv.Itoa(
					result.SubdomainChecked,
				),
			},
			{
				Name: "Current DNS Results",

				Value: strconv.Itoa(
					len(
						result.Subdomains,
					),
				),
			},
			{
				Name: "Potential Takeover Risks",

				Value: strconv.Itoa(
					potential,
				),
			},
		}

	hasResults :=
		len(result.Subdomains) > 0

	for index, field := range fields {

		last :=
			index ==
				len(fields)-1 &&
				!hasResults

		output.TreeValue(
			"",
			last,
			field.Name,
			field.Value,
		)
	}

	if !hasResults {

		return
	}

	output.TreeItem(
		"",
		true,
		"Results",
	)

	resultsPrefix :=
		output.TreeChildPrefix(
			"",
			true,
		)

	limit :=
		len(
			result.Subdomains,
		)

	if !output.Full() &&
		limit >
			terminalSubdomainLimit {

		limit =
			terminalSubdomainLimit
	}

	for index := 0; index < limit; index++ {

		hasMoreMessage :=
			limit <
				len(
					result.Subdomains,
				)

		last :=
			index ==
				limit-1 &&
				!hasMoreMessage

		printSubdomain(
			resultsPrefix,
			last,
			result.Subdomains[index],
		)
	}

	if limit <
		len(result.Subdomains) {

		output.TreeItem(
			resultsPrefix,
			true,
			fmt.Sprintf(
				"... %d more subdomains",
				len(result.Subdomains)-limit,
			),
		)
	}
}

func printSubdomain(
	prefix string,
	last bool,
	result knowledge.DomainSubdomainResult,
) {
	output.TreeItem(
		prefix,
		last,
		result.Host,
	)

	childPrefix :=
		output.TreeChildPrefix(
			prefix,
			last,
		)

	sections :=
		subdomainSections(
			result,
		)

	for index, section := range sections {

		sectionLast :=
			index ==
				len(sections)-1

		if section.Name ==
			"CNAME" {

			output.TreeValue(
				childPrefix,
				sectionLast,
				"CNAME",
				section.Values[0],
			)

			continue
		}

		output.TreeItem(
			childPrefix,
			sectionLast,
			section.Name,
		)

		valuePrefix :=
			output.TreeChildPrefix(
				childPrefix,
				sectionLast,
			)

		for valueIndex, value := range section.Values {

			output.TreeItem(
				valuePrefix,
				valueIndex ==
					len(section.Values)-1,
				value,
			)
		}
	}
}

func subdomainSections(
	result knowledge.DomainSubdomainResult,
) []treeSection {
	var sections []treeSection

	if len(result.A) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Name: "A",

					Values: result.A,
				},
			)
	}

	if len(result.AAAA) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Name: "AAAA",

					Values: result.AAAA,
				},
			)
	}

	if result.CNAME != "" {

		sections =
			append(
				sections,
				treeSection{
					Name: "CNAME",

					Values: []string{
						result.CNAME,
					},
				},
			)
	}

	if result.Takeover.Potential {

		values :=
			[]string{
				"Potential",
			}

		if result.Takeover.Provider != "" {

			values =
				append(
					values,
					"Provider: "+
						result.Takeover.Provider,
				)
		}

		if result.Takeover.Target != "" {

			values =
				append(
					values,
					"Target: "+
						result.Takeover.Target,
				)
		}

		if result.Takeover.Reason != "" {

			values =
				append(
					values,
					"Reason: "+
						result.Takeover.Reason,
				)
		}

		sections =
			append(
				sections,
				treeSection{
					Name: "Takeover Risk",

					Values: values,
				},
			)
	}

	return sections
}
