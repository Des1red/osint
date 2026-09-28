package platforms

import (
	"fmt"

	"osint/internal/information/output"
	"osint/internal/knowledge"
)

func PrintReddit(
	result knowledge.RedditResult,
) {
	printHeader(
		"Reddit",
	)

	printValue(
		"Username",
		result.Username,
	)

	fmt.Fprintf(
		output.Writer(),
		"Post Karma: %d\n",
		result.PostKarma,
	)

	fmt.Fprintf(
		output.Writer(),
		"Comment Karma: %d\n",
		result.CommentKarma,
	)

	fmt.Fprintf(
		output.Writer(),
		"Total Karma: %d\n",
		result.TotalKarma,
	)

	printValue(
		"Reddit Age",
		result.RedditAge,
	)

	printValue(
		"Cake Day",
		result.CakeDay,
	)

	printValue(
		"Description",
		result.ProfileDescription,
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
