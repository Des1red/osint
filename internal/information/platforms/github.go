package platforms

import (
	"fmt"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func PrintGitHub(
	result knowledge.GitHubResult,
) {
	printHeader(
		"GitHub",
	)

	printValue(
		"Username",
		result.Username,
	)

	if result.ID != 0 {
		fmt.Fprintf(
			output.Writer(),
			"ID: %d\n",
			result.ID,
		)
	}

	printValue(
		"Node ID",
		result.NodeID,
	)

	printValue(
		"Type",
		result.Type,
	)

	fmt.Fprintf(
		output.Writer(),
		"Site Admin: %t\n",
		result.SiteAdmin,
	)

	printValue(
		"Name",
		result.Name,
	)

	printValue(
		"Company",
		result.Company,
	)

	printValue(
		"Location",
		result.Location,
	)

	printValue(
		"Email",
		result.Email,
	)

	printValue(
		"Bio",
		result.Bio,
	)

	printValue(
		"Twitter",
		result.Twitter,
	)

	printValue(
		"Blog",
		result.Blog,
	)

	fmt.Fprintf(
		output.Writer(),
		"Public Repositories: %d\n",
		result.PublicRepos,
	)

	fmt.Fprintf(
		output.Writer(),
		"Public Gists: %d\n",
		result.PublicGists,
	)

	fmt.Fprintf(
		output.Writer(),
		"Followers: %d\n",
		result.Followers,
	)

	fmt.Fprintf(
		output.Writer(),
		"Following: %d\n",
		result.Following,
	)

	printValue(
		"Created",
		result.CreatedAt,
	)

	printValue(
		"Updated",
		result.UpdatedAt,
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
		"API",
		result.APIURL,
	)

	printValue(
		"Source",
		result.Source,
	)
}
