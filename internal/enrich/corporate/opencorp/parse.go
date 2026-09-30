package opencorp

import (
	"html"
	"regexp"
	"strings"
)

var officerLinkPattern = regexp.MustCompile(
	`(?i)href=["'](/officers/\d+)["']`,
)

var headingPattern = regexp.MustCompile(
	`(?is)<h1\b[^>]*>(.*?)</h1>`,
)

var companyLinkPattern = regexp.MustCompile(
	`(?is)<a\b[^>]*href=["'](/companies/([^/"']+)/[^"']+)["'][^>]*>(.*?)</a>`,
)

var definitionPattern = regexp.MustCompile(
	`(?is)<dt\b[^>]*>(.*?)</dt>\s*<dd\b[^>]*>(.*?)</dd>`,
)

var tableFieldPattern = regexp.MustCompile(
	`(?is)<tr\b[^>]*>.*?<th\b[^>]*>(.*?)</th>.*?<td\b[^>]*>(.*?)</td>.*?</tr>`,
)

var tagPattern = regexp.MustCompile(
	`(?s)<[^>]+>`,
)

var whitespacePattern = regexp.MustCompile(
	`\s+`,
)

func parseResults(
	body string,
) Result {
	var result Result

	matches :=
		officerLinkPattern.FindAllStringSubmatch(
			body,
			-1,
		)

	seen :=
		make(
			map[string]struct{},
		)

	for _, match := range matches {

		if len(match) < 2 {

			continue
		}

		path :=
			strings.TrimSpace(
				match[1],
			)

		if path == "" {

			continue
		}

		officer :=
			officerFromPath(
				path,
			)

		if officer.ProfileURL == "" {

			continue
		}

		key :=
			strings.ToLower(
				officer.ProfileURL,
			)

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		result.Officers =
			append(
				result.Officers,
				officer,
			)

		if len(result.Officers) >=
			maxResults {

			break
		}
	}

	return result
}

func parseOfficer(
	body string,
	officer Officer,
) Officer {
	name :=
		parseOfficerName(
			body,
		)

	if name != "" {

		officer.Name =
			name
	}

	position :=
		fieldValue(
			body,
			"position",
		)

	if position != "" {

		officer.Position =
			position
	}

	startDate :=
		fieldValue(
			body,
			"start date",
			"start_date",
		)

	if startDate != "" {

		officer.StartDate =
			startDate
	}

	endDate :=
		fieldValue(
			body,
			"end date",
			"end_date",
		)

	if endDate != "" {

		officer.EndDate =
			endDate
	}

	company,
		companyURL,
		jurisdiction :=
		parseCompany(
			body,
		)

	if company != "" {

		officer.Company =
			company
	}

	if companyURL != "" {

		officer.CompanyURL =
			companyURL
	}

	pageJurisdiction :=
		fieldValue(
			body,
			"jurisdiction",
			"jurisdiction code",
			"jurisdiction_code",
		)

	if pageJurisdiction != "" {

		officer.Jurisdiction =
			pageJurisdiction

	} else if jurisdiction != "" {

		officer.Jurisdiction =
			jurisdiction
	}

	return officer
}

func parseOfficerName(
	body string,
) string {
	match :=
		headingPattern.FindStringSubmatch(
			body,
		)

	if len(match) < 2 {

		return ""
	}

	return cleanText(
		match[1],
	)
}

func parseCompany(
	body string,
) (
	string,
	string,
	string,
) {
	match :=
		companyLinkPattern.FindStringSubmatch(
			body,
		)

	if len(match) < 4 {

		return "",
			"",
			""
	}

	path :=
		strings.TrimSpace(
			match[1],
		)

	jurisdiction :=
		strings.TrimSpace(
			match[2],
		)

	name :=
		cleanText(
			match[3],
		)

	if path == "" {

		return name,
			"",
			jurisdiction
	}

	return name,
		baseURL + path,
		jurisdiction
}

func fieldValue(
	body string,
	labels ...string,
) string {
	value :=
		fieldValueFromMatches(
			definitionPattern.FindAllStringSubmatch(
				body,
				-1,
			),
			labels,
		)

	if value != "" {

		return value
	}

	return fieldValueFromMatches(
		tableFieldPattern.FindAllStringSubmatch(
			body,
			-1,
		),
		labels,
	)
}

func fieldValueFromMatches(
	matches [][]string,
	labels []string,
) string {
	for _, match := range matches {

		if len(match) < 3 {

			continue
		}

		label :=
			normalizeLabel(
				cleanText(
					match[1],
				),
			)

		if label == "" {

			continue
		}

		for _, expected := range labels {

			if label !=
				normalizeLabel(
					expected,
				) {

				continue
			}

			return cleanText(
				match[2],
			)
		}
	}

	return ""
}

func normalizeLabel(
	value string,
) string {
	value =
		strings.TrimSpace(
			strings.ToLower(
				value,
			),
		)

	value =
		strings.ReplaceAll(
			value,
			"_",
			" ",
		)

	value =
		strings.TrimSuffix(
			value,
			":",
		)

	return whitespacePattern.ReplaceAllString(
		value,
		" ",
	)
}

func cleanText(
	value string,
) string {
	value =
		tagPattern.ReplaceAllString(
			value,
			" ",
		)

	value =
		html.UnescapeString(
			value,
		)

	value =
		whitespacePattern.ReplaceAllString(
			value,
			" ",
		)

	return strings.TrimSpace(
		value,
	)
}
