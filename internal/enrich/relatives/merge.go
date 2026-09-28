package relatives

import "osint/internal/enrich/model"

func mergeRelativeFirstHop(
	result model.EnrichmentResult,
	relative model.EnrichmentResult,
) model.EnrichmentResult {
	//
	// Relative FirstHop now contributes the
	// complete raw fact set.
	//
	// Nothing in this function decides who owns
	// those facts.
	//
	// Every fact carries EvidenceReference values
	// which resolve to Evidence containing:
	//
	//     Anchor
	//     Subjects
	//
	// Semantic ownership is resolved later by
	// organize.Organize().
	//
	result.Usernames =
		append(
			result.Usernames,
			relative.Usernames...,
		)

	result.Emails =
		append(
			result.Emails,
			relative.Emails...,
		)

	result.Phones =
		append(
			result.Phones,
			relative.Phones...,
		)

	result.Socials =
		append(
			result.Socials,
			relative.Socials...,
		)

	result.Links =
		append(
			result.Links,
			relative.Links...,
		)

	result.Locations =
		append(
			result.Locations,
			relative.Locations...,
		)

	result.Organizations =
		append(
			result.Organizations,
			relative.Organizations...,
		)

	result.Employment =
		append(
			result.Employment,
			relative.Employment...,
		)

	result.Education =
		append(
			result.Education,
			relative.Education...,
		)

	result.RelatedAccounts =
		append(
			result.RelatedAccounts,
			relative.RelatedAccounts...,
		)

	//
	// Do not merge relative.People here.
	//
	// This pass deliberately remains one level
	// deep:
	//
	//     Person1
	//       ↓
	//     Person2
	//       ↓
	//     FirstHop
	//       ↓
	//     STOP
	//
	// People discovered during Person2 FirstHop can
	// become part of a later graph-expansion
	// stage without turning this pass into
	// recursive BFS.
	//

	return result
}
