package installation

import (
	"fmt"
	"os"

	instll "github.com/Des1red/goinstall/cmd"

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
		instll.SetBinaryName(
			models.BinName,
		)

	if err != nil {
		return fmt.Errorf(
			"set binary name: %w",
			err,
		)
	}

	err =
		instll.Uninstall(
			true,
			false,
		)

	if err != nil {
		return fmt.Errorf(
			"remove binary: %w",
			err,
		)
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
