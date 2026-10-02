package cmd

import (
	"os"

	"osint/internal/bootstrap"
	"osint/internal/bootstrap/virtualenv"
	"osint/internal/engines"
	"osint/internal/information"
	"osint/internal/logger"
)

func checkflags() {
	flags()

	if showHelp {
		help()

		os.Exit(
			0,
		)
	}

	if len(os.Args) == 1 {
		help()

		os.Exit(
			0,
		)
	}
}

func boot() {
	bootstrap.Boot()
}

func preboot() func() {
	return bootstrap.Preboot()
}

func initiate() {
	engines.Manager()

	information.Print()
}

func cleanup() {
	err :=
		virtualenv.Close()

	if err != nil {

		logger.LogError(
			"Virtual environment cleanup failed",
			err.Error(),
		)
	}
}
