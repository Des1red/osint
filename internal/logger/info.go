package logger

import "fmt"

const (
	reset = "\033[0m"

	red = "\033[31m"

	green = "\033[32m"

	yellow = "\033[33m"

	cyan = "\033[36m"

	headerLine = "============================="
)

func Info(
	str string,
) {
	fmt.Printf(
		"[%sINFO%s] %s\n",
		cyan,
		reset,
		str,
	)
}

func HeaderStart(
	title string,
) {
	header(
		title,
		"started",
		green,
	)
}

func HeaderEnd(
	title string,
) {
	header(
		title,
		"ended",
		red,
	)
}

func header(
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
