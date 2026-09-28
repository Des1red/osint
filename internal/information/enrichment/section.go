package enrichment

import (
	"fmt"

	"osint/internal/information/output"
)

type treeSection struct {
	Title string

	Print func(
		prefix string,
	)
}

func printTreeSection(
	title string,
	underline string,
	sections []treeSection,
) {
	if len(sections) == 0 {
		return
	}

	fmt.Fprintln(
		output.Writer(),
	)

	fmt.Fprintln(
		output.Writer(),
		title,
	)

	fmt.Fprintln(
		output.Writer(),
		underline,
	)

	for index, section := range sections {

		last :=
			index ==
				len(sections)-1

		output.TreeItem(
			"",
			last,
			section.Title,
		)

		section.Print(
			output.TreeChildPrefix(
				"",
				last,
			),
		)
	}
}
