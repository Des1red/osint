package extract

import (
	"regexp"
	"strings"

	"osint/internal/enrich/model"
)

func extractRelationships(
	input model.Input,
) (
	[]model.OrganizationReference,
	[]model.EmploymentReference,
) {
	fullName :=
		strings.TrimSpace(
			input.FullName,
		)

	if fullName == "" {
		return nil,
			nil
	}

	namePattern :=
		relationshipNamePattern(
			fullName,
		)

	if namePattern == "" {
		return nil,
			nil
	}

	var organizations []model.OrganizationReference

	var employment []model.EmploymentReference

	for _, value := range input.Text {

		text :=
			strings.TrimSpace(
				value.Text,
			)

		if text == "" {
			continue
		}

		relationshipText :=
			relationshipEvidenceText(
				text,
			)

		segments :=
			relationshipSegments(
				relationshipText,
			)

		for _, segment := range segments {

			segment =
				strings.TrimSpace(
					segment,
				)

			if segment == "" {
				continue
			}

			foundOrganizations,
				foundEmployment :=
				relationshipsFromSegment(
					segment,
					value.Source,
					value.EvidenceID,
					namePattern,
				)

			organizations =
				append(
					organizations,
					foundOrganizations...,
				)

			employment =
				append(
					employment,
					foundEmployment...,
				)
		}
	}

	return uniqueOrganizations(
			organizations,
		),
		uniqueRelationshipEmployment(
			employment,
		)
}

func relationshipsFromSegment(
	text string,
	source string,
	evidenceID string,
	namePattern string,
) (
	[]model.OrganizationReference,
	[]model.EmploymentReference,
) {
	var organizations []model.OrganizationReference

	var employment []model.EmploymentReference

	add :=
		func(
			organization string,
			title string,
			current bool,
		) {
			organization =
				cleanOrganizationCandidate(
					organization,
				)

			title =
				normalizeRelationshipTitle(
					title,
				)

			if organization == "" {
				return
			}

			organizations =
				append(
					organizations,
					model.OrganizationReference{
						Organization: organization,

						Source: source,

						Evidence: wholeEvidenceReference(
							evidenceID,
						),
					},
				)

			employment =
				append(
					employment,
					model.EmploymentReference{
						Title: title,

						Organization: organization,

						Current: current,

						Source: source,

						Evidence: wholeEvidenceReference(
							evidenceID,
						),
					},
				)
		}

	passiveFounder :=
		regexp.MustCompile(
			`(?i)^(.+?)\s+(was\s+)?(co[- ]?founded|founded)\s+by\s+(.+)$`,
		)

	if match :=
		passiveFounder.FindStringSubmatch(
			text,
		); len(match) == 5 {

		founders :=
			match[4]

		if regexp.MustCompile(
			namePattern,
		).FindStringIndex(
			founders,
		) != nil {

			title :=
				"Founder"

			if strings.Contains(
				strings.ToLower(
					match[3],
				),
				"co",
			) {

				title =
					"Co-Founder"
			}

			add(
				match[1],
				title,
				true,
			)
		}
	}

	activeFounder :=
		regexp.MustCompile(
			`(?i)` +
				namePattern +
				`\s+(co[- ]?founded|founded)\s+(.+)$`,
		)

	if match :=
		activeFounder.FindStringSubmatch(
			text,
		); len(match) == 3 {

		title :=
			"Founder"

		if strings.Contains(
			strings.ToLower(
				match[1],
			),
			"co",
		) {

			title =
				"Co-Founder"
		}

		add(
			match[2],
			title,
			true,
		)
	}

	rolePattern :=
		regexp.MustCompile(
			`(?i)` +
				namePattern +
				`(?:\s*,)?\s+(?:is\s+|was\s+)?(?:the\s+)?` +
				`(co[- ]?founder|founder|co[- ]?owner|owner|chief executive officer|ceo|managing director|director|president|partner)` +
				`\s+(?:of|at)\s+(.+)$`,
		)

	if match :=
		rolePattern.FindStringSubmatch(
			text,
		); len(match) == 3 {

		current :=
			!strings.Contains(
				strings.ToLower(
					text,
				),
				" was ",
			)

		add(
			match[2],
			match[1],
			current,
		)
	}

	worksPattern :=
		regexp.MustCompile(
			`(?i)` +
				namePattern +
				`\s+(works|worked)\s+(?:at|for)\s+(.+)$`,
		)

	if match :=
		worksPattern.FindStringSubmatch(
			text,
		); len(match) == 3 {

		current :=
			strings.EqualFold(
				match[1],
				"works",
			)

		add(
			match[2],
			"",
			current,
		)
	}

	return organizations,
		employment
}

func relationshipNamePattern(
	fullName string,
) string {
	parts :=
		strings.Fields(
			strings.TrimSpace(
				fullName,
			),
		)

	if len(parts) < 2 {
		return ""
	}

	var escaped []string

	for _, part := range parts {

		part =
			strings.TrimSpace(
				part,
			)

		if part == "" {
			continue
		}

		escaped =
			append(
				escaped,
				regexp.QuoteMeta(
					part,
				),
			)
	}

	if len(escaped) < 2 {
		return ""
	}

	return `\b` +
		strings.Join(
			escaped,
			`\s+`,
		) +
		`\b`
}

func relationshipSegments(
	value string,
) []string {
	value =
		strings.ReplaceAll(
			value,
			"…",
			".",
		)

	return strings.FieldsFunc(
		value,
		func(
			character rune,
		) bool {
			switch character {

			case '.',
				'!',
				'?',
				'\n',
				'\r':

				return true

			default:
				return false
			}
		},
	)
}

func cleanOrganizationCandidate(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	separators :=
		[]string{
			" | ",
			" – ",
			" — ",
			" - ",
		}

	for _, separator := range separators {

		if index :=
			strings.LastIndex(
				value,
				separator,
			); index >= 0 {

			candidate :=
				strings.TrimSpace(
					value[index+
						len(
							separator,
						):],
				)

			if candidate != "" {
				value =
					candidate
			}
		}
	}

	for _, separator := range []string{
		", ",
		"; ",
		" who ",
		" which ",
		" where ",
	} {

		if index :=
			strings.Index(
				strings.ToLower(
					value,
				),
				separator,
			); index > 0 {

			value =
				value[:index]
		}
	}

	value =
		strings.Trim(
			strings.TrimSpace(
				value,
			),
			`"'“”‘’()[]{}:;-`,
		)

	if value == "" {
		return ""
	}

	if len(
		strings.Fields(
			value,
		),
	) > 12 {

		return ""
	}

	if len(value) > 120 {
		return ""
	}

	return value
}

func normalizeRelationshipTitle(
	value string,
) string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	switch value {

	case "founder":
		return "Founder"

	case "co-founder",
		"co founder":

		return "Co-Founder"

	case "owner":
		return "Owner"

	case "co-owner",
		"co owner":

		return "Co-Owner"

	case "ceo",
		"chief executive officer":

		return "CEO"

	case "managing director":
		return "Managing Director"

	case "director":
		return "Director"

	case "president":
		return "President"

	case "partner":
		return "Partner"

	default:
		return ""
	}
}

func uniqueRelationshipEmployment(
	values []model.EmploymentReference,
) []model.EmploymentReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.EmploymentReference

	for _, value := range values {

		if strings.TrimSpace(
			value.Organization,
		) == "" {

			continue
		}

		key :=
			strings.ToLower(
				strings.TrimSpace(
					value.Title,
				) +
					":" +
					strings.TrimSpace(
						value.Organization,
					) +
					":" +
					strings.TrimSpace(
						value.Source,
					),
			)

		if index,
			exists :=
			indexes[key]; exists {

			result[index].Evidence =
				mergeEvidence(
					result[index].Evidence,
					value.Evidence,
				)

			if !result[index].Current &&
				value.Current {

				result[index].Current =
					true
			}

			continue
		}

		indexes[key] =
			len(result)

		result =
			append(
				result,
				value,
			)
	}

	return result
}

func relationshipEvidenceText(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {

		return ""
	}

	//
	// URLs are not part of an organization
	// name.
	//
	// Replace them with a sentence boundary
	// rather than simply deleting them.
	//
	// Example:
	//
	// Person2 owner of company
	// www.straightrazorway.com Person2 ...
	//
	// becomes:
	//
	// Person2 owner of company.
	// Person2 ...
	//
	// This prevents:
	//
	//     company www
	//
	// from being extracted as the organization.
	//
	return linkPattern.ReplaceAllString(
		value,
		".",
	)
}
