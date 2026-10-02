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

	LockFileName = "osint.lock"

	XvfbStateName = "xvfb-state.json"
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

func LockFilePath() (
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
			LockFileName,
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

func ClearChromeProfileLocks() error {
	profileDirectory,
		err :=
		ChromeProfileDir()

	if err != nil {

		return err
	}

	lockFiles :=
		[]string{
			"SingletonLock",
			"SingletonCookie",
			"SingletonSocket",
		}

	for _, name := range lockFiles {

		path :=
			filepath.Join(
				profileDirectory,
				name,
			)

		err =
			os.Remove(
				path,
			)

		if err == nil ||
			os.IsNotExist(
				err,
			) {

			continue
		}

		return fmt.Errorf(
			"failed to remove chromium profile lock %s: %w",
			name,
			err,
		)
	}

	return nil
}

func XvfbStatePath() (
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
			XvfbStateName,
		),
		nil
}
