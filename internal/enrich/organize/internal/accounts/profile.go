package accounts

import (
	"strings"

	"osint/internal/enrich/extract"
	"osint/internal/enrich/model"
)

func ownedProfileFromEvidence(
	input model.Input,
	item model.Evidence,
) (
	string,
	string,
	bool,
) {
	platform,
		username,
		ok :=
		extract.SocialFromURL(
			item.URL,
		)

	if !ok {

		return "",
			"",
			false
	}

	//
	// Accounts already known before this
	// enrichment stage are strong ownership
	// evidence.
	//
	for _, account := range input.Accounts {

		if !strings.EqualFold(
			strings.TrimSpace(
				account.Platform,
			),
			platform,
		) {

			continue
		}

		if strings.EqualFold(
			normalizeUsername(
				account.Username,
			),
			normalizeUsername(
				username,
			),
		) {

			return platform,
				username,
				true
		}
	}

	fullName :=
		strings.TrimSpace(
			input.FullName,
		)

	if fullName == "" {

		return "",
			"",
			false
	}

	//
	// Require the person's identity in the
	// direct profile title.
	//
	// This distinguishes:
	//
	// Person1 (@ Person1)
	//
	// from:
	//
	// company Way (@straightrazorway)
	//
	// when the latter merely mentions Person1
	// elsewhere in its evidence.
	//
	if identityPresent(
		item.Title,
		fullName,
	) {

		return platform,
			username,
			true
	}

	//
	// Some profile URLs themselves expose the
	// person's real identity.
	//
	if identityPresent(
		urlSearchText(
			item.URL,
		),
		fullName,
	) {

		return platform,
			username,
			true
	}

	return "",
		"",
		false
}
