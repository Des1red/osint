package verify

import (
	"fmt"
	"os"
	"osint/internal/models"
)

func VerifyCache() error {
	cacheDir, err := models.CacheDir()
	if err != nil {
		return err
	}

	info, err := os.Stat(cacheDir)

	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf(
				"cache path exists but is not a directory: %s",
				cacheDir,
			)
		}

		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf(
			"failed to check cache directory: %w",
			err,
		)
	}

	if err := os.MkdirAll(
		cacheDir,
		0755,
	); err != nil {
		return fmt.Errorf(
			"failed to create cache directory: %w",
			err,
		)
	}

	return nil
}
