package models

import (
	"path/filepath"
)

const DependenciesStateName = "dependencies.json"

type Dependency struct {
	Name string `json:"name"`

	Executable string `json:"executable"`

	Package string `json:"package"`

	InstalledByOSINT bool `json:"installed_by_osint"`
}

type DependencyState struct {
	Dependencies []Dependency `json:"dependencies"`
}

func DependenciesStatePath() (
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
			DependenciesStateName,
		),
		nil
}
