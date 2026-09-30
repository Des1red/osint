package accounts

import (
	"strings"

	identitymatch "osint/internal/enrich/match"
)

func identityPresent(
	value string,
	identity string,
) bool {
	return identitymatch.IdentityPresent(
		identity,
		value,
	)
}

func urlSearchText(
	value string,
) string {
	replacer :=
		strings.NewReplacer(
			"/",
			" ",

			"-",
			" ",

			"_",
			" ",

			".",
			" ",

			"%20",
			" ",
		)

	return strings.ToLower(
		replacer.Replace(
			value,
		),
	)
}
