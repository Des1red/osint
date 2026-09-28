package phones

import (
	"strings"

	"osint/internal/enrich/model"
)

func evidenceIndex(
	values []model.Evidence,
) map[string]model.Evidence {
	result :=
		make(
			map[string]model.Evidence,
			len(values),
		)

	for _, value := range values {

		id :=
			strings.TrimSpace(
				value.ID,
			)

		if id == "" {

			continue
		}

		result[id] =
			value
	}

	return result
}

func resolveOwner(
	phone model.PhoneReference,
	evidence map[string]model.Evidence,
	identities []identityCandidate,
) (
	owner,
	bool,
) {
	var resolved owner

	found :=
		false

	//
	// First try explicit occurrence-level
	// attribution.
	//
	// Example:
	//
	//     Person1 +30...
	//     Person2 +44...
	//
	// This is stronger than Evidence.Subject.
	//
	for _, reference := range phone.Evidence {

		evidenceID :=
			strings.TrimSpace(
				reference.EvidenceID,
			)

		if evidenceID == "" {

			continue
		}

		item,
			exists :=
			evidence[evidenceID]

		if !exists {

			continue
		}

		current,
			ok :=
			ownerFromOccurrence(
				item,
				reference,
				identities,
			)

		if !ok {

			continue
		}

		if !found {

			resolved =
				current

			found =
				true

			continue
		}

		//
		// Conflicting explicit attribution is
		// genuinely ambiguous.
		//
		// Do not fall back to Subject in this
		// situation.
		//
		if !sameOwner(
			resolved,
			current,
		) {

			return owner{},
				false
		}
	}

	if found {

		return resolved,
			true
	}

	//
	// No explicit name-near-phone attribution
	// was available.
	//
	// Now Evidence.Subject may provide the
	// broader semantic association.
	//
	return ownerFromSubjects(
		phone.Evidence,
		evidence,
		identities,
	)
}

func ownerFromOccurrence(
	evidence model.Evidence,
	reference model.EvidenceReference,
	identities []identityCandidate,
) (
	owner,
	bool,
) {
	text :=
		evidence.Text

	if text == "" {

		return owner{},
			false
	}

	start :=
		reference.Start

	end :=
		reference.End

	//
	// Whole-evidence references are not strong
	// enough for explicit occurrence ownership.
	//
	if start ==
		end {

		return owner{},
			false
	}

	if start < 0 ||
		end < start ||
		end > len(text) {

		return owner{},
			false
	}

	bestScore :=
		0

	var bestOwner owner

	tied :=
		false

	for _, identity := range identities {

		score :=
			identityScore(
				text,
				start,
				identity,
			)

		if score == 0 {

			continue
		}

		if score >
			bestScore {

			bestScore =
				score

			bestOwner =
				identity.Owner

			tied =
				false

			continue
		}

		if score ==
			bestScore &&
			!sameOwner(
				bestOwner,
				identity.Owner,
			) {

			tied =
				true
		}
	}

	if bestScore == 0 ||
		tied {

		return owner{},
			false
	}

	return bestOwner,
		true
}

func ownerFromSubjects(
	references []model.EvidenceReference,
	evidence map[string]model.Evidence,
	identities []identityCandidate,
) (
	owner,
	bool,
) {
	var resolved owner

	found :=
		false

	for _, reference := range references {

		evidenceID :=
			strings.TrimSpace(
				reference.EvidenceID,
			)

		if evidenceID == "" {

			continue
		}

		item,
			exists :=
			evidence[evidenceID]

		if !exists {

			continue
		}

		for _, subject := range item.Subjects {

			if subject.Kind !=
				model.EvidenceAnchorPerson {

				continue
			}

			current,
				ok :=
				ownerForSubject(
					subject.Value,
					identities,
				)

			if !ok {

				continue
			}

			if !found {

				resolved =
					current

				found =
					true

				continue
			}

			if !sameOwner(
				resolved,
				current,
			) {

				return owner{},
					false
			}
		}
	}

	return resolved,
		found
}

func ownerForSubject(
	value string,
	identities []identityCandidate,
) (
	owner,
	bool,
) {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {

		return owner{},
			false
	}

	var resolved owner

	found :=
		false

	for _, identity := range identities {

		if !sameIdentityName(
			value,
			identity.Name,
		) {

			continue
		}

		if !found {

			resolved =
				identity.Owner

			found =
				true

			continue
		}

		if !sameOwner(
			resolved,
			identity.Owner,
		) {

			return owner{},
				false
		}
	}

	return resolved,
		found
}
