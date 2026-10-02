package verify

import (
	"fmt"
	"os"
	"os/exec"

	"osint/internal/models"
)

var browserCandidates = []string{
	"google-chrome",
	"google-chrome-stable",
	"chromium",
	"chromium-browser",
}

func VerifyBrowser() error {
	executable,
		err :=
		BrowserExecutable()

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

func BrowserExecutable() (
	string,
	error,
) {
	for _, candidate := range browserCandidates {

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
