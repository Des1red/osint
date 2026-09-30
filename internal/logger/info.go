package logger

import "fmt"

const headerLine = "============================="

func Info(
	str string,
) {
	fmt.Println(
		"[INFO] " +
			str,
	)
}

func HeaderStart(
	title string,
) {
	header(
		title,
		"started",
	)
}

func HeaderEnd(
	title string,
) {
	header(
		title,
		"ended",
	)
}

func header(
	title string,
	status string,
) {
	fmt.Println(
		headerLine,
	)

	fmt.Printf(
		"%s: %s\n",
		title,
		status,
	)

	fmt.Println(
		headerLine,
	)

	fmt.Println()
}
