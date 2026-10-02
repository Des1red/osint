package dependencies

import (
	"fmt"
)

func Uninstall() error {
	state,
		err :=
		readState()

	if err != nil {

		return err
	}

	manager :=
		""

	//
	// Remove in reverse installation order.
	//
	for index :=
		len(state.Dependencies) - 1; index >= 0; index-- {

		dependency :=
			&state.Dependencies[index]

		if !dependency.InstalledByOSINT {

			continue
		}

		if dependency.Package == "" {

			return fmt.Errorf(
				"dependency %s was installed by OSINT but has no recorded package",
				dependency.Name,
			)
		}

		//
		// Only resolve the package manager if we
		// actually have something to uninstall.
		//
		if manager == "" {

			manager,
				err =
				packageManager()

			if err != nil {

				return err
			}
		}

		fmt.Println(
			"removing dependency:",
			dependency.Package,
		)

		err =
			removePackage(
				manager,
				dependency.Package,
			)

		if err != nil {

			return fmt.Errorf(
				"remove dependency %s: %w",
				dependency.Name,
				err,
			)
		}

		//
		// Persist immediately so that if a later
		// dependency fails to uninstall, a retry
		// does not lose track of what was already
		// successfully removed.
		//
		dependency.InstalledByOSINT =
			false

		err =
			writeState(
				state,
			)

		if err != nil {

			return err
		}
	}

	return nil
}
