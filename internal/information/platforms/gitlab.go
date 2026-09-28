package platforms

import (
	"fmt"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func PrintGitLab(
	result knowledge.GitLabResult,
) {
	printHeader(
		"GitLab",
	)

	if result.ID != 0 {
		fmt.Fprintf(
			output.Writer(),
			"ID: %d\n",
			result.ID,
		)
	}

	printValue(
		"Username",
		result.Username,
	)

	printValue(
		"Name",
		result.Name,
	)

	printValue(
		"State",
		result.State,
	)

	fmt.Fprintf(
		output.Writer(),
		"Locked: %t\n",
		result.Locked,
	)

	printValue(
		"Avatar",
		result.AvatarURL,
	)

	printValue(
		"Profile",
		result.ProfileURL,
	)

	printValue(
		"Source",
		result.Source,
	)
}
