package social

import (
	"strings"

	"osint/internal/enrich/model"
	"osint/internal/enrich/provenance"
)

func (state *builder) processOwner(
	owner string,
	links []model.ExternalLink,
	owned map[string]struct{},
) {
	owner =
		strings.TrimSpace(
			owner,
		)

	if owner == "" {

		return
	}

	processed :=
		0

	for _, link := range links {

		if processed >=
			maxPostsPerOwner {

			break
		}

		postURL,
			ok :=
			instagramPostFromLink(
				link,
			)

		if !ok {

			continue
		}

		key :=
			postKey(
				postURL,
			)

		if key == "" {

			continue
		}

		post,
			ok :=
			state.post(
				key,
				postURL,
			)

		processed++

		if !ok {

			continue
		}

		evidence :=
			model.Evidence{
				Anchor: provenance.Anchor(
					model.Input{
						FullName: owner,
					},
				),

				Source: "Social Post [Instagram]",

				Title: post.Title,

				Text: post.Text,

				URL: post.URL,
			}

		evidence.ID =
			provenance.ID(
				"social-post",
				evidence,
			)

		used :=
			false

		for _, mention := range post.Mentions {

			username :=
				accountKey(
					mention.Username,
				)

			if username == "" ||
				ignoredUsername(
					username,
				) {

				continue
			}

			//
			// Do not create:
			//
			//     Person -> their own account
			//
			if _, exists :=
				owned[username]; exists {

				continue
			}

			to :=
				"@" +
					username

			//
			// Resolve only against identities that
			// were already known before Social ran.
			//
			if identity,
				exists :=
				state.identities[username]; exists {

				if samePerson(
					identity,
					owner,
				) {

					continue
				}

				to =
					identity
			}

			reference :=
				model.EvidenceReference{
					EvidenceID: evidence.ID,

					Start: mention.Start,

					End: mention.End,
				}

			state.addConnection(
				owner,
				to,
				"instagram",
				username,
				post.URL,
				reference,
			)

			used =
				true
		}

		if used {

			state.circle.Evidence =
				provenance.Merge(
					state.circle.Evidence,
					[]model.Evidence{
						evidence,
					},
				)
		}
	}
}
