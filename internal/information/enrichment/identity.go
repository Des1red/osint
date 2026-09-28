package enrichment

import (
	"strings"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func printIdentity(
	result knowledge.EnrichmentResult,
) {
	var sections []treeSection

	if len(result.Usernames) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Usernames",

					Print: func(
						prefix string,
					) {
						printUsernames(
							prefix,
							result.Usernames,
						)
					},
				},
			)
	}

	if len(
		result.UnattributedUsernames,
	) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Unattributed Usernames",

					Print: func(
						prefix string,
					) {
						printUsernames(
							prefix,
							result.UnattributedUsernames,
						)
					},
				},
			)
	}

	if len(result.RelatedAccounts) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Related Accounts",

					Print: func(
						prefix string,
					) {
						printRelatedAccounts(
							prefix,
							result.RelatedAccounts,
						)
					},
				},
			)
	}

	if len(
		result.UnattributedRelatedAccounts,
	) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Unattributed Related Accounts",

					Print: func(
						prefix string,
					) {
						printRelatedAccounts(
							prefix,
							result.UnattributedRelatedAccounts,
						)
					},
				},
			)
	}

	if len(result.Emails) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Emails",

					Print: func(
						prefix string,
					) {
						printEmails(
							prefix,
							result.Emails,
						)
					},
				},
			)
	}

	if len(result.UnattributedEmails) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Unattributed Emails",

					Print: func(
						prefix string,
					) {
						printEmails(
							prefix,
							result.UnattributedEmails,
						)
					},
				},
			)
	}

	if len(result.Socials) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Social Profiles",

					Print: func(
						prefix string,
					) {
						printSocials(
							prefix,
							result.Socials,
						)
					},
				},
			)
	}

	if len(
		result.UnattributedSocials,
	) > 0 {

		sections =
			append(
				sections,
				treeSection{
					Title: "Unattributed Social Profiles",

					Print: func(
						prefix string,
					) {
						printSocials(
							prefix,
							result.UnattributedSocials,
						)
					},
				},
			)
	}

	printTreeSection(
		"Identity",
		"--------",
		sections,
	)
}

func printUsernames(
	prefix string,
	usernames []knowledge.UsernameReference,
) {
	for index, username := range usernames {

		last :=
			index ==
				len(usernames)-1

		output.TreeItem(
			prefix,
			last,
			username.Username,
		)

		fields :=
			[]output.TreeField{}

		if username.Match !=
			knowledge.MatchNone {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Username Match",

						Value: matchLabel(
							username.Match,
						),
					},
				)
		}

		fields =
			appendProvenanceFields(
				fields,
				username.Source,
				username.Evidence,
			)

		output.TreeFields(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			fields,
		)
	}
}

func printRelatedAccounts(
	prefix string,
	accounts []knowledge.RelatedAccountReference,
) {
	for index, account := range accounts {

		last :=
			index ==
				len(accounts)-1

		title :=
			strings.TrimSpace(
				account.Username,
			)

		if title == "" {

			title =
				strings.TrimSpace(
					account.ProfileURL,
				)
		}

		if title == "" {

			title =
				strings.TrimSpace(
					account.Platform,
				)
		}

		if title == "" {

			title =
				"Account"
		}

		output.TreeItem(
			prefix,
			last,
			title,
		)

		fields :=
			[]output.TreeField{}

		if account.Platform != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Platform",

						Value: account.Platform,
					},
				)
		}

		if account.Association != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Association",

						Value: account.Association,
					},
				)
		}

		if account.ProfileURL != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Profile",

						Value: account.ProfileURL,
					},
				)
		}

		if account.EvidenceURL != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Evidence URL",

						Value: account.EvidenceURL,
					},
				)
		}

		fields =
			appendProvenanceFields(
				fields,
				account.Source,
				account.Evidence,
			)

		output.TreeFields(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			fields,
		)
	}
}

func printEmails(
	prefix string,
	emails []knowledge.EmailReference,
) {
	for index, email := range emails {

		last :=
			index ==
				len(emails)-1

		output.TreeItem(
			prefix,
			last,
			email.Email,
		)

		fields :=
			[]output.TreeField{}

		if email.Type != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Type",

						Value: email.Type,
					},
				)
		}

		fields =
			appendProvenanceFields(
				fields,
				email.Source,
				email.Evidence,
			)

		output.TreeFields(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			fields,
		)
	}
}

func printSocials(
	prefix string,
	socials []knowledge.SocialReference,
) {
	for index, social := range socials {

		last :=
			index ==
				len(socials)-1

		title :=
			strings.TrimSpace(
				social.Platform,
			)

		if title == "" {

			title =
				"Social Profile"
		}

		output.TreeItem(
			prefix,
			last,
			title,
		)

		fields :=
			[]output.TreeField{}

		if social.Username != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Username",

						Value: social.Username,
					},
				)
		}

		if social.Match !=
			knowledge.MatchNone {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "Username Match",

						Value: matchLabel(
							social.Match,
						),
					},
				)
		}

		if social.URL != "" {

			fields =
				append(
					fields,
					output.TreeField{
						Name: "URL",

						Value: social.URL,
					},
				)
		}

		fields =
			appendProvenanceFields(
				fields,
				social.Source,
				social.Evidence,
			)

		output.TreeFields(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			fields,
		)
	}
}

func matchLabel(
	level knowledge.MatchLevel,
) string {
	switch level {

	case knowledge.MatchExact:

		return "Exact"

	case knowledge.MatchClose:

		return "Close"

	case knowledge.MatchBroad:

		return "Broad"

	default:

		return "None"
	}
}
