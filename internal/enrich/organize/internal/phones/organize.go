package phones

import (
	"strings"

	"osint/internal/enrich/model"
)

func Organize(
	result model.EnrichmentResult,
	input model.Input,
	evidence []model.Evidence,
) model.EnrichmentResult {
	//
	// Preserve username-only enrichment
	// semantics.
	//
	if strings.TrimSpace(
		input.FullName,
	) == "" {

		return result
	}

	if len(result.Phones) == 0 {

		return result
	}

	evidenceByID :=
		evidenceIndex(
			evidence,
		)

	people :=
		append(
			[]model.PersonReference(nil),
			result.People...,
		)

	identities :=
		identityCandidates(
			input,
			people,
		)

	var targetPhones []model.PhoneReference

	unattributed :=
		append(
			[]model.PhoneReference(nil),
			result.UnattributedPhones...,
		)

	for _, phone := range result.Phones {

		resolved,
			ok :=
			resolveOwner(
				phone,
				evidenceByID,
				identities,
			)

		if !ok {

			unattributed =
				append(
					unattributed,
					phone,
				)

			continue
		}

		switch resolved.Kind {

		case ownerTarget:

			targetPhones =
				append(
					targetPhones,
					phone,
				)

		case ownerPerson:

			if resolved.PersonIndex < 0 ||
				resolved.PersonIndex >=
					len(people) {

				unattributed =
					append(
						unattributed,
						phone,
					)

				continue
			}

			people[resolved.PersonIndex].Phones =
				append(
					people[resolved.PersonIndex].Phones,
					phone,
				)

		default:

			unattributed =
				append(
					unattributed,
					phone,
				)
		}
	}

	result.Phones =
		targetPhones

	result.UnattributedPhones =
		unattributed

	result.People =
		people

	return result
}
