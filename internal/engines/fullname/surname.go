package fullname

import (
	"fmt"
	"strings"

	"osint/internal/models"
)

func resolveSurname(
	fullName string,
	position int,
) error {
	//
	// Position 0 means the user did not specify
	// the surname position.
	//
	// Leave the surname empty so surname-based
	// enrichment can simply skip its providers.
	//
	if position == 0 {

		models.ScopeInput.Surname =
			""

		return nil
	}

	fields :=
		strings.Fields(
			fullName,
		)

	if len(fields) < 2 {

		return fmt.Errorf(
			"full name requires at least two name parts",
		)
	}

	switch position {

	case 1:

		models.ScopeInput.Surname =
			fields[0]

	case 2:

		models.ScopeInput.Surname =
			fields[len(fields)-1]

	default:

		return fmt.Errorf(
			"surname position must be 1 or 2",
		)
	}

	return nil
}
