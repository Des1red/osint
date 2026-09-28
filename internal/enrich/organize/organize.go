package organize

import (
	"osint/internal/enrich/model"
	"osint/internal/enrich/organize/internal/accounts"
	"osint/internal/enrich/organize/internal/phones"
	"osint/internal/enrich/organize/internal/subjects"
)

func Organize(
	result model.EnrichmentResult,
	input model.Input,
	evidence []model.Evidence,
) model.EnrichmentResult {
	//
	// First route generic facts according to
	// Evidence.Subjects.
	//
	// This establishes the broad semantic owner:
	//
	//     root target
	//     related person
	//     unattributed
	//
	result =
		subjects.Organize(
			result,
			input,
			evidence,
		)

	//
	// Accounts have stronger semantics than
	// generic subject association.
	//
	// A subject-associated username/social may
	// still only be a mentioned account rather
	// than an account owned by that person.
	//
	result =
		accounts.Organize(
			result,
			input,
			evidence,
		)

	//
	// Phones have the strongest occurrence-level
	// ownership logic.
	//
	// Explicit:
	//
	//     Person1 +30...
	//
	// overrides the broader Evidence.Subject.
	//
	result =
		phones.Organize(
			result,
			input,
			evidence,
		)

	return result
}
