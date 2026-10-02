package xvfb

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readDisplay(
	readPipe *os.File,
) chan displayResult {
	ready :=
		make(
			chan displayResult,
			1,
		)

	go func() {
		reader :=
			bufio.NewReader(
				readPipe,
			)

		value,
			err :=
			reader.ReadString(
				'\n',
			)

		ready <- displayResult{
			Value: value,

			Err: err,
		}
	}()

	return ready
}

func setDisplay(
	value string,
) (
	string,
	error,
) {
	displayNumber :=
		strings.TrimSpace(
			value,
		)

	number,
		err :=
		strconv.Atoi(
			displayNumber,
		)

	if err != nil {

		return "",
			fmt.Errorf(
				"invalid Xvfb display number %q: %w",
				displayNumber,
				err,
			)
	}

	display :=
		fmt.Sprintf(
			":%d",
			number,
		)

	previousDisplay,
		previousDisplaySet =
		os.LookupEnv(
			"DISPLAY",
		)

	err =
		os.Setenv(
			"DISPLAY",
			display,
		)

	if err != nil {

		restoreDisplay()

		return "",
			fmt.Errorf(
				"set DISPLAY to %s: %w",
				display,
				err,
			)
	}

	return display,
		nil
}

func restoreDisplay() {
	if previousDisplaySet {

		_ =
			os.Setenv(
				"DISPLAY",
				previousDisplay,
			)

	} else {

		_ =
			os.Unsetenv(
				"DISPLAY",
			)
	}

	previousDisplay =
		""

	previousDisplaySet =
		false
}
