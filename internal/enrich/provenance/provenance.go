package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"osint/internal/enrich/model"
)

func Anchor(
	input model.Input,
) model.EvidenceAnchor {
	fullName :=
		strings.TrimSpace(
			input.FullName,
		)

	if fullName != "" {

		return model.EvidenceAnchor{
			Kind: model.EvidenceAnchorPerson,

			Value: fullName,
		}
	}

	rootUsername :=
		strings.TrimSpace(
			input.RootUsername,
		)

	if rootUsername != "" {

		return model.EvidenceAnchor{
			Kind: model.EvidenceAnchorAccount,

			Value: rootUsername,
		}
	}

	return model.EvidenceAnchor{
		Kind: model.EvidenceAnchorUnknown,
	}
}

// Subjects returns the currently known
// investigation graph.
//
// For FullName enrichment this currently means:
//
// root person
// +
// people already discovered around that root.
//
// This does not recursively expand the graph.
func Subjects(
	input model.Input,
	people []model.PersonReference,
) []model.EvidenceSubject {
	var result []model.EvidenceSubject

	seen :=
		make(
			map[string]struct{},
		)

	add :=
		func(
			kind model.EvidenceAnchorKind,
			value string,
		) {
			value =
				strings.TrimSpace(
					value,
				)

			if value == "" {

				return
			}

			key :=
				string(kind) +
					":" +
					strings.ToLower(
						strings.Join(
							strings.Fields(
								value,
							),
							" ",
						),
					)

			if _, exists :=
				seen[key]; exists {

				return
			}

			seen[key] =
				struct{}{}

			result =
				append(
					result,
					model.EvidenceSubject{
						Kind: kind,

						Value: value,
					},
				)
		}

	add(
		model.EvidenceAnchorPerson,
		input.FullName,
	)

	for _, person := range people {

		add(
			model.EvidenceAnchorPerson,
			person.Name,
		)
	}

	return result
}

func ID(
	kind string,
	value model.Evidence,
) string {
	kind =
		strings.ToLower(
			strings.TrimSpace(
				kind,
			),
		)

	if kind == "" {

		kind =
			"evidence"
	}

	parts :=
		[]string{
			string(
				value.Anchor.Kind,
			),

			normalizeIDPart(
				value.Anchor.Value,
			),

			subjectIDPart(
				value.Subjects,
			),

			normalizeIDPart(
				value.Query,
			),

			normalizeIDPart(
				value.Engine,
			),

			normalizeIDPart(
				value.Title,
			),

			normalizeIDPart(
				value.Text,
			),

			normalizeIDPart(
				value.URL,
			),
		}

	sum :=
		sha256.Sum256(
			[]byte(
				strings.Join(
					parts,
					"\x1f",
				),
			),
		)

	return kind +
		":" +
		hex.EncodeToString(
			sum[:16],
		)
}

func Merge(
	groups ...[]model.Evidence,
) []model.Evidence {
	seen :=
		make(
			map[string]struct{},
		)

	var result []model.Evidence

	for _, group := range groups {

		for _, value := range group {

			id :=
				strings.TrimSpace(
					value.ID,
				)

			if id != "" {

				if _, exists :=
					seen[id]; exists {

					continue
				}

				seen[id] =
					struct{}{}
			}

			result =
				append(
					result,
					value,
				)
		}
	}

	return result
}

func subjectIDPart(
	values []model.EvidenceSubject,
) string {
	if len(values) == 0 {

		return ""
	}

	parts :=
		make(
			[]string,
			0,
			len(values),
		)

	for _, value := range values {

		subjectValue :=
			normalizeIDPart(
				value.Value,
			)

		if subjectValue == "" {

			continue
		}

		parts =
			append(
				parts,
				string(value.Kind)+
					":"+
					strings.ToLower(
						subjectValue,
					),
			)
	}

	sort.Strings(
		parts,
	)

	return strings.Join(
		parts,
		"|",
	)
}

func normalizeIDPart(
	value string,
) string {
	return strings.Join(
		strings.Fields(
			strings.TrimSpace(
				value,
			),
		),
		" ",
	)
}
