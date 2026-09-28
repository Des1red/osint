package internal

import (
	"strings"

	"osint/internal/enrich/dedupe/internal/evidence"
	"osint/internal/enrich/dedupe/internal/normalize"
	"osint/internal/enrich/model"
)

func AccountResult(
	values []model.RelatedAccountReference,
) []model.RelatedAccountReference {
	var result []model.RelatedAccountReference

	for _, value := range values {

		username :=
			normalize.Text(
				value.Username,
			)

		if username == "" {

			continue
		}

		platform :=
			normalize.Text(
				value.Platform,
			)

		matched :=
			false

		for index := range result {

			existingUsername :=
				normalize.Text(
					result[index].Username,
				)

			if existingUsername !=
				username {

				continue
			}

			existingPlatform :=
				normalize.Text(
					result[index].Platform,
				)

			//
			// Same username on two explicitly
			// different platforms means two
			// different accounts.
			//
			if platform != "" &&
				existingPlatform != "" &&
				platform !=
					existingPlatform {

				continue
			}

			source :=
				evidence.PreferSource(
					result[index].Source,
					value.Source,
					result[index].Evidence,
					value.Evidence,
				)

			result[index].Evidence =
				evidence.Merge(
					result[index].Evidence,
					value.Evidence,
				)

			result[index].Source =
				source

			if existingPlatform == "" &&
				platform != "" {

				result[index].Platform =
					value.Platform
			}

			if strings.TrimSpace(
				result[index].ProfileURL,
			) == "" {

				result[index].ProfileURL =
					value.ProfileURL
			}

			if strings.TrimSpace(
				result[index].EvidenceURL,
			) == "" {

				result[index].EvidenceURL =
					value.EvidenceURL
			}

			//
			// Referenced is stronger than merely
			// Mentioned.
			//
			if strings.TrimSpace(
				result[index].Association,
			) == "" ||
				(strings.EqualFold(
					result[index].Association,
					"Mentioned",
				) &&
					strings.EqualFold(
						value.Association,
						"Referenced",
					)) {

				result[index].Association =
					value.Association
			}

			matched =
				true

			break
		}

		if matched {

			continue
		}

		result =
			append(
				result,
				value,
			)
	}

	return result
}
