package accounts

import (
	"regexp"
	"strings"

	"osint/internal/enrich/model"
)

var mentionedAccountPattern = regexp.MustCompile(
	`@([A-Za-z0-9._-]{1,64})`,
)

func relatedAccountsFromEvidence(
	values []model.Evidence,
	owned ownedAccounts,
) []model.RelatedAccountReference {
	var result []model.RelatedAccountReference

	for _, value := range values {

		//
		// Do not create account relationships
		// from follower/following directories or
		// unknown mirror sites.
		//
		if !accountAssociationEvidence(
			value,
		) {

			continue
		}

		text :=
			value.Text

		if strings.TrimSpace(
			text,
		) == "" {

			continue
		}

		matches :=
			mentionedAccountPattern.FindAllStringSubmatchIndex(
				text,
				-1,
			)

		if len(matches) == 0 {

			continue
		}

		platform :=
			platformFromURL(
				value.URL,
			)

		if platform == "" {

			continue
		}

		for _, match := range matches {

			if len(match) != 4 {

				continue
			}

			start :=
				match[0]

			end :=
				match[1]

			usernameStart :=
				match[2]

			usernameEnd :=
				match[3]

			if start < 0 ||
				end < start ||
				usernameStart < 0 ||
				usernameEnd < usernameStart ||
				end > len(text) ||
				usernameEnd > len(text) {

				continue
			}

			username :=
				cleanAccountMention(
					text[usernameStart:usernameEnd],
					text[end:],
				)

			if username == "" ||
				owned.hasUsername(
					username,
				) {

				continue
			}

			evidenceEnd :=
				usernameStart +
					len(username)

			if evidenceEnd >
				end {

				evidenceEnd =
					end
			}

			result =
				append(
					result,
					model.RelatedAccountReference{
						Platform: platform,

						Username: username,

						EvidenceURL: strings.TrimSpace(
							value.URL,
						),

						Association: "Mentioned",

						Source: evidenceSource(
							value,
						),

						Evidence: evidenceReference(
							value,
							start,
							evidenceEnd,
						),
					},
				)
		}
	}

	return result
}

func evidenceSource(
	value model.Evidence,
) string {
	source :=
		strings.TrimSpace(
			value.Source,
		)

	if source != "" {

		return source
	}

	source =
		"Enrichment Evidence"

	if strings.TrimSpace(
		value.Engine,
	) != "" {

		source +=
			" [" +
				strings.TrimSpace(
					value.Engine,
				) +
				"]"
	}

	if strings.TrimSpace(
		value.Query,
	) != "" {

		source +=
			" | Query: " +
				strings.TrimSpace(
					value.Query,
				)
	}

	return source
}
