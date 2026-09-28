package output

import (
	"fmt"
	"strings"
)

const (
	treeBranch = "├─ "

	treeLast = "└─ "

	treePipe = "│  "

	treeSpace = "   "
)

type TreeField struct {
	Name string

	Value string
}

func TreeItem(
	prefix string,
	last bool,
	value string,
) {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return
	}

	branch :=
		treeBranch

	if last {
		branch =
			treeLast
	}

	fmt.Fprintf(
		Writer(),
		"%s%s%s\n",
		prefix,
		branch,
		value,
	)
}

func TreeValue(
	prefix string,
	last bool,
	name string,
	value string,
) {
	name =
		strings.TrimSpace(
			name,
		)

	value =
		strings.TrimSpace(
			value,
		)

	if name == "" ||
		value == "" {

		return
	}

	branch :=
		treeBranch

	if last {
		branch =
			treeLast
	}

	fmt.Fprintf(
		Writer(),
		"%s%s%s: %s\n",
		prefix,
		branch,
		name,
		value,
	)
}

func TreeFields(
	prefix string,
	fields []TreeField,
) {
	filtered :=
		make(
			[]TreeField,
			0,
			len(fields),
		)

	for _, field := range fields {

		field.Name =
			strings.TrimSpace(
				field.Name,
			)

		field.Value =
			strings.TrimSpace(
				field.Value,
			)

		if field.Name == "" ||
			field.Value == "" {

			continue
		}

		filtered =
			append(
				filtered,
				field,
			)
	}

	for index, field := range filtered {

		TreeValue(
			prefix,
			index ==
				len(filtered)-1,
			field.Name,
			field.Value,
		)
	}
}

func TreeChildPrefix(
	prefix string,
	last bool,
) string {
	if last {

		return prefix +
			treeSpace
	}

	return prefix +
		treePipe
}
