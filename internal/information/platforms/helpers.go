package platforms

import (
	"fmt"
	"strings"

	"osint/internal/information/output"
)

func printHeader(
	name string,
) {
	fmt.Fprintln(
		output.Writer(),
	)

	fmt.Fprintln(
		output.Writer(),
		name,
	)

	fmt.Fprintln(
		output.Writer(),
		strings.Repeat(
			"-",
			len(name),
		),
	)
}

func printValue(
	name string,
	value string,
) {
	if strings.TrimSpace(
		value,
	) == "" {
		return
	}

	fmt.Fprintf(
		output.Writer(),
		"%s: %s\n",
		name,
		value,
	)
}

func printAccessState(
	accessible bool,
	loginRequired bool,
) {
	fmt.Fprintf(
		output.Writer(),
		"Accessible: %t\n",
		accessible,
	)

	fmt.Fprintf(
		output.Writer(),
		"Login Required: %t\n",
		loginRequired,
	)
}

func PrintNothingFound(
	platform string,
) {
	fmt.Fprintf(
		output.Writer(),
		"Nothing found on Platform: %s\n",
		platform,
	)
}
