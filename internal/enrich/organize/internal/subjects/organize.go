package subjects

import (
	"strings"

	"osint/internal/enrich/model"
)

func Organize(
	result model.EnrichmentResult,
	input model.Input,
	evidence []model.Evidence,
) model.EnrichmentResult {
	//
	// Preserve username-only enrichment
	// semantics.
	//
	if strings.TrimSpace(
		input.FullName,
	) == "" {

		return result
	}

	evidenceByID :=
		EvidenceIndex(
			evidence,
		)

	people :=
		append(
			[]model.PersonReference(nil),
			result.People...,
		)

	result.Usernames,
		result.UnattributedUsernames =
		route(
			result.Usernames,
			result.UnattributedUsernames,
			evidenceByID,
			input,
			people,
			func(
				value model.UsernameReference,
			) []model.EvidenceReference {

				return value.Evidence
			},
			func(
				person *model.PersonReference,
				value model.UsernameReference,
			) {

				person.Usernames =
					append(
						person.Usernames,
						value,
					)
			},
		)

	result.Emails,
		result.UnattributedEmails =
		route(
			result.Emails,
			result.UnattributedEmails,
			evidenceByID,
			input,
			people,
			func(
				value model.EmailReference,
			) []model.EvidenceReference {

				return value.Evidence
			},
			func(
				person *model.PersonReference,
				value model.EmailReference,
			) {

				person.Emails =
					append(
						person.Emails,
						value,
					)
			},
		)

	result.Socials,
		result.UnattributedSocials =
		route(
			result.Socials,
			result.UnattributedSocials,
			evidenceByID,
			input,
			people,
			func(
				value model.SocialReference,
			) []model.EvidenceReference {

				return value.Evidence
			},
			func(
				person *model.PersonReference,
				value model.SocialReference,
			) {

				person.Socials =
					append(
						person.Socials,
						value,
					)
			},
		)

	result.Links,
		result.UnattributedLinks =
		route(
			result.Links,
			result.UnattributedLinks,
			evidenceByID,
			input,
			people,
			func(
				value model.ExternalLink,
			) []model.EvidenceReference {

				return value.Evidence
			},
			func(
				person *model.PersonReference,
				value model.ExternalLink,
			) {

				person.Links =
					append(
						person.Links,
						value,
					)
			},
		)

	result.Locations,
		result.UnattributedLocations =
		route(
			result.Locations,
			result.UnattributedLocations,
			evidenceByID,
			input,
			people,
			func(
				value model.LocationReference,
			) []model.EvidenceReference {

				return value.Evidence
			},
			func(
				person *model.PersonReference,
				value model.LocationReference,
			) {

				person.Locations =
					append(
						person.Locations,
						value,
					)
			},
		)

	result.Organizations,
		result.UnattributedOrganizations =
		route(
			result.Organizations,
			result.UnattributedOrganizations,
			evidenceByID,
			input,
			people,
			func(
				value model.OrganizationReference,
			) []model.EvidenceReference {

				return value.Evidence
			},
			func(
				person *model.PersonReference,
				value model.OrganizationReference,
			) {

				person.Organizations =
					append(
						person.Organizations,
						value,
					)
			},
		)

	result.Employment,
		result.UnattributedEmployment =
		route(
			result.Employment,
			result.UnattributedEmployment,
			evidenceByID,
			input,
			people,
			func(
				value model.EmploymentReference,
			) []model.EvidenceReference {

				return value.Evidence
			},
			func(
				person *model.PersonReference,
				value model.EmploymentReference,
			) {

				person.Employment =
					append(
						person.Employment,
						value,
					)
			},
		)

	result.Education,
		result.UnattributedEducation =
		route(
			result.Education,
			result.UnattributedEducation,
			evidenceByID,
			input,
			people,
			func(
				value model.EducationReference,
			) []model.EvidenceReference {

				return value.Evidence
			},
			func(
				person *model.PersonReference,
				value model.EducationReference,
			) {

				person.Education =
					append(
						person.Education,
						value,
					)
			},
		)

	result.RelatedAccounts,
		result.UnattributedRelatedAccounts =
		route(
			result.RelatedAccounts,
			result.UnattributedRelatedAccounts,
			evidenceByID,
			input,
			people,
			func(
				value model.RelatedAccountReference,
			) []model.EvidenceReference {

				return value.Evidence
			},
			func(
				person *model.PersonReference,
				value model.RelatedAccountReference,
			) {

				person.RelatedAccounts =
					append(
						person.RelatedAccounts,
						value,
					)
			},
		)

	result.People =
		people

	return result
}

func route[
	T any,
](
	values []T,
	unattributed []T,
	evidence map[string]model.Evidence,
	input model.Input,
	people []model.PersonReference,
	references func(
		T,
	) []model.EvidenceReference,
	addPerson func(
		*model.PersonReference,
		T,
	),
) (
	[]T,
	[]T,
) {
	var target []T

	for _, value := range values {

		refs :=
			references(
				value,
			)

		//
		// A fact with no provenance at all was
		// already semantic before this WebSearch
		// organization stage.
		//
		// Keep it at the root instead of silently
		// turning pre-existing data into
		// unattributed data.
		//
		if len(refs) == 0 {

			target =
				append(
					target,
					value,
				)

			continue
		}

		owners :=
			ResolveAll(
				refs,
				evidence,
				input,
				people,
			)

		if len(owners) == 0 {

			unattributed =
				append(
					unattributed,
					value,
				)

			continue
		}

		assigned :=
			false

		for _, resolved := range owners {

			switch resolved.Kind {

			case OwnerTarget:

				target =
					append(
						target,
						value,
					)

				assigned =
					true

			case OwnerPerson:

				if resolved.PersonIndex < 0 ||
					resolved.PersonIndex >=
						len(people) {

					continue
				}

				addPerson(
					&people[resolved.PersonIndex],
					value,
				)

				assigned =
					true
			}
		}

		if !assigned {

			unattributed =
				append(
					unattributed,
					value,
				)
		}
	}

	return target,
		unattributed
}
