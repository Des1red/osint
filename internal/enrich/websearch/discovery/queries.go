package discovery

import (
	"strings"

	"osint/internal/enrich/model"
)

const maxQueries = 12

const maxAccountQueries = 2

func searchQueries(
	input model.Input,
) []string {
	var queries []string

	seen :=
		make(
			map[string]struct{},
		)

	add :=
		func(
			value string,
		) bool {
			if len(queries) >=
				maxQueries {

				return false
			}

			value =
				strings.TrimSpace(
					value,
				)

			if value == "" {
				return false
			}

			key :=
				strings.ToLower(
					value,
				)

			if _, exists :=
				seen[key]; exists {

				return false
			}

			seen[key] =
				struct{}{}

			queries =
				append(
					queries,
					value,
				)

			return true
		}

	fullName :=
		strings.TrimSpace(
			input.FullName,
		)

	rootUsername :=
		strings.TrimSpace(
			input.RootUsername,
		)

	//
	// Primary full-name identity searches.
	//
	// Use both:
	//
	//   "Person1"
	//   Person1
	//
	// The quoted search gives us precision.
	//
	// The unquoted search gives us broader
	// discovery and lets the search engine
	// find pages where the identity appears
	// naturally in titles, descriptions or
	// page content.
	//
	if fullName != "" {

		add(
			`"` +
				fullName +
				`"`,
		)

		add(
			fullName,
		)
	}

	//
	// Primary username search.
	//
	if rootUsername != "" {

		add(
			`"` +
				rootUsername +
				`"`,
		)
	}

	//
	// When both identifiers are available,
	// correlate them directly.
	//
	if fullName != "" &&
		rootUsername != "" {

		add(
			fullName +
				" " +
				rootUsername,
		)
	}

	//
	// Full-name contextual discovery.
	//
	// These searches intentionally do NOT
	// quote the full name.
	//
	// Discovery should have high recall.
	// Ranking, extraction and FirstHop perform
	// the later relevance checks.
	//
	if fullName != "" {

		//
		// Organization / employment evidence.
		//
		add(
			fullName +
				" company",
		)

		add(
			fullName +
				" owner",
		)

		add(
			fullName +
				" founder",
		)

		//
		// Professional profiles.
		//
		add(
			fullName +
				" linkedin",
		)

		//
		// Public contact information.
		//
		add(
			fullName +
				" email",
		)

		add(
			fullName +
				" phone",
		)

		add(
			fullName +
				" contact",
		)

		//
		// Broad geographic / identity context.
		//
		add(
			fullName +
				" location",
		)

		add(
			fullName +
				" profile",
		)
	}

	//
	// Full name + known account username.
	//
	// These searches are useful when there is
	// still room in the query budget.
	//
	accountQueries :=
		0

	if fullName != "" {

		for _, account := range input.Accounts {

			if accountQueries >=
				maxAccountQueries {

				break
			}

			username :=
				strings.TrimSpace(
					account.Username,
				)

			if username == "" {
				continue
			}

			if add(
				fullName +
					" " +
					username,
			) {

				accountQueries++
			}
		}
	}

	//
	// Username + known platform domain.
	//
	hostQueries :=
		0

	seenHosts :=
		make(
			map[string]struct{},
		)

	if rootUsername != "" {

		for _, account := range input.Accounts {

			if hostQueries >=
				maxAccountQueries {

				break
			}

			host :=
				hostFromURL(
					account.ProfileURL,
				)

			if host == "" {
				continue
			}

			if _, exists :=
				seenHosts[host]; exists {

				continue
			}

			seenHosts[host] =
				struct{}{}

			if add(
				`"` +
					rootUsername +
					`" site:` +
					host,
			) {

				hostQueries++
			}
		}
	}

	//
	// Correlate the root username with another
	// username discovered from known evidence.
	//
	if rootUsername != "" {

		for _, value := range input.Usernames {

			username :=
				strings.TrimSpace(
					value.Username,
				)

			if username == "" ||
				strings.EqualFold(
					username,
					rootUsername,
				) {

				continue
			}

			if add(
				rootUsername +
					" " +
					username,
			) {

				break
			}
		}
	}

	//
	// Username-specific contextual searches.
	//
	if rootUsername != "" {

		add(
			rootUsername +
				" email",
		)

		add(
			rootUsername +
				" contact",
		)

		add(
			rootUsername +
				" profile",
		)
	}

	return queries
}

func uniqueQueries(
	queries []string,
) []string {
	seen :=
		make(
			map[string]struct{},
		)

	var result []string

	for _, query := range queries {

		query =
			strings.TrimSpace(
				query,
			)

		if query == "" {
			continue
		}

		key :=
			strings.ToLower(
				query,
			)

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		result =
			append(
				result,
				query,
			)
	}

	return result
}
