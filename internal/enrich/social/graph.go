package social

import (
	"strings"

	"osint/internal/enrich/model"
)

func (state *builder) addConnection(
	from string,
	to string,
	platform string,
	username string,
	postURL string,
	evidence model.EvidenceReference,
) {
	from =
		strings.TrimSpace(
			from,
		)

	to =
		strings.TrimSpace(
			to,
		)

	platform =
		strings.ToLower(
			strings.TrimSpace(
				platform,
			),
		)

	username =
		accountKey(
			username,
		)

	if from == "" ||
		to == "" ||
		username == "" {

		return
	}

	key :=
		identityKey(
			from,
		) +
			"\x00" +
			platform +
			"\x00" +
			username

	if index,
		exists :=
		state.connections[key]; exists {

		connection :=
			&state.circle.Connections[index]

		connection.Evidence =
			mergeEvidenceReferences(
				connection.Evidence,
				[]model.EvidenceReference{
					evidence,
				},
			)

		connection.Posts =
			appendUniqueString(
				connection.Posts,
				postURL,
			)

		return
	}

	state.circle.Connections =
		append(
			state.circle.Connections,
			model.SocialCircleConnection{
				From: from,

				To: to,

				Platform: platform,

				Username: username,

				ProfileURL: instagramProfileURL(
					username,
				),

				Posts: []string{
					postURL,
				},

				Evidence: []model.EvidenceReference{
					evidence,
				},
			},
		)

	state.connections[key] =
		len(
			state.circle.Connections,
		) -
			1
}
