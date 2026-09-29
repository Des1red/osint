package enrichment

import (
	"strings"

	"osint/internal/enrich/model"
	"osint/internal/information/output"
	"osint/internal/knowledge"
)

type socialCircleBranch struct {
	Name string

	Connections []model.SocialCircleConnection
}

func printSocialCircle(
	result knowledge.EnrichmentResult,
) {
	branches :=
		socialCircleTree(
			result.SocialCircle.Connections,
		)

	if len(branches) == 0 {

		return
	}

	printTreeSection(
		"Social Circle",
		"-------------",
		[]treeSection{
			{
				Title: "Connections",

				Print: func(
					prefix string,
				) {
					printSocialCircleBranches(
						prefix,
						branches,
					)
				},
			},
		},
	)
}

func socialCircleTree(
	connections []model.SocialCircleConnection,
) []socialCircleBranch {
	indexes :=
		make(
			map[string]int,
		)

	var result []socialCircleBranch

	for _, connection := range connections {

		from :=
			strings.TrimSpace(
				connection.From,
			)

		to :=
			strings.TrimSpace(
				connection.To,
			)

		if from == "" ||
			to == "" {

			continue
		}

		key :=
			strings.ToLower(
				strings.Join(
					strings.Fields(
						from,
					),
					" ",
				),
			)

		index,
			exists :=
			indexes[key]

		if !exists {

			index =
				len(result)

			indexes[key] =
				index

			result =
				append(
					result,
					socialCircleBranch{
						Name: from,
					},
				)
		}

		result[index].Connections =
			append(
				result[index].Connections,
				connection,
			)
	}

	return result
}

func printSocialCircleBranches(
	prefix string,
	branches []socialCircleBranch,
) {
	for index, branch := range branches {

		last :=
			index ==
				len(branches)-1

		output.TreeItem(
			prefix,
			last,
			branch.Name,
		)

		printSocialConnections(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			branch.Connections,
		)
	}
}

func printSocialConnections(
	prefix string,
	connections []model.SocialCircleConnection,
) {
	for index, connection := range connections {

		last :=
			index ==
				len(connections)-1

		title :=
			socialConnectionTitle(
				connection,
			)

		if title == "" {

			continue
		}

		output.TreeItem(
			prefix,
			last,
			title,
		)

		printSocialConnectionDetails(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			connection,
		)
	}
}

func printSocialConnectionDetails(
	prefix string,
	connection model.SocialCircleConnection,
) {
	fields :=
		[]output.TreeField{}

	if strings.TrimSpace(
		connection.Platform,
	) != "" {

		fields =
			append(
				fields,
				output.TreeField{
					Name: "Platform",

					Value: connection.Platform,
				},
			)
	}

	//
	// Only print Username separately when To was
	// resolved to a known person.
	//
	if strings.TrimSpace(
		connection.Username,
	) != "" &&
		!strings.EqualFold(
			strings.TrimSpace(
				connection.To,
			),
			"@"+
				strings.TrimSpace(
					connection.Username,
				),
		) {

		fields =
			append(
				fields,
				output.TreeField{
					Name: "Username",

					Value: connection.Username,
				},
			)
	}

	if strings.TrimSpace(
		connection.ProfileURL,
	) != "" {

		fields =
			append(
				fields,
				output.TreeField{
					Name: "Profile",

					Value: connection.ProfileURL,
				},
			)
	}

	hasPosts :=
		len(connection.Posts) > 0

	hasEvidence :=
		output.Full() &&
			len(connection.Evidence) > 0

	for index, field := range fields {

		last :=
			index ==
				len(fields)-1 &&
				!hasPosts &&
				!hasEvidence

		output.TreeValue(
			prefix,
			last,
			field.Name,
			field.Value,
		)
	}

	if hasPosts {

		last :=
			!hasEvidence

		output.TreeItem(
			prefix,
			last,
			"Posts",
		)

		printSocialPosts(
			output.TreeChildPrefix(
				prefix,
				last,
			),
			connection.Posts,
		)
	}

	if hasEvidence {

		output.TreeValue(
			prefix,
			true,
			"Evidence",
			evidenceSummary(
				connection.Evidence,
			),
		)
	}
}

func printSocialPosts(
	prefix string,
	posts []string,
) {
	filtered :=
		make(
			[]string,
			0,
			len(posts),
		)

	for _, post := range posts {

		post =
			strings.TrimSpace(
				post,
			)

		if post == "" {

			continue
		}

		filtered =
			append(
				filtered,
				post,
			)
	}

	for index, post := range filtered {

		output.TreeItem(
			prefix,
			index ==
				len(filtered)-1,
			post,
		)
	}
}

func socialConnectionTitle(
	connection model.SocialCircleConnection,
) string {
	to :=
		strings.TrimSpace(
			connection.To,
		)

	username :=
		strings.TrimSpace(
			connection.Username,
		)

	if to == "" {

		if username == "" {

			return ""
		}

		return "@" +
			username
	}

	//
	// Unresolved account:
	//
	//     @giorgosmanikas
	//
	if strings.HasPrefix(
		to,
		"@",
	) {

		return to
	}

	//
	// Known identity:
	//
	//     Giorgos Manikas (@giorgosmanikas)
	//
	if username != "" {

		return to +
			" (@" +
			username +
			")"
	}

	return to
}
