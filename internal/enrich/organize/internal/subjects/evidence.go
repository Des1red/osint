package subjects

import (
	"fmt"
	"strings"

	"osint/internal/enrich/model"
)

func EvidenceIndex(
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

func ResolveAll(
	references []model.EvidenceReference,
	evidence map[string]model.Evidence,
	input model.Input,
	people []model.PersonReference,
) []Owner {
	seen :=
		make(
			map[string]struct{},
		)

	var result []Owner

	for _, reference := range references {

		evidenceID :=
			strings.TrimSpace(
				reference.EvidenceID,
			)

		if evidenceID == "" {

			continue
		}

		item,
			exists :=
			evidence[evidenceID]

		if !exists {

			continue
		}

		for _, subject := range item.Subjects {

			current,
				ok :=
				ownerFromSubject(
					subject,
					input,
					people,
				)

			if !ok {

				continue
			}

			key :=
				ownerKey(
					current,
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
					current,
				)
		}
	}

	return result
}

func ownerKey(
	value Owner,
) string {
	switch value.Kind {

	case OwnerTarget:

		return "target"

	case OwnerPerson:

		return fmt.Sprintf(
			"person:%d",
			value.PersonIndex,
		)

	default:

		return "unknown"
	}
}

func ownerFromSubject(
	subject model.EvidenceSubject,
	input model.Input,
	people []model.PersonReference,
) (
	Owner,
	bool,
) {
	if subject.Kind !=
		model.EvidenceAnchorPerson {

		return Owner{},
			false
	}

	name :=
		strings.TrimSpace(
			subject.Value,
		)

	if name == "" {

		return Owner{},
			false
	}

	target :=
		strings.TrimSpace(
			input.FullName,
		)

	if target != "" &&
		sameIdentityName(
			target,
			name,
		) {

		return Owner{
				Kind: OwnerTarget,
			},
			true
	}

	foundIndex :=
		-1

	for index, person := range people {

		personName :=
			strings.TrimSpace(
				person.Name,
			)

		if personName == "" {

			continue
		}

		if !sameIdentityName(
			personName,
			name,
		) {

			continue
		}

		if foundIndex >= 0 &&
			foundIndex !=
				index {

			return Owner{},
				false
		}

		foundIndex =
			index
	}

	if foundIndex < 0 {

		return Owner{},
			false
	}

	return Owner{
			Kind: OwnerPerson,

			PersonIndex: foundIndex,
		},
		true
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
