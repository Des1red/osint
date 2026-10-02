package dependencies

import (
	"fmt"

	"osint/internal/bootstrap/verify"
	"osint/internal/models"
)

func ensureXvfb(
	previous map[string]models.Dependency,
) (
	models.Dependency,
	error,
) {
	const name = "Xvfb"

	executable,
		err :=
		verify.XvfbExecutable()

	if err == nil {

		dependency :=
			models.Dependency{
				Name: name,

				Executable: executable,
			}

		if old,
			exists :=
			previous[name]; exists &&
			old.InstalledByOSINT {

			dependency.Package =
				old.Package

			dependency.InstalledByOSINT =
				true
		}

		return dependency,
			nil
	}

	manager,
		err :=
		packageManager()

	if err != nil {

		return models.Dependency{},
			err
	}

	var packageName string

	switch manager {

	case packageManagerDNF:

		packageName =
			"xorg-x11-server-Xvfb"

	case packageManagerAPT:

		packageName =
			"xvfb"

	default:

		return models.Dependency{},
			fmt.Errorf(
				"unsupported package manager %q",
				manager,
			)
	}

	fmt.Println(
		"installing dependency:",
		packageName,
	)

	err =
		installPackage(
			manager,
			packageName,
		)

	if err != nil {

		return models.Dependency{},
			fmt.Errorf(
				"install Xvfb: %w",
				err,
			)
	}

	executable,
		err =
		verify.XvfbExecutable()

	if err != nil {

		return models.Dependency{},
			fmt.Errorf(
				"Xvfb was installed but executable could not be found: %w",
				err,
			)
	}

	return models.Dependency{
			Name: name,

			Executable: executable,

			Package: packageName,

			InstalledByOSINT: true,
		},
		nil
}
