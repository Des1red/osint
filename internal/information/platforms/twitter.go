package platforms

import (
	"osint/internal/knowledge"
)

func PrintTwitter(
	result knowledge.TwitterResult,
) {
	printHeader(
		"Twitter / X",
	)

	printAccessState(
		result.Accessible,
		result.LoginRequired,
	)

	printValue(
		"Username",
		result.Username,
	)

	printValue(
		"Name",
		result.Name,
	)

	printValue(
		"Description",
		result.Description,
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
