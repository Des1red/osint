package subjects

type OwnerKind int

const (
	OwnerUnknown OwnerKind = iota

	OwnerTarget

	OwnerPerson
)

type Owner struct {
	Kind OwnerKind

	PersonIndex int
}

func sameOwner(
	left Owner,
	right Owner,
) bool {
	if left.Kind !=
		right.Kind {

		return false
	}

	if left.Kind ==
		OwnerPerson {

		return left.PersonIndex ==
			right.PersonIndex
	}

	return true
}
