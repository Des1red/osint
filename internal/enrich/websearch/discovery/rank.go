package discovery

import (
	"sort"
	"strings"
	"unicode"

	"osint/internal/enrich/model"
)

const maxResultsPerHost = 2

type rankedResult struct {
	Item SearchResult

	Score int
}

func prioritizeResults(
	input model.Input,
	values []SearchResult,
) []SearchResult {
	if len(values) == 0 {
		return nil
	}

	groups :=
		make(
			map[string][]rankedResult,
		)

	var queryOrder []string

	seenQueries :=
		make(
			map[string]struct{},
		)

	for _, value := range values {

		//
		// Do not let obviously unrelated results
		// into the ranking pool.
		//
		// Search engines are deliberately broad.
		// Ranking is where identity correlation
		// becomes strict.
		//
		if !resultEligible(
			input,
			value,
		) {
			continue
		}

		query :=
			strings.TrimSpace(
				value.Query,
			)

		queryKey :=
			strings.ToLower(
				query,
			)

		if queryKey == "" {

			queryKey =
				"_unknown"
		}

		if _, exists :=
			seenQueries[queryKey]; !exists {

			seenQueries[queryKey] =
				struct{}{}

			queryOrder =
				append(
					queryOrder,
					queryKey,
				)
		}

		groups[queryKey] =
			append(
				groups[queryKey],
				rankedResult{
					Item: value,

					Score: scoreResult(
						input,
						value,
					),
				},
			)
	}

	if len(groups) == 0 {
		return nil
	}

	//
	// Rank results inside each query first.
	//
	for key, group := range groups {

		sort.SliceStable(
			group,
			func(
				left int,
				right int,
			) bool {

				if group[left].Score ==
					group[right].Score {

					return strings.ToLower(
						group[left].Item.URL,
					) <
						strings.ToLower(
							group[right].Item.URL,
						)
				}

				return group[left].Score >
					group[right].Score
			},
		)

		groups[key] =
			group
	}

	positions :=
		make(
			map[string]int,
		)

	seenURLs :=
		make(
			map[string]struct{},
		)

	hostCounts :=
		make(
			map[string]int,
		)

	var result []SearchResult

	var deferred []rankedResult

	//
	// First pass:
	//
	// Round-robin across queries.
	//
	// Every search intent gets a chance to
	// contribute without allowing the first
	// query to consume the entire result budget.
	//
	for len(result) <
		maxResults {

		added :=
			false

		progressed :=
			false

		for _, queryKey := range queryOrder {

			group :=
				groups[queryKey]

			position :=
				positions[queryKey]

			for position <
				len(group) {

				candidate :=
					group[position]

				position++

				progressed =
					true

				canonical :=
					canonicalURL(
						candidate.Item.URL,
					)

				if canonical == "" {
					continue
				}

				if _, exists :=
					seenURLs[canonical]; exists {

					continue
				}

				host :=
					hostFromURL(
						canonical,
					)

				//
				// Soft host diversity limit.
				//
				// Do not discard extra results from
				// the same host completely. Keep them
				// for the second pass.
				//
				if host != "" &&
					hostCounts[host] >=
						maxResultsPerHost {

					deferred =
						append(
							deferred,
							candidate,
						)

					continue
				}

				seenURLs[canonical] =
					struct{}{}

				if host != "" {

					hostCounts[host]++
				}

				result =
					append(
						result,
						candidate.Item,
					)

				added =
					true

				break
			}

			positions[queryKey] =
				position

			if len(result) >=
				maxResults {

				break
			}
		}

		if !progressed ||
			!added {

			break
		}
	}

	//
	// Second pass:
	//
	// If hostname diversity left free slots,
	// use the strongest deferred candidates.
	//
	if len(result) <
		maxResults &&
		len(deferred) > 0 {

		sort.SliceStable(
			deferred,
			func(
				left int,
				right int,
			) bool {

				return deferred[left].Score >
					deferred[right].Score
			},
		)

		for _, candidate := range deferred {

			if len(result) >=
				maxResults {

				break
			}

			canonical :=
				canonicalURL(
					candidate.Item.URL,
				)

			if canonical == "" {
				continue
			}

			if _, exists :=
				seenURLs[canonical]; exists {

				continue
			}

			seenURLs[canonical] =
				struct{}{}

			result =
				append(
					result,
					candidate.Item,
				)
		}
	}

	return result
}

func resultEligible(
	input model.Input,
	item SearchResult,
) bool {
	title :=
		strings.ToLower(
			strings.TrimSpace(
				item.Title,
			),
		)

	snippet :=
		strings.ToLower(
			strings.TrimSpace(
				item.Snippet,
			),
		)

	resultURL :=
		strings.ToLower(
			strings.TrimSpace(
				item.URL,
			),
		)

	urlText :=
		urlSearchText(
			resultURL,
		)

	fullName :=
		strings.TrimSpace(
			input.FullName,
		)

	//
	// Full-name searches require actual identity
	// evidence.
	//
	// It is not enough for a result to contain
	// only "Person1" or only " ".
	//
	if fullName != "" {

		if identityPresent(
			title,
			fullName,
		) ||
			identityPresent(
				snippet,
				fullName,
			) ||
			identityPresent(
				urlText,
				fullName,
			) {

			return true
		}

		//
		// A result may not expose the person's
		// real name at all.
		//
		// Example:
		//
		// instagram.com/ Person1
		//
		// A known username/account is therefore
		// also strong enough identity evidence.
		//
		if knownUsernamePresent(
			input,
			title,
			snippet,
			resultURL,
		) {

			return true
		}

		//
		// Preserve a previously-known exact
		// profile URL even if its search card
		// contains little useful text.
		//
		if knownProfileURL(
			input,
			item.URL,
		) {

			return true
		}

		return false
	}

	//
	// Username-only searches should still have
	// an identity gate when username evidence
	// already exists.
	//
	if hasKnownUsername(
		input,
	) {

		if knownUsernamePresent(
			input,
			title,
			snippet,
			resultURL,
		) {

			return true
		}

		if knownProfileURL(
			input,
			item.URL,
		) {

			return true
		}

		return false
	}

	//
	// No identity information was supplied.
	//
	// There is nothing meaningful to gate on.
	//
	return true
}

func scoreResult(
	input model.Input,
	item SearchResult,
) int {
	title :=
		strings.ToLower(
			strings.TrimSpace(
				item.Title,
			),
		)

	snippet :=
		strings.ToLower(
			strings.TrimSpace(
				item.Snippet,
			),
		)

	resultURL :=
		strings.ToLower(
			strings.TrimSpace(
				item.URL,
			),
		)

	resultURLText :=
		urlSearchText(
			resultURL,
		)

	resultHost :=
		hostFromURL(
			item.URL,
		)

	score :=
		0

	fullName :=
		strings.ToLower(
			strings.TrimSpace(
				input.FullName,
			),
		)

	//
	// Real-name evidence is the strongest
	// ranking signal.
	//
	if fullName != "" {

		if identityPresent(
			title,
			fullName,
		) {

			score +=
				20
		}

		if identityPresent(
			snippet,
			fullName,
		) {

			score +=
				10
		}

		if identityPresent(
			resultURLText,
			fullName,
		) {

			score +=
				8
		}
	}

	rootUsername :=
		strings.ToLower(
			strings.TrimSpace(
				input.RootUsername,
			),
		)

	//
	// Root username is also strong identity
	// evidence.
	//
	if rootUsername != "" {

		if usernamePresent(
			title,
			rootUsername,
		) {

			score +=
				12
		}

		if usernamePresent(
			resultURL,
			rootUsername,
		) {

			score +=
				16
		}

		if usernamePresent(
			snippet,
			rootUsername,
		) {

			score +=
				6
		}
	}

	//
	// Known platform accounts provide strong
	// correlation signals.
	//
	for _, account := range input.Accounts {

		username :=
			strings.ToLower(
				strings.TrimSpace(
					account.Username,
				),
			)

		if username != "" {

			if usernamePresent(
				title,
				username,
			) {

				score +=
					10
			}

			if usernamePresent(
				resultURL,
				username,
			) {

				score +=
					14
			}

			if usernamePresent(
				snippet,
				username,
			) {

				score +=
					5
			}
		}

		//
		// Host agreement is useful, but it is
		// intentionally weak.
		//
		// linkedin.com alone does not imply that
		// a LinkedIn result belongs to our target.
		//
		host :=
			hostFromURL(
				account.ProfileURL,
			)

		if host != "" &&
			resultHost != "" &&
			host == resultHost {

			score +=
				2
		}

		accountURL :=
			canonicalURL(
				account.ProfileURL,
			)

		candidateURL :=
			canonicalURL(
				item.URL,
			)

		if accountURL != "" &&
			candidateURL != "" &&
			accountURL ==
				candidateURL {

			score +=
				20
		}
	}

	//
	// Other known usernames are useful secondary
	// correlation signals.
	//
	for _, value := range input.Usernames {

		username :=
			strings.ToLower(
				strings.TrimSpace(
					value.Username,
				),
			)

		if username == "" ||
			username ==
				rootUsername {

			continue
		}

		if usernamePresent(
			title,
			username,
		) {

			score +=
				6
		}

		if usernamePresent(
			resultURL,
			username,
		) {

			score +=
				9
		}

		if usernamePresent(
			snippet,
			username,
		) {

			score +=
				3
		}
	}

	//
	// Query/result consistency is deliberately
	// only a light bonus.
	//
	score +=
		queryContextScore(
			item.Query,
			title,
			snippet,
			resultURL,
		)

	return score
}

func identityPresent(
	value string,
	identity string,
) bool {
	valueTokens :=
		textTokens(
			value,
		)

	identityTokens :=
		textTokens(
			identity,
		)

	if len(valueTokens) == 0 ||
		len(identityTokens) == 0 {

		return false
	}

	//
	// For names containing middle names or
	// initials, the first and last components
	// are the stable identity anchors.
	//
	// Examples:
	//
	// Person1 Konstantinos
	// Person1 G.
	//
	required :=
		identityTokens

	if len(identityTokens) > 2 {

		required =
			[]string{
				identityTokens[0],
				identityTokens[len(identityTokens)-1],
			}
	}

	available :=
		make(
			map[string]struct{},
			len(valueTokens),
		)

	for _, token := range valueTokens {

		available[token] =
			struct{}{}
	}

	for _, token := range required {

		if _, exists :=
			available[token]; !exists {

			return false
		}
	}

	return true
}

func textTokens(
	value string,
) []string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	if value == "" {
		return nil
	}

	return strings.FieldsFunc(
		value,
		func(
			character rune,
		) bool {

			return !unicode.IsLetter(
				character,
			) &&
				!unicode.IsDigit(
					character,
				)
		},
	)
}

func knownUsernamePresent(
	input model.Input,
	title string,
	snippet string,
	resultURL string,
) bool {
	check :=
		func(
			username string,
		) bool {

			username =
				strings.ToLower(
					strings.TrimSpace(
						username,
					),
				)

			if len(username) < 3 {
				return false
			}

			return usernamePresent(
				title,
				username,
			) ||
				usernamePresent(
					snippet,
					username,
				) ||
				usernamePresent(
					resultURL,
					username,
				)
		}

	if check(
		input.RootUsername,
	) {
		return true
	}

	for _, account := range input.Accounts {

		if check(
			account.Username,
		) {
			return true
		}
	}

	for _, value := range input.Usernames {

		if check(
			value.Username,
		) {
			return true
		}
	}

	return false
}

func hasKnownUsername(
	input model.Input,
) bool {
	if len(
		strings.TrimSpace(
			input.RootUsername,
		),
	) >= 3 {

		return true
	}

	for _, account := range input.Accounts {

		if len(
			strings.TrimSpace(
				account.Username,
			),
		) >= 3 {

			return true
		}
	}

	for _, value := range input.Usernames {

		if len(
			strings.TrimSpace(
				value.Username,
			),
		) >= 3 {

			return true
		}
	}

	return false
}

func usernamePresent(
	value string,
	username string,
) bool {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	username =
		strings.ToLower(
			strings.TrimSpace(
				username,
			),
		)

	if value == "" ||
		username == "" {

		return false
	}

	return strings.Contains(
		value,
		username,
	)
}

func knownProfileURL(
	input model.Input,
	value string,
) bool {
	candidate :=
		canonicalURL(
			value,
		)

	if candidate == "" {
		return false
	}

	for _, account := range input.Accounts {

		known :=
			canonicalURL(
				account.ProfileURL,
			)

		if known == "" {
			continue
		}

		if candidate ==
			known {

			return true
		}
	}

	return false
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

func queryContextScore(
	query string,
	title string,
	snippet string,
	resultURL string,
) int {
	query =
		strings.ToLower(
			strings.TrimSpace(
				query,
			),
		)

	if query == "" {
		return 0
	}

	replacer :=
		strings.NewReplacer(
			`"`,
			" ",

			`'`,
			" ",

			"(",
			" ",

			")",
			" ",
		)

	terms :=
		strings.Fields(
			replacer.Replace(
				query,
			),
		)

	score :=
		0

	for _, term := range terms {

		term =
			strings.TrimSpace(
				term,
			)

		if term == "" ||
			len(term) < 3 {

			continue
		}

		if strings.HasPrefix(
			term,
			"site:",
		) {

			continue
		}

		//
		// These describe what we are looking
		// for, not who the target is.
		//
		// They must not make an unrelated result
		// appear more identity-relevant.
		//
		switch term {

		case "company",
			"owner",
			"founder",
			"linkedin",
			"email",
			"phone",
			"contact",
			"location",
			"profile":

			continue
		}

		switch {

		case tokenPresent(
			title,
			term,
		):

			score +=
				2

		case tokenPresent(
			urlSearchText(
				resultURL,
			),
			term,
		):

			score++

		case tokenPresent(
			snippet,
			term,
		):

			score++
		}
	}

	return score
}

func tokenPresent(
	value string,
	term string,
) bool {
	term =
		strings.ToLower(
			strings.TrimSpace(
				term,
			),
		)

	if term == "" {
		return false
	}

	for _, token := range textTokens(
		value,
	) {

		if token ==
			term {

			return true
		}
	}

	return false
}
