package output

import (
	"bytes"
	"io"
	"os"
	"strings"
)

const (
	ansiReset = "\033[0m"

	ansiBold = "\033[1m"

	ansiDim = "\033[2m"

	ansiCyan = "\033[36m"

	ansiBlue = "\033[34m"

	ansiYellow = "\033[33m"
)

type colorWriter struct {
	writer io.Writer
}

func (
	w colorWriter,
) Write(
	data []byte,
) (
	int,
	error,
) {
	if !ColorEnabled() {

		return w.writer.Write(
			data,
		)
	}

	colored :=
		colorOutput(
			string(
				data,
			),
		)

	_,
		err :=
		w.writer.Write(
			[]byte(
				colored,
			),
		)

	if err != nil {
		return 0,
			err
	}

	//
	// io.Writer must report the number of
	// original input bytes consumed, not the
	// size after ANSI codes were added.
	//
	return len(data),
		nil
}

func ColorEnabled() bool {
	if FileOutput() {
		return false
	}

	if _, disabled :=
		os.LookupEnv(
			"NO_COLOR",
		); disabled {

		return false
	}

	if strings.EqualFold(
		strings.TrimSpace(
			os.Getenv(
				"TERM",
			),
		),
		"dumb",
	) {

		return false
	}

	info,
		err :=
		os.Stdout.Stat()

	if err != nil {
		return false
	}

	return info.Mode()&
		os.ModeCharDevice != 0
}

func colorOutput(
	value string,
) string {
	if value == "" {
		return value
	}

	//
	// Preserve multi-line writes.
	//
	lines :=
		strings.SplitAfter(
			value,
			"\n",
		)

	var buffer bytes.Buffer

	for _, line := range lines {

		if line == "" {
			continue
		}

		buffer.WriteString(
			colorLine(
				line,
			),
		)
	}

	return buffer.String()
}

func colorLine(
	line string,
) string {
	hasNewline :=
		strings.HasSuffix(
			line,
			"\n",
		)

	content :=
		strings.TrimSuffix(
			line,
			"\n",
		)

	if strings.TrimSpace(
		content,
	) == "" {

		return line
	}

	colored :=
		colorContent(
			content,
		)

	if hasNewline {
		colored +=
			"\n"
	}

	return colored
}

func colorContent(
	line string,
) string {
	trimmed :=
		strings.TrimSpace(
			line,
		)

	//
	// Main title underlines.
	//
	// =====================
	//
	if underline(
		trimmed,
		'=',
	) {

		return paint(
			line,
			ansiBold,
			ansiCyan,
		)
	}

	//
	// Section underlines.
	//
	// ---------------------
	//
	if underline(
		trimmed,
		'-',
	) {

		return paint(
			line,
			ansiBlue,
		)
	}

	//
	// Tree output.
	//
	if strings.Contains(
		line,
		"├─ ",
	) ||
		strings.Contains(
			line,
			"└─ ",
		) {

		return colorTreeLine(
			line,
		)
	}

	//
	// Vertical tree continuation.
	//
	if strings.Contains(
		line,
		"│  ",
	) {

		return colorTreeLine(
			line,
		)
	}

	return line
}

func colorTreeLine(
	line string,
) string {
	branchIndex :=
		treeBranchIndex(
			line,
		)

	if branchIndex < 0 {
		return line
	}

	branchLength :=
		len(
			"├─ ",
		)

	if strings.HasPrefix(
		line[branchIndex:],
		"└─ ",
	) {

		branchLength =
			len(
				"└─ ",
			)
	}

	before :=
		line[:branchIndex]

	branch :=
		line[branchIndex : branchIndex+
			branchLength]

	after :=
		line[branchIndex+
			branchLength:]

	before =
		colorTreePrefix(
			before,
		)

	branch =
		paint(
			branch,
			ansiDim,
			ansiCyan,
		)

	//
	// Field:
	//
	// ├─ Country: Greece
	//
	colon :=
		strings.Index(
			after,
			": ",
		)

	if colon >= 0 {

		name :=
			after[:colon]

		value :=
			after[colon+2:]

		return before +
			branch +
			paint(
				name,
				ansiYellow,
			) +
			": " +
			value
	}

	//
	// Tree node.
	//
	// ├─ Network
	//
	return before +
		branch +
		paint(
			after,
			ansiBlue,
		)
}

func colorTreePrefix(
	value string,
) string {
	value =
		strings.ReplaceAll(
			value,
			"│  ",
			paint(
				"│  ",
				ansiDim,
				ansiCyan,
			),
		)

	return value
}

func treeBranchIndex(
	value string,
) int {
	branch :=
		strings.Index(
			value,
			"├─ ",
		)

	last :=
		strings.Index(
			value,
			"└─ ",
		)

	switch {

	case branch < 0:
		return last

	case last < 0:
		return branch

	case branch < last:
		return branch

	default:
		return last
	}
}

func underline(
	value string,
	char rune,
) bool {
	if value == "" {
		return false
	}

	for _, current := range value {

		if current != char {
			return false
		}
	}

	return len(
		value,
	) >= 3
}

func paint(
	value string,
	codes ...string,
) string {
	if value == "" {
		return value
	}

	return strings.Join(
		codes,
		"",
	) +
		value +
		ansiReset
}
