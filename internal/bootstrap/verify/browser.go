package verify

import (
	"fmt"
	"os"
	"os/exec"

	"osint/internal/models"
)

func VerifyBrowser() error {
	executable,
		err :=
		chromeExecutable()

	if err != nil {
		return err
	}

	profileDirectory,
		err :=
		models.ChromeProfileDir()

	if err != nil {
		return fmt.Errorf(
			"resolve Chrome profile directory: %w",
			err,
		)
	}

	err =
		os.MkdirAll(
			profileDirectory,
			0700,
		)

	if err != nil {
		return fmt.Errorf(
			"create Chrome profile directory: %w",
			err,
		)
	}

	models.Browser.Executable =
		executable

	models.Browser.ProfileDirectory =
		profileDirectory

	return nil
}

func chromeExecutable() (
	string,
	error,
) {
	candidates :=
		[]string{
			"google-chrome",
			"google-chrome-stable",
			"chromium",
			"chromium-browser",
		}

	for _, candidate := range candidates {

		path,
			err :=
			exec.LookPath(
				candidate,
			)

		if err == nil {
			return path,
				nil
		}
	}

	return "",
		fmt.Errorf(
			"Chrome or Chromium is not installed",
		)
}
