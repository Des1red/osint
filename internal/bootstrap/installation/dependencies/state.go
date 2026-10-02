package dependencies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"osint/internal/models"
)

func readState() (
	models.DependencyState,
	error,
) {
	path,
		err :=
		models.DependenciesStatePath()

	if err != nil {

		return models.DependencyState{},
			err
	}

	data,
		err :=
		os.ReadFile(
			path,
		)

	if os.IsNotExist(
		err,
	) {

		return models.DependencyState{},
			nil
	}

	if err != nil {

		return models.DependencyState{},
			fmt.Errorf(
				"read dependency state: %w",
				err,
			)
	}

	var state models.DependencyState

	err =
		json.Unmarshal(
			data,
			&state,
		)

	if err != nil {

		return models.DependencyState{},
			fmt.Errorf(
				"decode dependency state: %w",
				err,
			)
	}

	return state,
		nil
}

func writeState(
	state models.DependencyState,
) error {
	path,
		err :=
		models.DependenciesStatePath()

	if err != nil {

		return err
	}

	err =
		os.MkdirAll(
			filepath.Dir(
				path,
			),
			0700,
		)

	if err != nil {

		return fmt.Errorf(
			"create dependency state directory: %w",
			err,
		)
	}

	data,
		err :=
		json.MarshalIndent(
			state,
			"",
			"  ",
		)

	if err != nil {

		return fmt.Errorf(
			"encode dependency state: %w",
			err,
		)
	}

	temporary,
		err :=
		os.CreateTemp(
			filepath.Dir(
				path,
			),
			".dependencies-*.tmp",
		)

	if err != nil {

		return fmt.Errorf(
			"create temporary dependency state: %w",
			err,
		)
	}

	temporaryPath :=
		temporary.Name()

	defer os.Remove(
		temporaryPath,
	)

	err =
		temporary.Chmod(
			0600,
		)

	if err != nil {

		temporary.Close()

		return fmt.Errorf(
			"set dependency state permissions: %w",
			err,
		)
	}

	_,
		err =
		temporary.Write(
			data,
		)

	if err != nil {

		temporary.Close()

		return fmt.Errorf(
			"write dependency state: %w",
			err,
		)
	}

	err =
		temporary.Sync()

	if err != nil {

		temporary.Close()

		return fmt.Errorf(
			"sync dependency state: %w",
			err,
		)
	}

	err =
		temporary.Close()

	if err != nil {

		return fmt.Errorf(
			"close dependency state: %w",
			err,
		)
	}

	err =
		os.Rename(
			temporaryPath,
			path,
		)

	if err != nil {

		return fmt.Errorf(
			"commit dependency state: %w",
			err,
		)
	}

	return nil
}
