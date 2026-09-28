package accounts

import (
	"osint/internal/enrich/model"
)

type ownedAccounts struct {
	Usernames map[string]struct{}

	Socials map[string]struct{}
}

func newOwnedAccounts(
	input model.Input,
) ownedAccounts {
	result :=
		ownedAccounts{
			Usernames: make(
				map[string]struct{},
			),

			Socials: make(
				map[string]struct{},
			),
		}

	rootUsername :=
		normalizeUsername(
			input.RootUsername,
		)

	if rootUsername != "" {

		result.Usernames[rootUsername] =
			struct{}{}
	}

	for _, account := range input.Accounts {

		result.add(
			account.Platform,
			account.Username,
		)
	}

	return result
}

func (
	value *ownedAccounts,
) add(
	platform string,
	username string,
) {
	username =
		normalizeUsername(
			username,
		)

	if username == "" {

		return
	}

	value.Usernames[username] =
		struct{}{}

	platform =
		normalizePlatform(
			platform,
		)

	if platform == "" {

		return
	}

	value.Socials[platform+
		":"+
		username] =
		struct{}{}
}

func (
	value ownedAccounts,
) hasUsername(
	username string,
) bool {
	username =
		normalizeUsername(
			username,
		)

	if username == "" {

		return false
	}

	_,
		exists :=
		value.Usernames[username]

	return exists
}

func (
	value ownedAccounts,
) hasSocial(
	platform string,
	username string,
) bool {
	platform =
		normalizePlatform(
			platform,
		)

	username =
		normalizeUsername(
			username,
		)

	if platform == "" ||
		username == "" {

		return false
	}

	_,
		exists :=
		value.Socials[platform+
			":"+
			username]

	return exists
}
