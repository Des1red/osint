package facebook

import "strings"

func NameLookup(
	fullName string,
) ([]FacebookResult, error) {
	candidate := strings.TrimSpace(fullName)

	if candidate == "" {
		return nil,
			nil
	}

	result, err :=
		UsernameLookup(
			candidate,
		)

	if err != nil {
		return nil,
			err
	}

	if !result.Found {
		return nil,
			nil
	}

	return []FacebookResult{
		result,
	}, nil
}
