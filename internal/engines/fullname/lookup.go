package fullname

import (
	"osint/internal/platforms"
)

func lookupCandidates(
	fullName string,
	searchCandidates []string,
) []Candidate {
	discovered :=
		platforms.NameLookup(
			fullName,
			searchCandidates,
		)

	if len(discovered) == 0 {
		return nil
	}

	candidates :=
		make(
			[]Candidate,
			0,
			len(discovered),
		)

	seen :=
		make(
			map[string]struct{},
		)

	for _, candidate := range discovered {

		key :=
			candidateKey(
				candidate.Platform,
				candidate.Username,
				candidate.ProfileURL,
			)

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		candidates =
			append(
				candidates,
				Candidate{
					Platform: candidate.Platform,

					Username: candidate.Username,

					Name: candidate.Name,

					ProfileURL: candidate.ProfileURL,
				},
			)
	}

	return candidates
}
