package accounts

import (
	"strings"

	"osint/internal/enrich/model"
)

func evidenceForIdentity(
	fullName string,
	values []model.Evidence,
) []model.Evidence {
	fullName =
		strings.TrimSpace(
			fullName,
		)

	if fullName == "" {

		return nil
	}

	var result []model.Evidence

	for _, value := range values {

		if !evidenceHasPersonSubject(
			value,
			fullName,
		) {

			continue
		}

		result =
			append(
				result,
				value,
			)
	}

	return result
}

func evidenceHasPersonSubject(
	value model.Evidence,
	fullName string,
) bool {
	fullName =
		strings.TrimSpace(
			fullName,
		)

	if fullName == "" {

		return false
	}

	for _, subject := range value.Subjects {

		if subject.Kind !=
			model.EvidenceAnchorPerson {

			continue
		}

		if sameIdentityName(
			subject.Value,
			fullName,
		) {

			return true
		}
	}

	return false
}

func sameIdentityName(
	left string,
	right string,
) bool {
	leftFields :=
		strings.Fields(
			strings.ToLower(
				strings.TrimSpace(
					left,
				),
			),
		)

	rightFields :=
		strings.Fields(
			strings.ToLower(
				strings.TrimSpace(
					right,
				),
			),
		)

	if len(leftFields) == 0 ||
		len(rightFields) == 0 {

		return false
	}

	if strings.Join(
		leftFields,
		" ",
	) ==
		strings.Join(
			rightFields,
			" ",
		) {

		return true
	}

	if len(leftFields) < 2 ||
		len(rightFields) < 2 {

		return false
	}

	return leftFields[0] ==
		rightFields[len(rightFields)-1] &&
		leftFields[len(leftFields)-1] ==
			rightFields[0]
}

func evidenceReference(
	value model.Evidence,
	start int,
	end int,
) []model.EvidenceReference {
	id :=
		strings.TrimSpace(
			value.ID,
		)

	if id == "" {

		return nil
	}

	if start < 0 {

		start =
			0
	}

	if end < start {

		end =
			start
	}

	return []model.EvidenceReference{
		{
			EvidenceID: id,

			Start: start,

			End: end,
		},
	}
}
