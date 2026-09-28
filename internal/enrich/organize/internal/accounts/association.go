package accounts

import (
	"net/url"
	"strings"

	"osint/internal/enrich/model"
)

var accountMonthSuffixes = []string{
	"september",
	"february",
	"november",
	"december",
	"january",
	"october",
	"august",
	"march",
	"april",
	"june",
	"july",
	"sept",
	"jan",
	"feb",
	"mar",
	"apr",
	"jun",
	"jul",
	"aug",
	"sep",
	"oct",
	"nov",
	"dec",
	"may",
}

func usernameAssociationAllowed(
	value model.UsernameReference,
	evidence []model.Evidence,
) bool {
	index :=
		accountEvidenceIndex(
			evidence,
		)

	for _, reference := range value.Evidence {

		item,
			exists :=
			index[strings.TrimSpace(
				reference.EvidenceID,
			)]

		if !exists {

			continue
		}

		if accountAssociationEvidence(
			item,
		) {

			return true
		}
	}

	return false
}

func socialAssociationAllowed(
	value model.SocialReference,
	evidence []model.Evidence,
) bool {
	index :=
		accountEvidenceIndex(
			evidence,
		)

	for _, reference := range value.Evidence {

		item,
			exists :=
			index[strings.TrimSpace(
				reference.EvidenceID,
			)]

		if !exists {

			continue
		}

		if accountAssociationEvidence(
			item,
		) {

			return true
		}
	}

	return false
}

func accountEvidenceIndex(
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

func accountAssociationEvidence(
	value model.Evidence,
) bool {
	evidenceURL :=
		strings.TrimSpace(
			value.URL,
		)

	if evidenceURL == "" {

		return false
	}

	//
	// A list of somebody else's followers,
	// following accounts, friends, connections,
	// etc. is not meaningful account-association
	// evidence.
	//
	if accountListURL(
		evidenceURL,
	) {

		return false
	}

	//
	// Raw @mentions are only promoted into
	// RelatedAccounts when they occur on a
	// recognised social platform.
	//
	// A mirror / directory / arbitrary website
	// containing thousands of handles should not
	// become an account graph.
	//
	return platformFromURL(
		evidenceURL,
	) != ""
}

func accountListURL(
	value string,
) bool {
	parsed,
		err :=
		url.Parse(
			strings.TrimSpace(
				value,
			),
		)

	if err != nil {

		return false
	}

	segments :=
		strings.Split(
			strings.ToLower(
				strings.Trim(
					parsed.Path,
					"/",
				),
			),
			"/",
		)

	for _, segment := range segments {

		switch segment {

		case "following",
			"followers",
			"friends",
			"connections",
			"subscribers",
			"subscriptions":

			return true
		}
	}

	query :=
		strings.ToLower(
			parsed.RawQuery,
		)

	for _, marker := range []string{
		"following",
		"followers",
		"friends",
		"connections",
		"subscribers",
		"subscriptions",
	} {

		if strings.Contains(
			query,
			marker,
		) {

			return true
		}
	}

	return false
}

func cleanAccountMention(
	value string,
	trailing string,
) string {
	username :=
		normalizeUsername(
			value,
		)

	if username == "" {

		return ""
	}

	if strings.HasSuffix(
		username,
		".",
	) &&
		accountTrailingBoundary(
			trailing,
		) {

		username =
			strings.TrimSuffix(
				username,
				".",
			)
	}

	lower :=
		strings.ToLower(
			username,
		)

	monthSuffix :=
		accountAttachedMonthSuffix(
			lower,
		)

	if monthSuffix != "" &&
		accountTrailingStartsDigit(
			trailing,
		) {

		username =
			username[:len(username)-
				len(monthSuffix)]

		lower =
			strings.ToLower(
				username,
			)
	}

	if strings.HasSuffix(
		lower,
		"joined",
	) &&
		(accountTrailingStartsMonth(
			trailing,
		) ||
			accountTrailingStartsDigit(
				trailing,
			)) {

		username =
			username[:len(username)-
				len("joined")]
	}

	return normalizeUsername(
		username,
	)
}

func accountAttachedMonthSuffix(
	value string,
) string {
	for _, suffix := range accountMonthSuffixes {

		if len(value) <=
			len(suffix) {

			continue
		}

		if strings.HasSuffix(
			value,
			suffix,
		) {

			return suffix
		}
	}

	return ""
}

func accountTrailingStartsDigit(
	value string,
) bool {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {

		return false
	}

	return value[0] >= '0' &&
		value[0] <= '9'
}

func accountTrailingStartsMonth(
	value string,
) bool {
	fields :=
		strings.Fields(
			strings.TrimSpace(
				value,
			),
		)

	if len(fields) == 0 {

		return false
	}

	word :=
		strings.ToLower(
			strings.Trim(
				fields[0],
				`"'“”‘’.,;:()[]{}-`,
			),
		)

	switch word {

	case "january",
		"february",
		"march",
		"april",
		"may",
		"june",
		"july",
		"august",
		"september",
		"october",
		"november",
		"december",
		"jan",
		"feb",
		"mar",
		"apr",
		"jun",
		"jul",
		"aug",
		"sep",
		"sept",
		"oct",
		"nov",
		"dec":

		return true
	}

	return false
}

func accountTrailingBoundary(
	value string,
) bool {
	if value == "" {

		return true
	}

	switch value[0] {

	case ' ',
		'\t',
		'\n',
		'\r',
		',',
		';',
		':',
		'!',
		'?':

		return true
	}

	return false
}
