package bootstrap

import (
	"os"

	"osint/internal/bootstrap/verify"
	"osint/internal/logger"
	"osint/internal/models"
)

const installPath = models.InstallPath

const binName = models.BinName

func bootflagusage() {
	if models.BootFlags.Install &&
		models.BootFlags.Uninstall {

		logger.LogError(
			"Boot failed",
			"Can't use install && uninstall",
		)

		os.Exit(
			0,
		)
	}
}

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

	if models.BootFlags.Install {

		logger.Info(
			"Installing OSINT...",
		)

		install()

		os.Exit(
			0,
		)
	}

	if models.BootFlags.Uninstall {

		logger.Info(
			"Uninstalling OSINT...",
		)

		uninstall()

		os.Exit(
			0,
		)
	}

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

	restoreState()

	logger.Info(
		"Runtime state prepared.",
	)
}

func restoreState() {
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
}
