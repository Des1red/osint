package bootstrap

import (
	"fmt"
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
	fmt.Println("=============================")
	fmt.Println("Bootstrap: started")
	fmt.Println("=============================")
	fmt.Println()
	fmt.Println("initializing....")
	bootflagusage()

	if models.BootFlags.Install {

		install()

		os.Exit(
			0,
		)
	}

	if models.BootFlags.Uninstall {

		uninstall()

		os.Exit(
			0,
		)
	}

	if models.BootFlags.Debug {

		fmt.Println(
			"Debug messages on",
		)
	}
	fmt.Println("Verifying dependencies for this run.")
	verifyDependencies()

	fmt.Println()
	fmt.Println("=============================")
	fmt.Println("Bootstrap: ended")
	fmt.Println("=============================")
	fmt.Println()
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
