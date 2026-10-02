package bootstrap

import (
	"os"

	"osint/internal/bootstrap/installation"
	"osint/internal/bootstrap/verify"
	"osint/internal/bootstrap/virtualenv"
	"osint/internal/logger"
	"osint/internal/models"
)

func Boot() {
	logger.HeaderStart(
		"Bootstrap",
	)

	defer logger.HeaderEnd(
		"Bootstrap",
	)

	logger.Info(
		"Initializing...",
	)

	bootflagusage()

	handleInstallation()

	if models.BootFlags.Debug {

		logger.Info(
			"Debug messages enabled.",
		)
	}

	logger.Info(
		"Verifying dependencies for this run...",
	)

	verifyDependencies()

	logger.Info(
		"Dependencies verified.",
	)

	logger.Info(
		"Preparing runtime state...",
	)

	prepareState()

	logger.Info(
		"Runtime state prepared.",
	)

	logger.Info(
		"Preparing virtual environment...",
	)

	prepareVirtualEnvironment()
	logger.Debug(
		"Virtual Environment",
		"DISPLAY",
		os.Getenv(
			"DISPLAY",
		),
	)
	logger.Info(
		"Virtual environment ready.",
	)
}

func bootflagusage() {
	if models.BootFlags.Install &&
		models.BootFlags.Uninstall {

		logger.LogError(
			"Boot failed",
			"Can't use install && uninstall",
		)

		os.Exit(
			1,
		)
	}
}

func handleInstallation() {
	if models.BootFlags.Install {

		logger.Info(
			"Installing OSINT...",
		)

		err :=
			installation.Install()

		if err != nil {

			logger.LogError(
				"Installation failed",
				err.Error(),
			)

			os.Exit(
				1,
			)
		}

		logger.Info(
			"OSINT installed successfully.",
		)

		os.Exit(
			0,
		)
	}

	if models.BootFlags.Uninstall {

		logger.Info(
			"Uninstalling OSINT...",
		)

		err :=
			installation.Uninstall()

		if err != nil {

			logger.LogError(
				"Uninstallation failed",
				err.Error(),
			)

			os.Exit(
				1,
			)
		}

		logger.Info(
			"OSINT uninstalled successfully.",
		)

		os.Exit(
			0,
		)
	}
}

func verifyDependencies() {
	verify.VerifyKeys()

	if err :=
		verify.VerifyCache(); err != nil {

		logger.LogError(
			"Cache verification failed",
			err.Error(),
		)

		os.Exit(
			1,
		)
	}

	if err :=
		verify.VerifyBrowser(); err != nil {

		logger.LogError(
			"Browser verification failed",
			err.Error(),
		)

		os.Exit(
			1,
		)
	}

	if err :=
		verify.VerifyXvfb(); err != nil {

		logger.LogError(
			"Xvfb verification failed",
			err.Error(),
		)

		os.Exit(
			1,
		)
	}
}

func prepareState() {
	err :=
		models.ClearChromeProfileLocks()

	if err != nil {

		logger.LogError(
			"Failed to clear Chromium profile locks",
			err.Error(),
		)

		os.Exit(
			1,
		)
	}

	err =
		virtualenv.PrepareState()

	if err != nil {

		logger.LogError(
			"Failed to prepare virtual environment state",
			err.Error(),
		)

		os.Exit(
			1,
		)
	}
}

func prepareVirtualEnvironment() {
	err :=
		virtualenv.Xvfb()

	if err != nil {

		logger.LogError(
			"Virtual environment failed",
			err.Error(),
		)

		os.Exit(
			1,
		)
	}
}
