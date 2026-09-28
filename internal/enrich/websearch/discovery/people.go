package discovery

import (
	"strings"

	"osint/internal/enrich/extract"
	"osint/internal/enrich/model"
	"osint/internal/logger"
)

const maxRelationshipCandidates = 3

func validatePeopleRelationships(
	fullName string,
	people []model.PersonReference,
) []model.PersonReference {
	fullName =
		strings.TrimSpace(
			fullName,
		)

	if fullName == "" ||
		len(people) == 0 {

		return people
	}

	checked :=
		0

	for index := range people {

		if people[index].RelationVerified {
			continue
		}

		if checked >=
			maxRelationshipCandidates {

			break
		}

		candidate :=
			strings.TrimSpace(
				people[index].Name,
			)

		if candidate == "" {
			continue
		}

		checked++

		queries :=
			familyValidationQueries(
				fullName,
				candidate,
			)

		webResult, err :=
			search(
				queries,
			)

		if err != nil {

			logger.LogError(
				"Relationship validation failed for "+candidate,
				err.Error(),
			)

			continue
		}

		logger.Debug(
			"Relationship Validation",
			"target",
			fullName,
			"candidate",
			candidate,
			"results",
			webResult.Results,
		)

		for _, item := range webResult.Results {

			text :=
				strings.TrimSpace(
					item.Title +
						" " +
						item.Snippet,
				)

			if text == "" {
				continue
			}

			relation,
				verified :=
				extract.FamilyRelation(
					fullName,
					candidate,
					text,
				)

			if !verified {
				continue
			}

			people[index].Relation =
				relation

			people[index].RelationVerified =
				true

			people[index].RelationSource =
				relationshipValidationSource(
					item,
				)

			break
		}
	}

	return people
}

func familyValidationQueries(
	target string,
	candidate string,
) []string {
	target =
		strings.TrimSpace(
			target,
		)

	candidate =
		strings.TrimSpace(
			candidate,
		)

	if target == "" ||
		candidate == "" {

		return nil
	}

	return []string{
		`"` +
			target +
			`" "` +
			candidate +
			`" family`,

		`"` +
			target +
			`" "` +
			candidate +
			`" son OR daughter OR father OR mother OR brother OR sister`,
	}
}

func relationshipValidationSource(
	item SearchResult,
) string {
	source :=
		"WebSearch Relationship Validation"

	if strings.TrimSpace(
		item.Engine,
	) != "" {

		source +=
			" [" +
				strings.TrimSpace(
					item.Engine,
				) +
				"]"
	}

	if strings.TrimSpace(
		item.Query,
	) != "" {

		source +=
			" | Query: " +
				strings.TrimSpace(
					item.Query,
				)
	}

	if strings.TrimSpace(
		item.URL,
	) != "" {

		source +=
			" | URL: " +
				strings.TrimSpace(
					item.URL,
				)
	}

	return source
}
