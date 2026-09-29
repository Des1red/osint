package social

import (
	"strings"

	"osint/internal/enrich/model"
)

func (state *builder) indexIdentity(
	name string,
	usernames []model.UsernameReference,
	socials []model.SocialReference,
	extra []string,
) {
	name =
		strings.TrimSpace(
			name,
		)

	if name == "" {

		return
	}

	for _, username := range usernames {

		state.indexAccount(
			name,
			username.Username,
		)
	}

	for _, social := range socials {

		state.indexAccount(
			name,
			social.Username,
		)
	}

	for _, value := range extra {

		state.indexAccount(
			name,
			value,
		)
	}
}

func (state *builder) indexAccount(
	name string,
	username string,
) {
	key :=
		accountKey(
			username,
		)

	if key == "" {

		return
	}

	if _, ambiguous :=
		state.ambiguous[key]; ambiguous {

		return
	}

	existing,
		exists :=
		state.identities[key]

	if !exists {

		state.identities[key] =
			name

		return
	}

	if samePerson(
		existing,
		name,
	) {

		return
	}

	//
	// Same account appears to belong to more than
	// one already-known identity.
	//
	// Stop resolving that username automatically.
	//
	delete(
		state.identities,
		key,
	)

	state.ambiguous[key] =
		struct{}{}
}

func rootExtraAccounts(
	input model.Input,
) []string {
	var result []string

	if value :=
		strings.TrimSpace(
			input.RootUsername,
		); value != "" {

		result =
			append(
				result,
				value,
			)
	}

	for _, account := range input.Accounts {

		value :=
			strings.TrimSpace(
				account.Username,
			)

		if value == "" {

			continue
		}

		result =
			append(
				result,
				value,
			)
	}

	return result
}

func ownedAccounts(
	usernames []model.UsernameReference,
	socials []model.SocialReference,
	extra []string,
) map[string]struct{} {
	result :=
		make(
			map[string]struct{},
		)

	add :=
		func(
			value string,
		) {
			key :=
				accountKey(
					value,
				)

			if key == "" {

				return
			}

			result[key] =
				struct{}{}
		}

	for _, username := range usernames {

		add(
			username.Username,
		)
	}

	for _, social := range socials {

		add(
			social.Username,
		)
	}

	for _, value := range extra {

		add(
			value,
		)
	}

	return result
}

func accountKey(
	value string,
) string {
	return strings.ToLower(
		strings.Trim(
			strings.TrimSpace(
				value,
			),
			"@/",
		),
	)
}

func identityKey(
	value string,
) string {
	return strings.ToLower(
		strings.Join(
			strings.Fields(
				strings.TrimSpace(
					value,
				),
			),
			" ",
		),
	)
}

func samePerson(
	left string,
	right string,
) bool {
	left =
		identityKey(
			left,
		)

	if left == "" {

		return false
	}

	return left ==
		identityKey(
			right,
		)
}

func ignoredUsername(
	value string,
) bool {
	switch accountKey(
		value,
	) {

	case "",
		"instagram",
		"facebook",
		"twitter",
		"linkedin",
		"youtube",
		"tiktok",
		"official",
		"profile",
		"explore":

		return true
	}

	return false
}

func instagramProfileURL(
	username string,
) string {
	username =
		accountKey(
			username,
		)

	if username == "" {

		return ""
	}

	return "https://www.instagram.com/" +
		username +
		"/"
}
