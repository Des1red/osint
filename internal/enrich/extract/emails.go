package extract

import (
	"strings"

	"osint/internal/enrich/model"
)

func extractEmails(
	input model.Input,
) []model.EmailReference {
	var emails []model.EmailReference

	for _, value := range input.Text {

		matches :=
			emailPattern.FindAllStringIndex(
				value.Text,
				-1,
			)

		for _, match := range matches {

			if len(match) != 2 {

				continue
			}

			start :=
				match[0]

			end :=
				match[1]

			if start < 0 ||
				end < start ||
				end > len(value.Text) {

				continue
			}

			email :=
				strings.TrimSpace(
					value.Text[start:end],
				)

			if email == "" {

				continue
			}

			emails =
				append(
					emails,
					model.EmailReference{
						Email: email,

						Source: value.Source,

						Evidence: evidenceReference(
							value.EvidenceID,
							start,
							end,
						),
					},
				)
		}
	}

	return uniqueEmails(
		emails,
	)
}

func uniqueEmails(
	values []model.EmailReference,
) []model.EmailReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.EmailReference

	for _, value := range values {

		email :=
			strings.TrimSpace(
				value.Email,
			)

		if email == "" {

			continue
		}

		key :=
			strings.ToLower(
				email +
					":" +
					value.Source,
			)

		if index,
			exists :=
			indexes[key]; exists {

			//
			// Preserve any information learned
			// from another occurrence of the
			// same extracted email.
			//
			if result[index].Type == "" &&
				value.Type != "" {

				result[index].Type =
					value.Type
			}

			result[index].Evidence =
				mergeEvidence(
					result[index].Evidence,
					value.Evidence,
				)

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
