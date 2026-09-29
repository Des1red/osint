package social

import (
	"strings"

	"osint/internal/enrich/model"
	"osint/internal/logger"
)

const maxPostsPerOwner = 10

type builder struct {
	circle model.SocialCircle

	//
	// username -> already-known identity name.
	//
	identities map[string]string

	//
	// A username associated with more than one
	// known identity must not resolve
	// automatically.
	//
	ambiguous map[string]struct{}

	//
	// from + platform + username -> connection.
	//
	connections map[string]int

	//
	// Shared post cache for this enrichment run.
	//
	posts map[string]cachedPost
}

func Enrich(
	result model.EnrichmentResult,
	input model.Input,
) model.SocialCircle {
	root :=
		strings.TrimSpace(
			input.FullName,
		)

	if root == "" {

		return model.SocialCircle{}
	}

	state :=
		&builder{
			identities: make(
				map[string]string,
			),

			ambiguous: make(
				map[string]struct{},
			),

			connections: make(
				map[string]int,
			),

			posts: make(
				map[string]cachedPost,
			),
		}

	rootAccounts :=
		rootExtraAccounts(
			input,
		)

	//
	// Build the already-known identity index.
	//
	state.indexIdentity(
		root,
		result.Usernames,
		result.Socials,
		rootAccounts,
	)

	for _, person := range result.People {

		state.indexIdentity(
			person.Name,
			person.Usernames,
			person.Socials,
			nil,
		)
	}

	//
	// Analyze posts attributed to the root.
	//
	state.processOwner(
		root,
		result.Links,
		ownedAccounts(
			result.Usernames,
			result.Socials,
			rootAccounts,
		),
	)

	//
	// Analyze posts attributed to already-known
	// people.
	//
	for _, person := range result.People {

		state.processOwner(
			person.Name,
			person.Links,
			ownedAccounts(
				person.Usernames,
				person.Socials,
				nil,
			),
		)
	}

	logger.Debug(
		"Social Circle",
		"connections",
		state.circle.Connections,
		"evidence count",
		len(
			state.circle.Evidence,
		),
	)

	return state.circle
}
