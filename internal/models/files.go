package models

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	BinName = "osint"

	InstallPath = "/usr/local/bin/" +
		BinName

	CacheName = "osint"

	AppDirectoryName = ".osint-master"

	ChromeProfileName = "chrome-profile"
)

func GeoCachePath() (
	string,
	error,
) {
	dir,
		err :=
		CacheDir()

	if err != nil {
		return "",
			err
	}

	return filepath.Join(
			dir,
			"geo.json",
		),
		nil
}

func CacheDir() (
	string,
	error,
) {
	cacheDir,
		err :=
		os.UserCacheDir()

	if err != nil {
		return "",
			fmt.Errorf(
				"failed to get user cache directory: %w",
				err,
			)
	}

	return filepath.Join(
			cacheDir,
			CacheName,
		),
		nil
}

func AppDir() (
	string,
	error,
) {
	homeDirectory,
		err :=
		os.UserHomeDir()

	if err != nil {
		return "",
			fmt.Errorf(
				"failed to get user home directory: %w",
				err,
			)
	}

	return filepath.Join(
			homeDirectory,
			AppDirectoryName,
		),
		nil
}

func ChromeProfileDir() (
	string,
	error,
) {
	appDirectory,
		err :=
		AppDir()

	if err != nil {
		return "",
			err
	}

	return filepath.Join(
			appDirectory,
			ChromeProfileName,
		),
		nil
}
