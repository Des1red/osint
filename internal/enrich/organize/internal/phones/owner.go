package phones

type ownerKind int

const (
	ownerUnknown ownerKind = iota

	ownerTarget

	ownerPerson
)

type owner struct {
	Kind ownerKind

	PersonIndex int
}

func sameOwner(
	left owner,
	right owner,
) bool {
	if left.Kind !=
		right.Kind {

		return false
	}

	if left.Kind ==
		ownerPerson {

		return left.PersonIndex ==
			right.PersonIndex
	}

	return true
}
