package fullname

import (
	"fmt"

	enrichmentinfo "osint/internal/information/enrichment"
	"osint/internal/information/output"
	platforminfo "osint/internal/information/platforms"
	"osint/internal/knowledge"
)

func FullNamePrinter() bool {
	result :=
		knowledge.Data.FullName

	printMatches(
		result.Matches,
	)

	hasEnrichment :=
		enrichmentinfo.Print(
			result.Enrichment,
		)

	return len(
		result.Matches,
	) > 0 ||
		hasEnrichment
}

func printMatches(
	matches []knowledge.FullNameMatch,
) {
	fmt.Fprintln(
		output.Writer(),
		"Full Name Matches",
	)

	fmt.Fprintln(
		output.Writer(),
		"-----------------",
	)

	found :=
		make(
			map[string]bool,
		)

	for index, match := range matches {

		found[match.Platform] =
			true

		last :=
			index ==
				len(matches)-1

		label :=
			match.Platform

		if label == "" {
			label =
				"Profile"
		}

		output.TreeItem(
			"",
			last,
			label,
		)

		output.TreeFields(
			output.TreeChildPrefix(
				"",
				last,
			),
			[]output.TreeField{
				{
					Name: "Username",

					Value: match.Username,
				},
				{
					Name: "Name",

					Value: match.Name,
				},
				{
					Name: "Profile",

					Value: match.ProfileURL,
				},
				{
					Name: "Description",

					Value: match.Description,
				},
			},
		)
	}

	fmt.Fprintln(
		output.Writer(),
	)

	if !found["GitHub"] {

		platforminfo.PrintNothingFound(
			"GitHub",
		)
	}

	if !found["GitLab"] {

		platforminfo.PrintNothingFound(
			"GitLab",
		)
	}

	if !found["Facebook"] {

		platforminfo.PrintNothingFound(
			"Facebook",
		)
	}
}
