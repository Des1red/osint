package fullname

import (
	"fmt"
	"strings"
)

func validate(
	value string,
) error {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return fmt.Errorf(
			"full name is empty",
		)
	}

	fields :=
		strings.Fields(
			value,
		)

	if len(fields) < 2 {
		return fmt.Errorf(
			"full name requires at least first and last name",
		)
	}

	return nil
}
