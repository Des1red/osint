package username

import (
	"fmt"
	"strings"
)

func validate(
	username string,
) error {
	username =
		strings.TrimSpace(
			username,
		)

	if username == "" {
		return fmt.Errorf(
			"username is empty",
		)
	}

	return nil
}
