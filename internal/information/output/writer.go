package output

import (
	"io"
	"os"

	"osint/internal/models"
)

var writer io.Writer = os.Stdout

var fileOutput bool

func SetWriter(
	value io.Writer,
	isFileOutput bool,
) {
	writer =
		value

	fileOutput =
		isFileOutput
}

func Writer() io.Writer {
	if !ColorEnabled() {
		return writer
	}

	return colorWriter{
		writer: writer,
	}
}

func FileOutput() bool {
	return fileOutput
}

func Full() bool {
	return models.OutputFlags.Full
}
