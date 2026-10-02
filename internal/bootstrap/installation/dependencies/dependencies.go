package dependencies

import (
	"osint/internal/models"
)

func Install() error {
	previous,
		err :=
		readState()

	if err != nil {

		return err
	}

	previousDependencies :=
		make(
			map[string]models.Dependency,
		)

	for _, dependency := range previous.Dependencies {

		previousDependencies[dependency.Name] =
			dependency
	}

	state :=
		models.DependencyState{}

	xvfb,
		err :=
		ensureXvfb(
			previousDependencies,
		)

	if err != nil {

		return err
	}

	state.Dependencies =
		append(
			state.Dependencies,
			xvfb,
		)

	err =
		writeState(
			state,
		)

	if err != nil {

		return err
	}

	browser,
		err :=
		ensureBrowser(
			previousDependencies,
		)

	if err != nil {

		return err
	}

	state.Dependencies =
		append(
			state.Dependencies,
			browser,
		)

	err =
		writeState(
			state,
		)

	if err != nil {

		return err
	}

	return nil
}
