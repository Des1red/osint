package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"osint/internal/models"
)

func WriteFile(
	engine string,
	target string,
	write func(),
) error {
	file, err := createOutputFile(
		engine,
		target,
	)
	if err != nil {
		return err
	}

	defer file.Close()

	SetWriter(
		file,
		true,
	)

	write()

	SetWriter(
		os.Stdout,
		false,
	)

	return nil
}

func createOutputFile(
	engine string,
	target string,
) (*os.File, error) {
	base := strings.TrimSpace(
		models.ScopeInput.OutputFile,
	)

	if base == "" {
		base = fmt.Sprintf(
			"osint_%s_%s",
			engine,
			sanitizeFileName(target),
		)
	} else {
		extension := filepath.Ext(base)

		if extension != "" {
			base = strings.TrimSuffix(
				base,
				extension,
			)
		}

		base = fmt.Sprintf(
			"%s_%s",
			sanitizeFileName(base),
			engine,
		)
	}

	fileName := base + ".txt"

	return os.Create(
		fileName,
	)
}

func sanitizeFileName(
	value string,
) string {
	value = strings.TrimSpace(value)

	var builder strings.Builder

	for _, char := range value {
		switch {
		case unicode.IsLetter(char):
			builder.WriteRune(char)

		case unicode.IsDigit(char):
			builder.WriteRune(char)

		case char == '.',
			char == '-',
			char == '_':

			builder.WriteRune(char)

		default:
			builder.WriteRune('_')
		}
	}

	return filepath.Base(
		builder.String(),
	)
}
