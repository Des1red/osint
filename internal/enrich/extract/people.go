package extract

import (
	"regexp"
	"strings"
	"unicode"

	"osint/internal/enrich/model"
)

const personNameCapture = `([\p{L}][\p{L}'’.-]*(?:\s+[\p{L}][\p{L}'’.-]*){1,3})`

const familyRelationCapture = `(son|daughter|father|mother|brother|sister|husband|wife|spouse|partner|parent|child)`

func extractPeople(
	input model.Input,
) []model.PersonReference {
	target :=
		strings.TrimSpace(
			input.FullName,
		)

	if target == "" {
		return nil
	}

	var people []model.PersonReference

	for _, value := range input.Text {

		text :=
			strings.TrimSpace(
				value.Text,
			)

		if text == "" {
			continue
		}

		segments :=
			relationshipSegments(
				text,
			)

		for _, segment := range segments {

			segment =
				strings.TrimSpace(
					segment,
				)

			if segment == "" {
				continue
			}

			people =
				append(
					people,
					peopleFromFounderStatement(
						target,
						segment,
						value.Source,
						value.EvidenceID,
					)...,
				)

			people =
				append(
					people,
					peopleFromExplicitFamilyStatement(
						target,
						segment,
						value.Source,
						value.EvidenceID,
					)...,
				)
		}
	}

	return uniquePeople(
		people,
	)
}

func peopleFromFounderStatement(
	target string,
	text string,
	source string,
	evidenceID string,
) []model.PersonReference {
	var people []model.PersonReference

	passive :=
		regexp.MustCompile(
			`(?i)(?:co[- ]?founded|founded)\s+by\s+(.+)$`,
		)

	if match :=
		passive.FindStringSubmatch(
			text,
		); len(match) == 2 {

		people =
			append(
				people,
				peopleFromList(
					target,
					match[1],
					source,
					evidenceID,
				)...,
			)
	}

	active :=
		regexp.MustCompile(
			`(?i)^(.+?)\s+(?:co[- ]?founded|founded)\s+`,
		)

	if match :=
		active.FindStringSubmatch(
			text,
		); len(match) == 2 {

		people =
			append(
				people,
				peopleFromList(
					target,
					match[1],
					source,
					evidenceID,
				)...,
			)
	}

	return people
}

func peopleFromList(
	target string,
	value string,
	source string,
	evidenceID string,
) []model.PersonReference {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return nil
	}

	if !personNamePresent(
		value,
		target,
	) {

		return nil
	}

	replacer :=
		strings.NewReplacer(
			" and ",
			",",

			" & ",
			",",

			" + ",
			",",

			";",
			",",
		)

	parts :=
		strings.Split(
			replacer.Replace(
				value,
			),
			",",
		)

	var people []model.PersonReference

	for _, part := range parts {

		name :=
			cleanPersonName(
				part,
			)

		if name == "" {
			continue
		}

		if samePersonName(
			name,
			target,
		) {

			continue
		}

		if !validPersonName(
			name,
		) {

			continue
		}

		people =
			append(
				people,
				model.PersonReference{
					Name: name,

					Source: source,

					Evidence: wholeEvidenceReference(
						evidenceID,
					),
				},
			)
	}

	return people
}

func peopleFromExplicitFamilyStatement(
	target string,
	text string,
	source string,
	evidenceID string,
) []model.PersonReference {
	targetPattern :=
		flexiblePersonPattern(
			target,
		)

	if targetPattern == "" {
		return nil
	}

	var people []model.PersonReference

	possessive :=
		regexp.MustCompile(
			`(?i)` +
				targetPattern +
				`(?:'s|’s)\s+` +
				familyRelationCapture +
				`\s+` +
				personNameCapture,
		)

	if match :=
		possessive.FindStringSubmatch(
			text,
		); len(match) == 3 {

		addVerifiedPerson(
			&people,
			target,
			match[2],
			match[1],
			source,
			evidenceID,
		)
	}

	pronoun :=
		regexp.MustCompile(
			`(?i)` +
				targetPattern +
				`\s+and\s+(?:his|her|their)\s+` +
				familyRelationCapture +
				`\s+` +
				personNameCapture,
		)

	if match :=
		pronoun.FindStringSubmatch(
			text,
		); len(match) == 3 {

		addVerifiedPerson(
			&people,
			target,
			match[2],
			match[1],
			source,
			evidenceID,
		)
	}

	candidateFirst :=
		regexp.MustCompile(
			`(?i)` +
				personNameCapture +
				`\s*,?\s*(?:is\s+|was\s+)?(?:the\s+)?` +
				familyRelationCapture +
				`\s+of\s+` +
				targetPattern,
		)

	if match :=
		candidateFirst.FindStringSubmatch(
			text,
		); len(match) == 3 {

		addVerifiedPerson(
			&people,
			target,
			match[1],
			match[2],
			source,
			evidenceID,
		)
	}

	return people
}

func addVerifiedPerson(
	people *[]model.PersonReference,
	target string,
	name string,
	relation string,
	source string,
	evidenceID string,
) {
	name =
		cleanPersonName(
			name,
		)

	if name == "" ||
		samePersonName(
			name,
			target,
		) ||
		!validPersonName(
			name,
		) {

		return
	}

	relation =
		normalizeFamilyRelation(
			relation,
		)

	if relation == "" {
		return
	}

	*people =
		append(
			*people,
			model.PersonReference{
				Name: name,

				Relation: relation,

				RelationVerified: true,

				Source: source,

				Evidence: wholeEvidenceReference(
					evidenceID,
				),

				RelationSource: source,

				RelationEvidence: wholeEvidenceReference(
					evidenceID,
				),
			},
		)
}

// FamilyRelation determines whether a piece of
// evidence explicitly states the relationship
// between target and candidate.
//
// It does NOT infer from surname, company,
// co-occurrence, age, gender, or anything else.
func FamilyRelation(
	target string,
	candidate string,
	text string,
) (
	string,
	bool,
) {
	target =
		strings.TrimSpace(
			target,
		)

	candidate =
		strings.TrimSpace(
			candidate,
		)

	text =
		strings.TrimSpace(
			text,
		)

	if target == "" ||
		candidate == "" ||
		text == "" {

		return "",
			false
	}

	targetPattern :=
		flexiblePersonPattern(
			target,
		)

	candidatePattern :=
		flexiblePersonPattern(
			candidate,
		)

	if targetPattern == "" ||
		candidatePattern == "" {

		return "",
			false
	}

	candidateFirst :=
		regexp.MustCompile(
			`(?i)` +
				candidatePattern +
				`\s*,?\s*(?:is\s+|was\s+)?(?:the\s+)?` +
				familyRelationCapture +
				`\s+of\s+` +
				targetPattern,
		)

	if match :=
		candidateFirst.FindStringSubmatch(
			text,
		); len(match) == 2 {

		relation :=
			normalizeFamilyRelation(
				match[1],
			)

		if relation != "" {
			return relation,
				true
		}
	}

	possessive :=
		regexp.MustCompile(
			`(?i)` +
				targetPattern +
				`(?:'s|’s)\s+` +
				familyRelationCapture +
				`\s+` +
				candidatePattern,
		)

	if match :=
		possessive.FindStringSubmatch(
			text,
		); len(match) == 2 {

		relation :=
			normalizeFamilyRelation(
				match[1],
			)

		if relation != "" {
			return relation,
				true
		}
	}

	pronoun :=
		regexp.MustCompile(
			`(?i)` +
				targetPattern +
				`\s+and\s+(?:his|her|their)\s+` +
				familyRelationCapture +
				`\s+` +
				candidatePattern,
		)

	if match :=
		pronoun.FindStringSubmatch(
			text,
		); len(match) == 2 {

		relation :=
			normalizeFamilyRelation(
				match[1],
			)

		if relation != "" {
			return relation,
				true
		}
	}

	targetFirst :=
		regexp.MustCompile(
			`(?i)` +
				targetPattern +
				`\s*,?\s*(?:is\s+|was\s+)?(?:the\s+)?` +
				familyRelationCapture +
				`\s+of\s+` +
				candidatePattern,
		)

	if match :=
		targetFirst.FindStringSubmatch(
			text,
		); len(match) == 2 {

		relation :=
			inverseFamilyRelation(
				match[1],
			)

		if relation != "" {
			return relation,
				true
		}
	}

	return "",
		false
}

func normalizeFamilyRelation(
	value string,
) string {
	switch strings.ToLower(
		strings.TrimSpace(
			value,
		),
	) {

	case "son":
		return "Son"

	case "daughter":
		return "Daughter"

	case "father":
		return "Father"

	case "mother":
		return "Mother"

	case "brother":
		return "Brother"

	case "sister":
		return "Sister"

	case "husband":
		return "Husband"

	case "wife":
		return "Wife"

	case "spouse":
		return "Spouse"

	case "partner":
		return "Partner"

	case "parent":
		return "Parent"

	case "child":
		return "Child"

	default:
		return ""
	}
}

func inverseFamilyRelation(
	value string,
) string {
	switch strings.ToLower(
		strings.TrimSpace(
			value,
		),
	) {

	case "father",
		"mother",
		"parent":

		return "Child"

	case "son",
		"daughter",
		"child":

		return "Parent"

	case "brother",
		"sister":

		return "Sibling"

	case "husband",
		"wife",
		"spouse":

		return "Spouse"

	case "partner":
		return "Partner"

	default:
		return ""
	}
}

func flexiblePersonPattern(
	value string,
) string {
	parts :=
		strings.Fields(
			strings.TrimSpace(
				value,
			),
		)

	if len(parts) < 2 {
		return ""
	}

	build :=
		func(
			values []string,
		) string {

			var escaped []string

			for _, value := range values {

				value =
					strings.TrimSpace(
						value,
					)

				if value == "" {
					continue
				}

				escaped =
					append(
						escaped,
						regexp.QuoteMeta(
							value,
						),
					)
			}

			return strings.Join(
				escaped,
				`\s+`,
			)
		}

	forward :=
		build(
			parts,
		)

	if len(parts) != 2 {
		return forward
	}

	reversed :=
		build(
			[]string{
				parts[1],
				parts[0],
			},
		)

	return `(?:` +
		forward +
		`|` +
		reversed +
		`)`
}

func personNamePresent(
	value string,
	name string,
) bool {
	pattern :=
		flexiblePersonPattern(
			name,
		)

	if pattern == "" {
		return false
	}

	return regexp.MustCompile(
		`(?i)`+
			pattern,
	).FindStringIndex(
		value,
	) != nil
}

func samePersonName(
	left string,
	right string,
) bool {
	leftParts :=
		personTokens(
			left,
		)

	rightParts :=
		personTokens(
			right,
		)

	if len(leftParts) !=
		len(rightParts) {

		return false
	}

	if len(leftParts) == 0 {
		return false
	}

	used :=
		make(
			[]bool,
			len(rightParts),
		)

	for _, leftPart := range leftParts {

		found :=
			false

		for index, rightPart := range rightParts {

			if used[index] {
				continue
			}

			if leftPart !=
				rightPart {

				continue
			}

			used[index] =
				true

			found =
				true

			break
		}

		if !found {
			return false
		}
	}

	return true
}

func personTokens(
	value string,
) []string {
	return strings.FieldsFunc(
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		),
		func(
			character rune,
		) bool {

			return !unicode.IsLetter(
				character,
			)
		},
	)
}

func cleanPersonName(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	value =
		strings.Trim(
			value,
			`"'“”‘’()[]{}:;,.`,
		)

	return strings.Join(
		strings.Fields(
			value,
		),
		" ",
	)
}

func validPersonName(
	value string,
) bool {
	parts :=
		strings.Fields(
			strings.TrimSpace(
				value,
			),
		)

	if len(parts) < 2 ||
		len(parts) > 4 {

		return false
	}

	for _, part := range parts {

		part =
			strings.Trim(
				part,
				`"'“”‘’()[]{}:;,.`,
			)

		if part == "" {
			return false
		}

		first,
			_ :=
			utf8FirstRune(
				part,
			)

		if first == 0 ||
			!unicode.IsUpper(
				first,
			) {

			return false
		}
	}

	return true
}

func utf8FirstRune(
	value string,
) (
	rune,
	int,
) {
	for _, character := range value {

		return character,
			1
	}

	return 0,
		0
}

func uniquePeople(
	values []model.PersonReference,
) []model.PersonReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.PersonReference

	for _, value := range values {

		name :=
			cleanPersonName(
				value.Name,
			)

		if name == "" {
			continue
		}

		key :=
			strings.ToLower(
				name,
			)

		index,
			exists :=
			indexes[key]

		if !exists {

			value.Name =
				name

			indexes[key] =
				len(result)

			result =
				append(
					result,
					value,
				)

			continue
		}

		result[index].Evidence =
			mergeEvidence(
				result[index].Evidence,
				value.Evidence,
			)

		if strings.TrimSpace(
			result[index].Source,
		) == "" {

			result[index].Source =
				value.Source
		}

		if !result[index].RelationVerified &&
			value.RelationVerified {

			result[index].Relation =
				value.Relation

			result[index].RelationVerified =
				true

			result[index].RelationSource =
				value.RelationSource
		}

		result[index].RelationEvidence =
			mergeEvidence(
				result[index].RelationEvidence,
				value.RelationEvidence,
			)

		if strings.TrimSpace(
			result[index].RelationSource,
		) == "" {

			result[index].RelationSource =
				value.RelationSource
		}
	}

	return result
}
