package dependencies

import (
	"fmt"

	"osint/internal/bootstrap/verify"
	"osint/internal/models"
)

const browserDependencyName = "Chrome / Chromium"

func ensureBrowser(
	previous map[string]models.Dependency,
) (
	models.Dependency,
	error,
) {
	executable,
		err :=
		verify.BrowserExecutable()

	if err == nil {

		dependency :=
			models.Dependency{
				Name: browserDependencyName,

				Executable: executable,
			}

		if old,
			exists :=
			previous[browserDependencyName]; exists &&
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

	//
	// Chromium is the browser OSINT installs when
	// no supported Chrome/Chromium installation
	// already exists.
	//
	packageName :=
		"chromium"

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
				"install Chromium: %w",
				err,
			)
	}

	executable,
		err =
		verify.BrowserExecutable()

	if err != nil {

		return models.Dependency{},
			fmt.Errorf(
				"Chromium was installed but executable could not be found: %w",
				err,
			)
	}

	return models.Dependency{
			Name: browserDependencyName,

			Executable: executable,

			Package: packageName,

			InstalledByOSINT: true,
		},
		nil
}
