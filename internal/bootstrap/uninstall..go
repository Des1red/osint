package bootstrap

import (
	"fmt"
	"os"

	"osint/internal/models"
)

func uninstall() {
	if err :=
		os.Remove(
			installPath,
		); err != nil {

		if !os.IsNotExist(
			err,
		) {

			fmt.Println(
				"error removing binary (are you root?):",
				err,
			)

			return
		}
	}

	cacheDir,
		err :=
		models.CacheDir()

	if err != nil {

		fmt.Println(
			"error finding cache directory:",
			err,
		)

		return
	}

	if err :=
		os.RemoveAll(
			cacheDir,
		); err != nil {

		fmt.Println(
			"error removing cache:",
			err,
		)

		return
	}

	appDir,
		err :=
		models.AppDir()

	if err != nil {

		fmt.Println(
			"error finding application directory:",
			err,
		)

		return
	}

	if err :=
		os.RemoveAll(
			appDir,
		); err != nil {

		fmt.Println(
			"error removing application data:",
			err,
		)

		return
	}

	fmt.Println(
		"removed",
		installPath,
	)

	fmt.Println(
		"removed cache",
		cacheDir,
	)

	fmt.Println(
		"removed application data",
		appDir,
	)
}
