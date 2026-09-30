package logger

import (
	"fmt"
	"strings"
)

const (
	reset = "\033[0m"

	red = "\033[31m"

	green = "\033[32m"

	yellow = "\033[33m"

	cyan = "\033[36m"

	orange = "\033[38;5;208m"

	headerLine = "============================="

	depthLine = "│   "
)

var sessionDepth int

func Info(
	str string,
) {
	fmt.Printf(
		"%s[%sINFO%s] %s\n",
		depthPrefix(),
		cyan,
		reset,
		str,
	)
}

func HeaderStart(
	title string,
) {
	if sessionDepth == 0 {

		rootHeader(
			title,
			"started",
			green,
		)

		sessionDepth++

		return
	}

	nestedHeaderStart(
		title,
	)

	sessionDepth++
}

func HeaderEnd(
	title string,
) {
	if sessionDepth > 0 {

		sessionDepth--
	}

	if sessionDepth == 0 {

		rootHeader(
			title,
			"ended",
			red,
		)

		return
	}

	nestedHeaderEnd(
		title,
	)
}

func rootHeader(
	title string,
	status string,
	statusColor string,
) {
	fmt.Printf(
		"%s%s%s\n",
		yellow,
		headerLine,
		reset,
	)

	fmt.Printf(
		"%s: %s%s%s\n",
		title,
		statusColor,
		status,
		reset,
	)

	fmt.Printf(
		"%s%s%s\n",
		yellow,
		headerLine,
		reset,
	)

	fmt.Println()
}

func nestedHeaderStart(
	title string,
) {
	parentIndent :=
		strings.Repeat(
			"    ",
			sessionDepth-1,
		)

	childIndent :=
		strings.Repeat(
			"    ",
			sessionDepth,
		)

	//
	// Continue the parent branch down to the
	// point where the child session begins.
	//
	fmt.Printf(
		"%s%s│%s\n",
		parentIndent,
		yellow,
		reset,
	)

	//
	// Orange transition into the child.
	//
	fmt.Printf(
		"%s%s└──%s\n",
		parentIndent,
		orange,
		reset,
	)

	//
	// The child header now lives at its own
	// depth.
	//
	fmt.Printf(
		"%s%s%s%s\n",
		childIndent,
		yellow,
		headerLine,
		reset,
	)

	fmt.Printf(
		"%s%s: %sstarted%s\n",
		childIndent,
		title,
		green,
		reset,
	)

	fmt.Printf(
		"%s%s%s%s\n",
		childIndent,
		yellow,
		headerLine,
		reset,
	)

	fmt.Printf(
		"%s%s│%s\n",
		childIndent,
		yellow,
		reset,
	)
}

func nestedHeaderEnd(
	title string,
) {
	indent :=
		strings.Repeat(
			"    ",
			sessionDepth,
		)

	fmt.Printf(
		"%s%s%s%s\n",
		indent,
		yellow,
		headerLine,
		reset,
	)

	fmt.Printf(
		"%s%s: %sended%s\n",
		indent,
		title,
		red,
		reset,
	)

	fmt.Printf(
		"%s%s%s%s\n",
		indent,
		yellow,
		headerLine,
		reset,
	)

	fmt.Println()
}

func depthPrefix() string {
	if sessionDepth == 0 {

		return ""
	}

	indent :=
		strings.Repeat(
			"    ",
			sessionDepth-1,
		)

	return indent +
		yellow +
		depthLine +
		reset
}
