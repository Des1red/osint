package installation

import (
	"fmt"
	"os"

	"osint/internal/bootstrap/installation/dependencies"
	"osint/internal/bootstrap/installation/privilege"
	"osint/internal/models"
)

func Uninstall() error {
	err :=
		privilege.EnsureUserInvocation()

	if err != nil {

		return err
	}

	//
	// Dependency state lives inside the application
	// directory, so dependencies must be handled
	// before application data is removed.
	//
	err =
		dependencies.Uninstall()

	if err != nil {

		return fmt.Errorf(
			"uninstall dependencies: %w",
			err,
		)
	}

	err =
		removeBinary()

	if err != nil {

		return err
	}

	cacheDir,
		err :=
		removeCache()

	if err != nil {

		return err
	}

	appDir,
		err :=
		removeApplicationData()

	if err != nil {

		return err
	}

	fmt.Println(
		"removed",
		models.InstallPath,
	)

	fmt.Println(
		"removed cache",
		cacheDir,
	)

	fmt.Println(
		"removed application data",
		appDir,
	)

	return nil
}

func removeBinary() error {
	command,
		err :=
		privilege.Command(
			"rm",
			"-f",
			models.InstallPath,
		)

	if err != nil {

		return err
	}

	command.Stdout =
		os.Stdout

	command.Stderr =
		os.Stderr

	command.Stdin =
		os.Stdin

	err =
		command.Run()

	if err != nil {

		return fmt.Errorf(
			"remove binary %s: %w",
			models.InstallPath,
			err,
		)
	}

	return nil
}

func removeCache() (
	string,
	error,
) {
	cacheDir,
		err :=
		models.CacheDir()

	if err != nil {

		return "",
			fmt.Errorf(
				"resolve cache directory: %w",
				err,
			)
	}

	err =
		os.RemoveAll(
			cacheDir,
		)

	if err != nil {

		return "",
			fmt.Errorf(
				"remove cache: %w",
				err,
			)
	}

	return cacheDir,
		nil
}

func removeApplicationData() (
	string,
	error,
) {
	appDir,
		err :=
		models.AppDir()

	if err != nil {

		return "",
			fmt.Errorf(
				"resolve application directory: %w",
				err,
			)
	}

	err =
		os.RemoveAll(
			appDir,
		)

	if err != nil {

		return "",
			fmt.Errorf(
				"remove application data: %w",
				err,
			)
	}

	return appDir,
		nil
}
