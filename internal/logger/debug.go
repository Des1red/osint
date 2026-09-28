package logger

import (
	"encoding/json"
	"fmt"
	"strings"

	"osint/internal/models"
)

func Debug(
	title string,
	values ...any,
) {
	if !models.BootFlags.Debug {
		return
	}

	fmt.Println()
	fmt.Println(
		"========== DEBUG " +
			strings.ToUpper(
				strings.TrimSpace(
					title,
				),
			) +
			" ==========",
	)

	for index := 0; index < len(values); index += 2 {

		name :=
			"value"

		value :=
			values[index]

		if index+1 <
			len(values) {

			if key, ok :=
				values[index].(string); ok {

				name =
					key

				value =
					values[index+1]
			}
		}

		fmt.Println(
			name + ":",
		)

		printDebugValue(
			value,
		)
	}

	fmt.Println(
		"==========================================",
	)
}

func printDebugValue(
	value any,
) {
	encoded, err :=
		json.MarshalIndent(
			value,
			"",
			"  ",
		)

	if err == nil {

		fmt.Println(
			string(
				encoded,
			),
		)

		return
	}

	fmt.Printf(
		"%+v\n",
		value,
	)
}
