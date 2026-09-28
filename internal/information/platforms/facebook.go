package platforms

import (
	"osint/internal/knowledge"
)

func PrintFacebook(
	result knowledge.FacebookResult,
) {
	printHeader(
		"Facebook",
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
		"Category",
		result.Category,
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
