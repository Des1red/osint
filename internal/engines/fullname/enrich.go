package fullname

import (
	"osint/internal/enrich"
	"osint/internal/logger"
)

func enrichFullName(
	fullName string,
	matches []Match,
) enrich.EnrichmentResult {
	input :=
		enrich.Input{
			FullName: fullName,

			Accounts: make(
				[]enrich.AccountInput,
				0,
				len(matches),
			),
		}

	for _, match := range matches {

		input.Accounts =
			append(
				input.Accounts,
				enrich.AccountInput{
					Platform: match.Platform,

					Username: match.Username,

					ProfileURL: match.ProfileURL,
				},
			)
	}

	result, err :=
		enrich.Enrich(
			input,
		)

	if err != nil {

		logger.LogError(
			"Full name enrichment failed",
			err.Error(),
		)

		return enrich.EnrichmentResult{}
	}

	return result
}
