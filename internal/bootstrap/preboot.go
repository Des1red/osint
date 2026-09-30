package bootstrap

import (
	"os"
	"osint/internal/bootstrap/lockfile"
	"osint/internal/logger"
)

func preBoot() (
	*lockfile.Lock,
	error,
) {
	processLock,
		err :=
		lockfile.Acquire()

	if err != nil {

		return nil,
			err
	}

	return processLock,
		nil
}

func Preboot() func() {
	processLock,
		err :=
		preBoot()

	if err != nil {

		logger.LogError(
			"Startup failed",
			err.Error(),
		)

		os.Exit(
			1,
		)
	}

	return func() {
		err :=
			processLock.Release()

		if err != nil {

			logger.LogError(
				"Lock release failed",
				err.Error(),
			)
		}
	}
}
